package gcp

import (
	"context"

	"cloud.google.com/go/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
)

func getAllBuckets(cfg *config.Gcp) ([]*resource.Resource, error) {
	var buckets []*resource.Resource
	log.Debug("Fetching Cloud storage buckets")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return []*resource.Resource{}, err
	}
	it := client.Buckets(ctx, cfg.ProjectId)
	for {
		bucketAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf(err.Error())
		}
		bucket := NewResource(bucketAttrs.Name, cloudStorageName)
		bucket.Status = resource.Running
		bucket.CreationDate = bucketAttrs.Created
		bucket.Tags = bucketAttrs.Labels
		buckets = append(buckets, bucket)
	}
	log.Debugf("Found %s storage buckets", len(buckets))
	return buckets, nil
}

// Empties a Bucket
func emptyBucket(bucket *resource.Resource) error {
	return nil
}

func destroyBuckets(cfg *config.Gcp, buckets []*resource.Resource) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	for _, bucket := range buckets {
		if err := client.Bucket(bucket.UUID).Delete(ctx); err != nil {
			log.Errorf("Failed deleting bucket %s: %s", bucket.UUID, err.Error())
		}
	}
	return nil
}
