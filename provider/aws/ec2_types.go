package aws

import (
	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
)

var ec2Manager resource.Manager

const (
	// Name of resource
	ec2Name = "EC2"
	// LongName descriptive name for resource
	ec2LongName = "Elastic Compute Cloud"
)

func newEC2Manager(cfg *config.Config, logPath string) resource.Manager {
	logger := config.GetLogger(ec2Name, logPath)

	ec2Manager = resource.Manager{
		Name:     ec2Name,
		LongName: ec2LongName,
		Config:   cfg,
		Logger:   logger,
		GetAll: func() ([]*resource.Resource, error) {
			return GetAllEC2Instances(*cfg.Aws, logger)
		},
		Destroy: func(resources []*resource.Resource) error {
			return TerminateEC2Instances(*cfg.Aws, resources, logger)
		},
		Stop: func(resources []*resource.Resource) error {
			return StopEC2Instances(*cfg.Aws, resources, logger)
		},
		Resume: func(resources []*resource.Resource) error {
			return ResumeEC2Instances(*cfg.Aws, resources, logger)
		},
	}
	return ec2Manager
}
