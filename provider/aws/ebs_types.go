package aws

import (
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
)

// Manages Ebs instances on the AWS.
// Ebs resources support stopping/resuming and terminating instances.

var ebsManager resource.Manager

const (
	// Name of resource
	ebsName = "ebs"
	// LongName descriptive name for resource
	ebsLongName = "Elastic Block Storage"
)

var ebsLogger *log.Entry

func newEbsManager(cfg *config.Config, logPath string) resource.Manager {
	ebsLogger = config.GetLogger(ebsName, logPath)

	ebsManager = resource.Manager{
		Name:     ebsName,
		LongName: ebsLongName,
		Config:   cfg,
		Logger:   logger,
		GetAll: func() ([]*resource.Resource, error) {
			return GetAllEbsVolumes(*cfg.Aws)
		},
		Destroy: func(resources []*resource.Resource) error {
			return TerminateEbsVolumes(*cfg.Aws, resources)
		},
	}
	return ebsManager
}
