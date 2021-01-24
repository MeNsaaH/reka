package aws

import (
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
)

// Manages Ebs instances on the AWS.
// Ebs resources support stopping/resuming and terminating instances.

var amiManager resource.Manager

const (
	// Name of resource
	amiName = "ami"
	// LongName descriptive name for resource
	amiLongName = "Amazon Machine Images"
)

var amiLogger *log.Entry

func newAmiManager(cfg *config.Config, logPath string) resource.Manager {
	amiLogger = config.GetLogger(amiName, logPath)

	amiManager = resource.Manager{
		Name:     amiName,
		LongName: amiLongName,
		Config:   cfg,
		Logger:   logger,
		GetAll: func() ([]*resource.Resource, error) {
			return GetAllImages(*cfg.Aws)
		},
		Destroy: func(resources []*resource.Resource) error {
			return TerminateImages(*cfg.Aws, resources)
		},
	}
	return amiManager
}
