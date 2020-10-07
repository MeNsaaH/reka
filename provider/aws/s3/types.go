package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/types"
)

const (
	// Name of resource
	Name = "S3"
	// LongName descriptive name for resource
	LongName = "Simple Storage Service"
)

var logger *log.Entry

// ResourceManager : Implements the ResourceManager Interface to expose methods implemented by each resource module
type ResourceManager struct {
	types.DefaultResourceManager
	config aws.Config
}

func InitManager() ResourceManager {
	mgr := ResourceManager{}
	mgr.Name = Name
	mgr.LongName = LongName
	mgr.config = config.GetAWS().Config
	mgr.Logger = types.GetLogger(mgr.Name)

	logger = mgr.Logger

	return mgr
}

func (r ResourceManager) GetName() string {
	return r.Name
}

func (r ResourceManager) GetAll() ([]*types.Resource, error) {
	return getAllS3Buckets(r.config)
}

func (r ResourceManager) Destroy(resources []*types.Resource) error {
	return destroyS3Buckets(r.config, resources)
}

func (r ResourceManager) GetReapable() ([]*types.Resource, error) {
	return []*types.Resource{}, nil
}

func (r ResourceManager) GetLogger() *log.Entry {
	return r.Logger
}

// New Returns a new S3 Resource object
func New(id string) *types.Resource {
	resource := types.Resource{}
	resource.ID = id
	resource.Name = Name
	resource.LongName = LongName

	return &resource
}
