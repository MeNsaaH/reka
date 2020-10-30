package aws

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/types"
)

var (
	providerName     = "aws"
	logger           *log.Entry
	resourceManagers map[string]*types.ResourceManager
)

func GetName() string {
	return providerName
}

// NewResource Returns a new Resource object
func NewResource(id, manager string) *types.Resource {
	resource := types.Resource{}
	resource.UUID = id
	resource.Manager = resourceManagers[manager]

	return &resource
}

// NewProvider : Creates a New AWS Provider
func NewProvider() (*types.Provider, error) {

	aws := types.Provider{}
	aws.Name = providerName

	logFile := fmt.Sprintf("%s/logger.log", config.GetConfig().LogPath)
	logger = types.GetLogger(providerName, logFile)
	// Setup Logger
	aws.Logger = logger

	cfg := config.GetConfig()

	ec2Manager := newEC2Manager(cfg, logFile)
	s3Manager := newS3Manager(cfg, logFile)

	resourceManagers = map[string]*types.ResourceManager{
		ec2Manager.Name: &ec2Manager,
		s3Manager.Name:  &s3Manager,
	}

	aws.ResourceManagers = resourceManagers
	return &aws, nil
}
