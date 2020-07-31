package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/mensaah/reka/provider"
)

const (
	Name     = "S3"
	LongName = "Simple Storage Service"
)

// ResourceManager : Implements the ResourceManager Interface to expose methods implemented by each resource module
type ResourceManager struct {
	provider.DefaultResourceManager
	config aws.Config
}

func InitManager(cfg aws.Config) ResourceManager {
	S3s := ResourceManager{}
	S3s.Name = Name
	S3s.LongName = LongName
	S3s.config = cfg

	return S3s
}

func (r ResourceManager) GetName() string {
	return r.Name
}

func (r ResourceManager) GetAll() ([]*provider.Resource, error) {
	return getAllS3Buckets(r.config)
}

func (r ResourceManager) Destroy(resources []*provider.Resource) error {
	return destroyS3Buckets(r.config, resources)
}

func (r ResourceManager) GetReapable(config provider.Config) ([]*provider.Resource, error) {
	return []*provider.Resource{}, nil
}

// NewS3 Returns a new S3 Resource object
func NewS3(id string) *provider.Resource {
	resource := provider.Resource{}
	resource.ID = id
	resource.Name = Name
	resource.LongName = LongName

	return &resource
}
