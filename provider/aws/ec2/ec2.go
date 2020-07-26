package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/provider"
	"github.com/mensaah/reka/provider/aws/utils"
)

// returns only instance IDs of unprotected ec2 instances
func getInstanceDetails(svc *ec2.Client, output *ec2.DescribeInstancesResponse, region string) ([]*provider.Resource, error) {
	var ec2Instances []*provider.Resource
	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			tags := utils.ParseResourceTags(instance.Tags)
			// We need the creation-date when parsing Tags for relative defintions
			// EC2 Instances Launch Time is not the creation date. It's the time it was last launched.
			// To get the creation date we might want to get the creation date of the EBS attached to the EC2 instead
			tags["creation-date"] = (*instance.LaunchTime).String()
			ec2 := NewEC2(*instance.InstanceId)
			ec2.Region = region
			// Get CreationDate by getting LaunchTime of attached Volume
			ec2.CreationDate = *instance.LaunchTime
			ec2.Tags = tags
			ec2.State = utils.GetResourceState(*instance.State.Code)
			log.Info(tags)
			ec2Instances = append(ec2Instances, ec2)
		}
	}

	return ec2Instances, nil
}

// GetAllEC2Instances Get all instances
func GetAllEC2Instances(cfg aws.Config, region string) ([]*provider.Resource, error) {
	svc := ec2.New(cfg)
	params := &ec2.DescribeInstancesInput{}

	// Build the request with its input parameters
	req := svc.DescribeInstancesRequest(params)

	// Send the request, and get the response or error back
	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	instances, err := getInstanceDetails(svc, resp, region)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

// StopEC2Instances Stop Running Instances
func StopEC2Instances(cfg aws.Config, instances []*provider.Resource) error {
	svc := ec2.New(cfg)
	var instanceIds []string

	for _, instance := range instances {
		if instance.IsActive() {
			instanceIds = append(instanceIds, instance.ID)
		}
	}

	if len(instanceIds) <= 0 {
		return nil
	}

	log.Debug("Stopping EC2 Instances ", instanceIds, " ...")

	params := &ec2.StopInstancesInput{
		InstanceIds: instanceIds,
	}

	req := svc.StopInstancesRequest(params)
	resp, err := req.Send(context.Background())
	// TODO Attach error to specific instance where the error occurred if possible
	if err != nil {
		fmt.Println(resp, err)
	}
	return err
}

// StartEC2Instances Start Stopped instances
func StartEC2Instances(cfg aws.Config, instances []*provider.Resource) error {
	svc := ec2.New(cfg)
	var instanceIds []string

	for _, instance := range instances {
		if instance.IsStopped() {
			instanceIds = append(instanceIds, instance.ID)
		}
	}

	if len(instanceIds) <= 0 {
		return nil
	}

	params := &ec2.StartInstancesInput{
		InstanceIds: instanceIds,
	}
	log.Debug("Starting EC2 Instances ", instanceIds, " ...")

	req := svc.StartInstancesRequest(params)
	resp, err := req.Send(context.Background())
	// TODO Attach error to specific instance where the error occurred if possible
	if err != nil {
		fmt.Println(resp, err)
	}
	return err
}

// StartEC2Instances Start Stopped instances
func TerminateEC2Instances(cfg aws.Config, instances []*provider.Resource) error {
	svc := ec2.New(cfg)
	var instanceIds []string

	for _, instance := range instances {
		if instance.IsStopped() || instance.IsActive() {
			instanceIds = append(instanceIds, instance.ID)
		}
	}

	if len(instanceIds) <= 0 {
		return nil
	}

	log.Debug("Terminating EC2 Instances ", instanceIds, " ...")

	params := &ec2.TerminateInstancesInput{
		InstanceIds: instanceIds,
	}

	req := svc.TerminateInstancesRequest(params)
	resp, err := req.Send(context.Background())
	// TODO Attach error to specific instance where the error occurred if possible
	if err != nil {
		fmt.Println(resp, err)
	}
	return err
}
