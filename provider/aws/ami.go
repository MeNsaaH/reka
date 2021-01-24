package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/mensaah/reka/resource"
)

// returns only image IDs of unprotected ec2 images
func getImageDetails(svc *ec2.Client, output *ec2.DescribeImagesOutput, region string) ([]*resource.Resource, error) {
	var images []*resource.Resource
	amiLogger.Debug("Fetching image Details")
	for _, image := range output.Images {
		tags := make(resource.Tags)
		for _, t := range image.Tags {
			tags[*t.Key] = *t.Value
		}
		tags["creation-date"] = *image.CreationDate
		amiResource := NewResource(*image.ImageId, ec2Name)
		amiResource.Region = *image.ImageLocation
		// Get CreationDate by getting LaunchTime of attached Image
		amiLogger.Debugf("TIME: %s", *image.CreationDate)
		// amiResource.CreationDate = *image.CreationDate
		amiResource.Tags = tags
		amiResource.Status = resource.Running
		images = append(images, amiResource)
	}
	return images, nil
}

// GetAllImages Get all images
func GetAllImages(cfg aws.Config) ([]*resource.Resource, error) {
	amiLogger.Debug("Fetching images...")

	svc := ec2.NewFromConfig(cfg)
	params := &ec2.DescribeImagesInput{
		Owners: []string{"self"},
	}

	// Build the request with its input parameters
	resp, err := svc.DescribeImages(context.TODO(), params)

	if err != nil {
		return nil, err
	}
	images, err := getImageDetails(svc, resp, cfg.Region)
	if err != nil {
		return nil, err
	}
	amiLogger.Debugf("Found %d Images", len(images))
	return images, nil
}

// Terminatemages deregisters images
func TerminateImages(cfg aws.Config, images []*resource.Resource) error {
	svc := ec2.NewFromConfig(cfg)
	var imageIds []string

	for _, image := range images {
		if image.IsStopped() || image.IsActive() {
			imageIds = append(imageIds, image.UUID)
		}
	}

	if len(imageIds) <= 0 {
		return nil
	}

	amiLogger.Debug("Terminating Ami Images ", imageIds, " ...")

	for _, imageId := range imageIds {
		params := &ec2.DeregisterImageInput{
			ImageId: &imageId,
		}

		_, err := svc.DeregisterImage(context.TODO(), params)
		if err != nil {
			amiLogger.Errorf("Failed to Deregistering Image %s: %s", imageId, err)
		}
	}
	return nil
}
