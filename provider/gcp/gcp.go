package gcp

import (
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider"
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

	return &resource
}

// NewProvider : Creates a New AWS Provider
func NewProvider() (*provider.Provider, error) {

	gcp := provider.Provider{}
	gcp.Name = providerName

	gcp.SetLogger("logger.log")

	cfg := config.GetConfig()

	cloudStorageManager := newCloudStorageManager(cfg, gcp.LogPath)

	resourceManagers = map[string]*resource.Manager{
		cloudStorageManager.Name: &cloudStorageManager,
	}

	gcp.Managers = resourceManagers
	return &gcp, nil
}
