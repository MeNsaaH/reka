package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/resource"
)

func getS3BucketRegion(cfg aws.Config, bucketName string) (string, error) {

	region, err := s3manager.GetBucketRegion(context.TODO(), s3.NewFromConfig(cfg), bucketName)
	if err != nil {
		var notFoundErr *s3Types.NoSuchBucket
		if errors.As(err, &notFoundErr) {
			log.Printf("scan failed because the table was not found, %v",
				notFoundErr.ErrorMessage())
			return "", fmt.Errorf("unable to find bucket %s's region not found", bucketName)
		}
		return "", err
	}
	s3Logger.Debugf("Bucket %s is in %s region\n", bucketName, region)
	return region, err
}

func getS3BucketTags(svc *s3.Client, bucketName string) (resource.Tags, error) {
	s3Logger.Debug("Fetching S3 Tags")
	input := &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketName),
	}

	result, err := svc.GetBucketTagging(context.TODO(), input)
	if err != nil {
		return resource.Tags{}, err
	}
	tags := make(resource.Tags)
	for _, t := range result.TagSet {
		tags[*t.Key] = *t.Value
	}
	return tags, nil
}

// returns only s3Bucket IDs of unprotected s3 instances
func getS3BucketsDetails(svc *s3.Client, cfg aws.Config, output *s3.ListBucketsOutput) ([]*resource.Resource, error) {
	var s3Buckets []*resource.Resource
	for _, s3Bucket := range output.Buckets {
		// Get tags
		tags, err := getS3BucketTags(svc, *s3Bucket.Name)
		if err != nil {
			s3Logger.Error(err)
		}
		tags["creation-date"] = (*s3Bucket.CreationDate).String()
		// Get region
		s3Region, err := getS3BucketRegion(cfg, *s3Bucket.Name)
		if err != nil {
			s3Logger.Errorf("Could not get region for Bucket %s", *s3Bucket.Name)
			continue
		}
		s3 := NewResource(*s3Bucket.Name, s3Name)
		s3.Status = resource.Running
		s3.Region = s3Region
		s3.CreationDate = *s3Bucket.CreationDate
		s3.Tags = tags
		s3Buckets = append(s3Buckets, s3)
	}

	s3Logger.Debugf("Found %d buckets", len(s3Buckets))

	return s3Buckets, nil
}

// GetAllS3Buckets Get all s3Buckets
func getAllS3Buckets(cfg aws.Config) ([]*resource.Resource, error) {
	s3Logger.Debug("Fetching S3 Buckets")
	svc := s3.NewFromConfig(cfg)
	params := &s3.ListBucketsInput{}

	// Build the request with its input parameters
	resp, err := svc.ListBuckets(context.TODO(), params)
	if err != nil {
		return nil, err
	}
	buckets, err := getS3BucketsDetails(svc, cfg, resp)
	if err != nil {
		return nil, err
	}
	return buckets, nil
}

// Empties a Bucket
func emptyBucket(svc *s3.Client, bucket *resource.Resource) error {
	return nil
}

// Destroys a Single Bucket
func destroyBucket(svc *s3.Client, bucket *resource.Resource) error {
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket.UUID),
	}

	_, err := svc.DeleteBucket(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

func destroyS3Buckets(cfg aws.Config, s3Buckets []*resource.Resource) error {
	bucketsPerRegion := make(map[string][]*resource.Resource)
	delCount := 0
	if len(s3Buckets) <= 0 {
		return nil
	}

	for _, bucket := range s3Buckets {
		bucketsPerRegion[bucket.Region] = append(bucketsPerRegion[bucket.Region], bucket)
	}

	// TODO Use Goroutines and also destroy all objects in the bucket before executing destroy on bucket
	for region, buckets := range bucketsPerRegion {
		svc := s3.NewFromConfig(cfg, func(options *s3.Options) {
			options.Region = region
		})
		for _, bucket := range buckets {
			err := destroyBucket(svc, bucket)
			if err != nil {
				s3Logger.Errorf("Failed to delete Bucket %d - Error %s ", bucket.ID, err.Error())
			} else {
				delCount++
			}
		}
	}
	s3Logger.Infof("Destroyed %d S3 buckets", delCount)
	return nil
}
