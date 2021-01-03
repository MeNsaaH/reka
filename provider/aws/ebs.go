package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/mensaah/reka/resource"
)

// returns only volume IDs of unprotected ec2 volumes
func getVolumeDetails(svc *ec2.Client, output *ec2.DescribeVolumesOutput, region string) ([]*resource.Resource, error) {
	var ebsVolumes []*resource.Resource
	ebsLogger.Debug("Fetching Ebs Details")
	for _, volume := range output.Volumes {
		tags := make(resource.Tags)
		for _, t := range volume.Tags {
			tags[*t.Key] = *t.Value
		}
		// We need the creation-date when parsing Tags for relative defintions
		// Ebs Volumes Launch Time is not the creation date. It's the time it was last launched.
		// TODO To get the creation date we might want to get the creation date of the EBS attached to the Ebs instead
		tags["creation-date"] = (*volume.CreateTime).String()
		ebsResource := NewResource(*volume.VolumeId, ec2Name)
		ebsResource.Region = region
		// Get CreationDate by getting LaunchTime of attached Volume
		ebsResource.CreationDate = *volume.CreateTime
		ebsResource.Tags = tags
		if len(volume.Attachments) == 0 {
			ebsResource.Status = resource.Unused
		} else {

			ebsResource.Status = resource.Running
		}
		ebsVolumes = append(ebsVolumes, ebsResource)
	}

	return ebsVolumes, nil
}

// GetAllEbsVolumes Get all volumes
func GetAllEbsVolumes(cfg aws.Config) ([]*resource.Resource, error) {
	ebsLogger.Debug("Fetching Ebs Volumes")

	svc := ec2.NewFromConfig(cfg)
	params := &ec2.DescribeVolumesInput{}

	// Build the request with its input parameters
	resp, err := svc.DescribeVolumes(context.TODO(), params)

	if err != nil {
		return nil, err
	}
	volumes, err := getVolumeDetails(svc, resp, cfg.Region)
	if err != nil {
		return nil, err
	}
	ebsLogger.Debugf("Found %d Ebs volumes", len(volumes))
	return volumes, nil
}

// TerminateEbsVolumes Shutdown volumes
func TerminateEbsVolumes(cfg aws.Config, volumes []*resource.Resource) error {
	svc := ec2.NewFromConfig(cfg)
	var volumeIds []string

	for _, volume := range volumes {
		if volume.IsStopped() || volume.IsActive() {
			volumeIds = append(volumeIds, volume.UUID)
		}
	}

	if len(volumeIds) <= 0 {
		return nil
	}

	ebsLogger.Debug("Terminating Ebs Volumes ", volumeIds, " ...")

	for _, volumeId := range volumeIds {
		params := &ec2.DeleteVolumeInput{
			VolumeId: &volumeId,
		}

		_, err := svc.DeleteVolume(context.TODO(), params)
		// TODO Attach error to specific volume where the error occurred if possible
		if err != nil {
			ebsLogger.Errorf("Failed to Delete VOlume %s: %s", volumeId, err)
		}
	}
	return nil
}
