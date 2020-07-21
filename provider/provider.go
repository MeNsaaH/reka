package provider

import (
	log "github.com/sirupsen/logrus"
)

// ResourcesWithManager Contains map where key is the name of the manager and the value is an array of the
// resources managed by the manager
type Resources map[string][]*Resource

type Provider struct {
	Name             string
	ResourceManagers map[string]ResourceManager // An array `ResourceManagers` interfaces
}

func (p *Provider) GetAllResources() Resources {
	resources := make(Resources)
	for _, resMgr := range p.ResourceManagers {
		resMgrResources, err := resMgr.GetAll()
		if err != nil {
			log.Error(err)
		}
		resources[resMgr.GetName()] = resMgrResources
	}
	return resources
}

// Auth returns the Config object to be used for authentication of requests
func (p *Provider) Auth(config Config) interface{} {
	return nil
}

// GetReapableResources : Return the resources which can be destroyed
func (p *Provider) GetReapableResources(resources *Resources, config Config) Resources {
	return Resources{}
}

// DestroyResources : Return the resources which can be destroyed
func (p *Provider) DestroyResources(resources Resources) map[string]string {
	errs := make(map[string]string)

	for mgrName, res := range resources {
		mgr := p.getManager(mgrName)
		if err := mgr.Destroy(res); err != nil {
			errs[mgrName] = err.Error()
		}
	}
	return errs
}

func (p *Provider) getManager(name string) ResourceManager {
	return p.ResourceManagers[name]
}

func (p *Provider) Nuke() (string, error) {
	return "", nil
}
