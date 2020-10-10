package aws

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/provider/aws/utils"
	"github.com/mensaah/reka/types"
)

func getS3BucketRegion(cfg aws.Config, bucketName string, logger *log.Entry) (string, error) {

	region, err := s3manager.GetBucketRegion(context.Background(), cfg, bucketName, "us-west-2")
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			return "", fmt.Errorf("unable to find bucket %s's region not found", bucketName)
		}
		return "", err
	}
	logger.Debugf("Bucket %s is in %s region\n", bucketName, region)
	return region, err
}

func getS3BucketTags(svc *s3.Client, bucketName string, logger *log.Entry) (types.ResourceTags, error) {
	logger.Debug("Fetching S3 Tags")
	input := &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketName),
	}

	req := svc.GetBucketTaggingRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return types.ResourceTags{}, aerr
			}
		}
		return types.ResourceTags{}, err
	}
	// https://stackoverflow.com/a/48554123/7167357
	tags := utils.ParseResourceTags(*(*[]utils.AWSTag)(unsafe.Pointer(&result.TagSet)))
	return tags, nil
}

// returns only s3Bucket IDs of unprotected s3 instances
func getS3BucketsDetails(svc *s3.Client, cfg aws.Config, output *s3.ListBucketsResponse, logger *log.Entry) ([]*types.Resource, error) {
	var s3Buckets []*types.Resource
	for _, s3Bucket := range output.Buckets {
		// Get tags
		tags, err := getS3BucketTags(svc, *s3Bucket.Name, logger)
		if err != nil {
			logger.Error(err)
		}
		tags["creation-date"] = (*s3Bucket.CreationDate).String()
		// Get region
		s3Region, err := getS3BucketRegion(cfg, *s3Bucket.Name, logger)
		if err != nil {
			logger.Errorf("Could not get region for Bucket %s", *s3Bucket.Name)
			continue
		}
		s3 := NewResource(*s3Bucket.Name, s3Name)
		s3.State = types.Running
		s3.Region = s3Region
		// Get CreationDate by getting LaunchTime of attached Volume
		s3.CreationDate = *s3Bucket.CreationDate
		s3.Tags = tags
		s3Buckets = append(s3Buckets, s3)
	}

	logger.Debugf("Found %d buckets", len(s3Buckets))

	return s3Buckets, nil
}

// GetAllS3Buckets Get all s3Buckets
func getAllS3Buckets(cfg aws.Config, logger *log.Entry) ([]*types.Resource, error) {
	logger.Debug("Fetching S3 Buckets")
	svc := s3.New(cfg)
	params := &s3.ListBucketsInput{}

	// Build the request with its input parameters
	req := svc.ListBucketsRequest(params)

	// Send the request, and get the response or error back
	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	buckets, err := getS3BucketsDetails(svc, cfg, resp, logger)
	if err != nil {
		return nil, err
	}
	return buckets, nil
}

// Destroys a Single Bucket
func destroyBucket(svc *s3.Client, bucket *types.Resource, logger *log.Entry) error {
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket.UUID),
	}

	req := svc.DeleteBucketRequest(input)
	_, err := req.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return aerr
			}
		}
		return err
	}

	return nil
}

func destroyS3Buckets(cfg aws.Config, s3Buckets []*types.Resource, logger *log.Entry) error {
	bucketsPerRegion := make(map[string][]*types.Resource)
	delCount := 0
	if len(s3Buckets) <= 0 {
		return nil
	}

	for _, bucket := range s3Buckets {
		bucketsPerRegion[bucket.Region] = append(bucketsPerRegion[bucket.Region], bucket)
	}

	// TODO Use Goroutines
	for region, buckets := range bucketsPerRegion {
		svc := s3.New(cfg)
		svc.Client.Config.Region = region
		for _, bucket := range buckets {
			err := destroyBucket(svc, bucket, logger)
			if err != nil {
				logger.Errorf("Failed to delete Bucket %d - Error %s ", bucket.ID, err.Error())
				bucket.DestroyError = err
			} else {
				delCount++
			}
		}
	}
	logger.Infof("Destroyed %d S3 buckets", delCount)
	return nil
}
