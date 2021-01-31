package aws

import (
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
)

// Manages RDS instances on the AWS.
// RDS resources support stopping/resuming and terminating instances.

var rdsManager resource.Manager

const (
	// Name of resource
	rdsName = "rds"
	// LongName descriptive name for resource
	rdsLongName = "Relational Database Service"
)

var rdsLogger *log.Entry

func newRDSManager(cfg *config.Config, logPath string) resource.Manager {
	rdsLogger = config.GetLogger(rdsName, logPath)

	rdsManager = resource.Manager{
		Name:     rdsName,
		LongName: rdsLongName,
		Config:   cfg,
		Logger:   logger,
		GetAll: func() ([]*resource.Resource, error) {
			return GetAllRDSInstances(*cfg.Aws)
		},
		Destroy: func(resources []*resource.Resource) error {
			return TerminateRDSInstances(*cfg.Aws, resources)
		},
		Stop: func(resources []*resource.Resource) error {
			return StopRDSInstances(*cfg.Aws, resources)
		},
		Resume: func(resources []*resource.Resource) error {
			return ResumeRDSInstances(*cfg.Aws, resources)
		},
	}
	return rdsManager
}
