package ec2

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/types"
)

const (
	Name     = "EC2"
	LongName = "Elastic Compute Cloud"
)

// ResourceManager : Implements the ResourceManager Interface to expose methods implemented by each resource module
type ResourceManager struct {
	types.DefaultResourceManager
	config aws.Config
}

var (
	logger *log.Entry
)

func InitManager() ResourceManager {
	mgr := ResourceManager{}
	mgr.Name = Name
	mgr.LongName = LongName
	mgr.Provider = "aws"
	mgr.config = config.GetAWS().Config
	mgr.Logger = types.GetLogger(mgr.Name)

	logger = mgr.Logger

	return mgr
}

func (r ResourceManager) GetName() string {
	return r.Name
}

func (r ResourceManager) GetAll() ([]*types.Resource, error) {
	region := "us-east-2"
	return GetAllEC2Instances(r.config, region)
}

func (r ResourceManager) Destroy(resources []*types.Resource) error {
	return TerminateEC2Instances(r.config, resources)
}

func (r ResourceManager) Stop(resources []*types.Resource) error {
	return StopEC2Instances(r.config, resources)
}

func (r ResourceManager) Resume(resources []*types.Resource) error {
	return StartEC2Instances(r.config, resources)
}

func (r ResourceManager) GetReapable() ([]*types.Resource, error) {
	return []*types.Resource{}, nil
}

func (r ResourceManager) GetLogger() *log.Entry {
	return r.Logger
}

// New Returns a new EC2 Resource object
func New(id string) *types.Resource {
	resource := types.Resource{}
	resource.ID = id
	resource.Name = Name
	resource.LongName = LongName

	return &resource
}
