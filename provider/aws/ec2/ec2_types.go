package ec2

import (
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/mensaah/reka/provider"
)

const (
	Name     = "EC2"
	LongName = "Elastic Compute Cloud"
)

// ResourceManager : Implements the ResourceManager Interface to expose methods implemented by each resource module
type ResourceManager struct {
	Name     string
	LongName string
	config   aws.Config
}

func InitManager(cfg aws.Config) ResourceManager {
	EC2s := ResourceManager{}
	EC2s.Name = Name
	EC2s.LongName = LongName
	EC2s.config = cfg

	return EC2s
}

func (r ResourceManager) GetName() string {
	return r.Name
}

func (r ResourceManager) GetAll() ([]*provider.Resource, error) {
	region := "us-east-2"
	return GetAllEC2Instances(r.config, region)
}

func (r ResourceManager) Destroy(resources []*provider.Resource) error {
	return TerminateEC2Instances(r.config, resources)
}

func (r ResourceManager) Stop(resources []*provider.Resource) error {
	return StopEC2Instances(r.config, resources)
}

func (r ResourceManager) Resume(resources []*provider.Resource) error {
	return StartEC2Instances(r.config, resources)
}

func (r ResourceManager) GetReapable(config provider.Config) ([]*provider.Resource, error) {
	return []*provider.Resource{}, nil
}

// NewEC2 Returns a new EC2 Resource object
func NewEC2(id string) *provider.Resource {
	resource := provider.Resource{}
	resource.ID = id
	resource.Name = Name
	resource.LongName = LongName

	return &resource
}
