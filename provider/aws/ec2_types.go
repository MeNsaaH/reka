package aws

import (
	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/types"
)

var ec2Manager types.ResourceManager

const (
	// Name of resource
	ec2Name = "EC2"
	// LongName descriptive name for resource
	ec2LongName = "Elastic Compute Cloud"
)

func newEC2Manager(cfg *config.Config, logPath string) types.ResourceManager {
	logger := types.GetLogger(ec2Name, logPath)

	ec2Manager = types.ResourceManager{
		Name:     ec2Name,
		LongName: ec2LongName,
		Config:   cfg,
		Logger:   logger,
		GetAll: func() ([]*types.Resource, error) {
			return GetAllEC2Instances(*cfg.Aws, logger)
		},
		Destroy: func(resources []*types.Resource) error {
			return TerminateEC2Instances(*cfg.Aws, resources, logger)
		},
		Stop: func(resources []*types.Resource) error {
			return StopEC2Instances(*cfg.Aws, resources, logger)
		},
		Resume: func(resources []*types.Resource) error {
			return ResumeEC2Instances(*cfg.Aws, resources, logger)
		},
	}
	return ec2Manager
}
