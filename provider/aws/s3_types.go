package aws

import (
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
)

const (
	// Name of resource
	s3Name = "S3"
	// LongName descriptive name for resource
	s3LongName = "Simple Storage Service"
)

var s3Logger *log.Entry

func newS3Manager(cfg *config.Config, logPath string) resource.Manager {

	s3Logger = config.GetLogger(s3Name, logPath)

	return resource.Manager{
		Name:     s3Name,
		LongName: s3LongName,
		Config:   cfg,
		Logger:   logger,
		GetAll: func() ([]*resource.Resource, error) {
			return getAllS3Buckets(*cfg.Aws)
		},
		Destroy: func(resources []*resource.Resource) error {
			return destroyS3Buckets(*cfg.Aws, resources)
		},
	}
}
