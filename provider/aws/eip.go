package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/mensaah/reka/resource"
)

// returns only ip IDs of unprotected ec2 ips
func getIPDetails(svc *ec2.Client, output *ec2.DescribeAddressesOutput, region string) ([]*resource.Resource, error) {
	var eIps []*resource.Resource
	eipLogger.Debug("Fetching EIp Details")
	for _, ip := range output.Addresses {
		tags := make(resource.Tags)
		for _, t := range ip.Tags {
			tags[*t.Key] = *t.Value
		}
		eResource := NewResource(*ip.AllocationId, eipName)
		eResource.Region = region
		eResource.Tags = tags
		eResource.Status = resource.Running
		if *ip.AssociationId == "" || *ip.InstanceId == "" {
			eResource.Status = resource.Unused
		}
		eIps = append(eIps, eResource)
	}

	return eIps, nil
}

// GetAllIPAddresses Get all ips
func GetAllIPAddresses(cfg aws.Config) ([]*resource.Resource, error) {
	eipLogger.Debug("Fetching EIp Ips")

	svc := ec2.NewFromConfig(cfg)
	params := &ec2.DescribeAddressesInput{}

	// Build the request with its input parameters
	resp, err := svc.DescribeAddresses(context.TODO(), params)

	if err != nil {
		return nil, err
	}
	ips, err := getIPDetails(svc, resp, cfg.Region)
	if err != nil {
		return nil, err
	}
	eipLogger.Debugf("Found %d EIp ips", len(ips))
	return ips, nil
}

// TerminateIPAddresses Shutdown ips
func TerminateIPAddresses(cfg aws.Config, ips []*resource.Resource) error {
	svc := ec2.NewFromConfig(cfg)
	var targetIps []*resource.Resource

	for _, ip := range ips {
		if ip.IsStopped() || ip.IsActive() {
			targetIps = append(targetIps, ip)
		}
	}

	if len(targetIps) <= 0 {
		return nil
	}

	eipLogger.Debug("Terminating Ips ", targetIps, " ...")

	for _, ip := range targetIps {
		if ip.Status != resource.Unused {
			params := &ec2.DisassociateAddressInput{
				AssociationId: &ip.UUID,
			}
			_, err := svc.DisassociateAddress(context.TODO(), params)
			if err != nil {
				eipLogger.Errorf("Failed to Dissociated IP %s: %s", ip.UUID, err)
				continue
			}
		}
		params := &ec2.ReleaseAddressInput{
			AllocationId: &ip.UUID,
		}
		_, err := svc.ReleaseAddress(context.TODO(), params)
		if err != nil {
			eipLogger.Errorf("Failed to Delete VOlume %s: %s", ip.UUID, err)
		}
	}
	return nil
}
