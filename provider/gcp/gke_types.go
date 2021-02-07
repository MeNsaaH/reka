package gcp

import (
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
)

const (
	// Name of resource
	gkeName      = "gke"
	nodePoolName = "nodepool"
	// LongName descriptive name for resource
	gkeLongName = "Compute Engine"
)

var gkeLogger *log.Entry

func newGkeManager(cfg *config.Config, logPath string) resource.Manager {

	gkeLogger = config.GetLogger(gkeName, logPath)

	return resource.Manager{
		Name:     gkeName,
		LongName: gkeLongName,
		Config:   cfg,
		Logger:   gkeLogger,
		GetAll: func() ([]*resource.Resource, error) {
			return getAllGkeClusters(cfg.Gcp)
		},
		Destroy: func(resources []*resource.Resource) error {
			return destroyGkeClusters(cfg.Gcp, resources)
		},
		Stop: func(resources []*resource.Resource) error {
			return stopGkeClusters(cfg.Gcp, resources)
		},
		Resume: func(resources []*resource.Resource) error {
			return startGkeClusters(cfg.Gcp, resources)
		},
	}
}
