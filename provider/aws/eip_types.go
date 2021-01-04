package aws

import (
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
)

// Manages Eip instances on the AWS.
var eipManager resource.Manager

const (
	// Name of resource
	eipName = "EIP"
	// LongName descriptive name for resource
	eipLongName = "Elastic IP"
)

var eipLogger *log.Entry

func newEipManager(cfg *config.Config, logPath string) resource.Manager {
	eipLogger = config.GetLogger(eipName, logPath)

	eipManager = resource.Manager{
		Name:     eipName,
		LongName: eipLongName,
		Config:   cfg,
		Logger:   logger,
		GetAll: func() ([]*resource.Resource, error) {
			return GetAllIPAddresses(*cfg.Aws)
		},
		Destroy: func(resources []*resource.Resource) error {
			return TerminateIPAddresses(*cfg.Aws, resources)
		},
	}
	return eipManager
}
