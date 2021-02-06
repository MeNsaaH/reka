package gcp

import (
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
)

const (
	// Name of resource
	computeInstanceName = "compute"
	// LongName descriptive name for resource
	computeLongName = "Compute Engine"
)

var computeLogger *log.Entry

func newComputeInstanceManager(cfg *config.Config, logPath string) resource.Manager {

	computeLogger = config.GetLogger(computeInstanceName, logPath)

	return resource.Manager{
		Name:     computeInstanceName,
		LongName: computeLongName,
		Config:   cfg,
		Logger:   computeLogger,
		GetAll: func() ([]*resource.Resource, error) {
			return getAllComputeInstances(cfg.Gcp)
		},
		Destroy: func(resources []*resource.Resource) error {
			return destroyComputeInstances(cfg.Gcp, resources)
		},
		Stop: func(resources []*resource.Resource) error {
			return stopComputeInstances(cfg.Gcp, resources)
		},
		Resume: func(resources []*resource.Resource) error {
			return startComputeInstances(cfg.Gcp, resources)
		},
	}
}
