package aws

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider"
	"github.com/mensaah/reka/resource"
)

var (
	providerName     = "aws"
	logger           *log.Entry
	resourceManagers map[string]*resource.Manager
)

func GetName() string {
	return providerName
}

// NewResource Returns a new Resource object
func NewResource(id, manager string) *resource.Resource {
	resource := resource.Resource{}
	resource.UUID = id
	resource.Manager = resourceManagers[manager]

	return &resource
}

// NewProvider : Creates a New AWS Provider
func NewProvider() (*provider.Provider, error) {

	aws := provider.Provider{}
	aws.Name = providerName

	logFile := fmt.Sprintf("%s/logger.log", config.GetConfig().LogPath)
	logger = config.GetLogger(providerName, logFile)
	// Setup Logger
	aws.Logger = logger

	cfg := config.GetConfig()

	ec2Manager := newEC2Manager(cfg, logFile)
	s3Manager := newS3Manager(cfg, logFile)

	resourceManagers = map[string]*resource.Manager{
		ec2Manager.Name: &ec2Manager,
		s3Manager.Name:  &s3Manager,
	}

	aws.Managers = resourceManagers
	return &aws, nil
}
