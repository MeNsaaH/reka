package gcp

import (
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider/types"
	"github.com/mensaah/reka/resource"
)

var (
	providerName     = "gcp"
	logger           *log.Entry
	resourceManagers map[string]*resource.Manager
)

func GetName() string {
	return providerName
}

// NewResource Returns a new Resource object
func NewResource(id, manager string) *resource.Resource {
	resource := resource.Resource{}
	resource.UUID = id
	resource.Manager = resourceManagers[manager]
	resource.ProviderName = providerName

	return &resource
}

// NewProvider : Creates a New AWS Provider
func NewProvider() (*types.Provider, error) {

	gcp := types.Provider{}
	gcp.Name = providerName

	gcp.SetLogger("logger.log")

	cfg := config.GetConfig()

	cloudStorageManager := newCloudStorageManager(cfg, gcp.LogPath)
	computeInstanceManager := newComputeInstanceManager(cfg, gcp.LogPath)
	gkeManager := newGkeManager(cfg, gcp.LogPath)

	resourceManagers = map[string]*resource.Manager{
		cloudStorageManager.Name:    &cloudStorageManager,
		computeInstanceManager.Name: &computeInstanceManager,
		gkeManager.Name:             &gkeManager,
	}

	gcp.Managers = resourceManagers
	return &gcp, nil
}
