package aws

import (
	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/types"
)

const (
	// Name of resource
	s3Name = "S3"
	// LongName descriptive name for resource
	s3LongName = "Simple Storage Service"
)

func newS3Manager(cfg *config.Config, logPath string) types.ResourceManager {

	logger := types.GetLogger(s3Name, logPath)

	return types.ResourceManager{
		Name:     s3Name,
		LongName: s3LongName,
		Config:   cfg,
		Logger:   logger,
		GetAll: func() ([]*types.Resource, error) {
			return getAllS3Buckets(*cfg.Aws, logger)
		},
		Destroy: func(resources []*types.Resource) error {
			return destroyS3Buckets(*cfg.Aws, resources, logger)
		},
	}
}
