package aws

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/mensaah/reka/provider/aws/utils"
	"github.com/mensaah/reka/resource"
)

// returns only instance IDs of unprotected ec2 instances
func getInstanceDetails(svc *ec2.Client, output *ec2.DescribeInstancesOutput, region string) ([]*resource.Resource, error) {
	var ec2Instances []*resource.Resource
	ec2Logger.Debug("Fetching EC2 Details")
	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			// https://stackoverflow.com/a/48554123/7167357
			tags := utils.ParseTags(*(*[]*utils.AWSTag)(unsafe.Pointer(&instance.Tags)))

			// We need the creation-date when parsing Tags for relative defintions
			// EC2 Instances Launch Time is not the creation date. It's the time it was last launched.
			// To get the creation date we might want to get the creation date of the EBS attached to the EC2 instead
			tags["creation-date"] = (*instance.LaunchTime).String()
			ec2 := NewResource(*instance.InstanceId, ec2Name)
			ec2.Region = region
			// Get CreationDate by getting LaunchTime of attached Volume
			ec2.CreationDate = *instance.LaunchTime
			ec2.Tags = tags
			ec2.State = utils.GetResourceState(*instance.State.Code)
			ec2Instances = append(ec2Instances, ec2)
		}
	}

	return ec2Instances, nil
}

// GetAllEC2Instances Get all instances
func GetAllEC2Instances(cfg aws.Config) ([]*resource.Resource, error) {
	ec2Logger.Debug("Fetching EC2 Instances")

	svc := ec2.NewFromConfig(cfg)
	params := &ec2.DescribeInstancesInput{}

	// Build the request with its input parameters
	resp, err := svc.DescribeInstances(context.Background(), params)

	if err != nil {
		return nil, err
	}
	instances, err := getInstanceDetails(svc, resp, cfg.Region)
	if err != nil {
		return nil, err
	}
	ec2Logger.Debugf("Found %d EC2 instances", len(instances))
	return instances, nil
}

// StopEC2Instances Stop Running Instances
func StopEC2Instances(cfg aws.Config, instances []*resource.Resource) error {
	svc := ec2.NewFromConfig(cfg)
	var instanceIds []*string

	for _, instance := range instances {
		if instance.IsActive() {
			instanceIds = append(instanceIds, &instance.UUID)
		}
	}

	if len(instanceIds) <= 0 {
		return nil
	}

	ec2Logger.Debug("Stopping EC2 Instances ", instanceIds, " ...")

	params := &ec2.StopInstancesInput{
		InstanceIds: instanceIds,
	}

	resp, err := svc.StopInstances(context.Background(), params)
	// TODO Attach error to specific instance where the error occurred if possible
	if err != nil {
		fmt.Println(resp, err)
	}
	return err
}

// ResumeEC2Instances Resume Stopped instances
func ResumeEC2Instances(cfg aws.Config, instances []*resource.Resource) error {
	svc := ec2.NewFromConfig(cfg)
	var instanceIds []*string

	for _, instance := range instances {
		if instance.IsStopped() {
			instanceIds = append(instanceIds, &instance.UUID)
		}
	}

	if len(instanceIds) <= 0 {
		return nil
	}

	params := &ec2.StartInstancesInput{
		InstanceIds: instanceIds,
	}
	ec2Logger.Debug("Starting EC2 Instances ", instanceIds, " ...")

	resp, err := svc.StartInstances(context.Background(), params)
	// TODO Attach error to specific instance where the error occurred if possible
	if err != nil {
		fmt.Println(resp, err)
	}
	return err
}

// TerminateEC2Instances Shutdown instances
func TerminateEC2Instances(cfg aws.Config, instances []*resource.Resource) error {
	svc := ec2.NewFromConfig(cfg)
	var instanceIds []*string

	for _, instance := range instances {
		if instance.IsStopped() || instance.IsActive() {
			instanceIds = append(instanceIds, &instance.UUID)
		}
	}

	if len(instanceIds) <= 0 {
		return nil
	}

	ec2Logger.Debug("Terminating EC2 Instances ", instanceIds, " ...")

	params := &ec2.TerminateInstancesInput{
		InstanceIds: instanceIds,
	}

	resp, err := svc.TerminateInstances(context.Background(), params)
	// TODO Attach error to specific instance where the error occurred if possible
	if err != nil {
		fmt.Println(resp, err)
	}
	return err
}
