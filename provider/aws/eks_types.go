package aws

import (
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
)

// Manages EKS instances on the AWS.
// EKS resources support stopping/resuming and terminating instances.

var eksManager resource.Manager

const (
	// Name of resource
	eksName = "eks"
	// LongName descriptive name for resource
	eksLongName = "Elastic Compute Cloud"

	nodegroupName = "Nodegroup"
)

var eksLogger *log.Entry

func newEksManager(cfg *config.Config, logPath string) resource.Manager {
	eksLogger = config.GetLogger(eksName, logPath)

	eksManager = resource.Manager{
		Name:     eksName,
		LongName: eksLongName,
		Config:   cfg,
		Logger:   logger,
		GetAll: func() ([]*resource.Resource, error) {
			return GetAllEKSClusters(*cfg.Aws)
		},
		Destroy: func(resources []*resource.Resource) error {
			return TerminateEKSClusters(*cfg.Aws, resources)
		},
		Stop: func(resources []*resource.Resource) error {
			return StopEKSClusters(*cfg.Aws, resources)
		},
		Resume: func(resources []*resource.Resource) error {
			return ResumeEKSClusters(*cfg.Aws, resources)
		},
	}
	return eksManager
}
