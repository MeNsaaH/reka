package gcp

import (
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
)

const (
	// Name of resource
	cloudStorageName = "CloudStorage"
	// LongName descriptive name for resource
	cloudStorageLongName = "Simple Storage Service"
)

var cloudStorageLogger *log.Entry

func newCloudStorageManager(cfg *config.Config, logPath string) resource.Manager {

	cloudStorageLogger = config.GetLogger(cloudStorageName, logPath)

	return resource.Manager{
		Name:     cloudStorageName,
		LongName: cloudStorageLongName,
		Config:   cfg,
		Logger:   logger,
		GetAll: func() ([]*resource.Resource, error) {
			return getAllBuckets(&cfg.Gcp)
		},
		Destroy: func(resources []*resource.Resource) error {
			return destroyBuckets(&cfg.Gcp, resources)
		},
	}
}
