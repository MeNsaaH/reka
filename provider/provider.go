package provider

import (
	log "github.com/sirupsen/logrus"
)

// Resources Contains map where key is the name of the manager and the value is an array of the
// resources managed by the manager
type Resources map[string][]*Resource

// Provider : Provider definition
// Implements all logic for controlling Resource Managers
type Provider struct {
	Name             string
	ResourceManagers map[string]ResourceManager // An array `ResourceManagers` interfaces
	Config           *Config
}

// GetAllResources : Returns all resources which reka can find
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

// GetDestroyableResources : Return the resources which can be destroyed
func (p *Provider) GetDestroyableResources(resources Resources, config Config) Resources {
	destroyableResources := make(Resources)
	for mgrName, resList := range resources {
		var destroyableResList []*Resource
		for _, resource := range resList {
			log.Info(resource, mgrName)
			if ShouldInitiateDestruction(resource.Tags) {
				destroyableResList = append(destroyableResList, resource)
			}
		}
		destroyableResources[mgrName] = destroyableResList
	}
	return destroyableResources
}

// GetStoppableResources : Return the resources which can be stopped
func (p *Provider) GetStoppableResources(resources Resources, config Config) Resources {
	stoppableResources := make(Resources)
	for mgrName, resList := range resources {
		var stoppableResList []*Resource
		for _, resource := range resList {
			if resource.IsActive() && ShouldInitiateStopping(resource.Tags) {
				stoppableResList = append(stoppableResList, resource)
			}
		}
		stoppableResources[mgrName] = stoppableResList
	}
	return stoppableResources
}

// GetResumableResources : Return the resources which can be Resumed
func (p *Provider) GetResumableResources(resources Resources, config Config) Resources {
	resumableResource := make(Resources)
	for mgrName, resList := range resources {
		var resumableResList []*Resource
		for _, resource := range resList {
			log.Info(resource, mgrName)
			if resource.IsStopped() && ShouldInitiateResumption(resource.Tags) {
				resumableResList = append(resumableResList, resource)
			}
		}
		resumableResource[mgrName] = resumableResList
	}
	return resumableResource
}

// DestroyResources : Return the resources which can be destroyed
func (p *Provider) DestroyResources(resources Resources) map[string]string {
	errs := make(map[string]string)

	for mgrName, res := range resources {
		mgr := p.getManager(mgrName)
		log.Infof("Destroying %s ", mgrName)
		if err := mgr.Destroy(res); err != nil {
			errs[mgrName] = err.Error()
		}
	}
	return errs
}

// StopResources : Return the resources which can be destroyed
func (p *Provider) StopResources(resources Resources) map[string]string {
	errs := make(map[string]string)

	for mgrName, res := range resources {
		mgr := p.getManager(mgrName)
		if _, ok := mgr.(ResourceStopperResumer); ok {
			log.Infof("Stopping %s ", mgrName)
			if err := mgr.Stop(res); err != nil {
				errs[mgrName] = err.Error()
			}
		}
	}
	return errs
}

// ResumeResources : Return the resources which can be destroyed
func (p *Provider) ResumeResources(resources Resources) map[string]string {
	errs := make(map[string]string)

	for mgrName, res := range resources {
		mgr := p.getManager(mgrName)
		if _, ok := mgr.(ResourceStopperResumer); ok {
			log.Infof("Resuming %s ", mgrName)
			if err := mgr.Resume(res); err != nil {
				errs[mgrName] = err.Error()
			}
		}
	}
	return errs
}

func (p *Provider) getManager(name string) ResourceManager {
	return p.ResourceManagers[name]
}

// Nuke : POOF !!!
// destroys everything tracked by reka
func (p *Provider) Nuke() (string, error) {
	return "", nil
}
