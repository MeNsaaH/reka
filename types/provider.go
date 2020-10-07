package types

import (
	"sync"

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
}

// GetAllResources : Returns all resources which reka can find
func (p *Provider) GetAllResources() Resources {
	log.Info("Fetching All Resources")
	var wg sync.WaitGroup
	resources := make(Resources)
	for _, resMgr := range p.ResourceManagers {
		wg.Add(1)
		go func(res Resources, resMgr ResourceManager) {
			defer wg.Done()
			resMgrResources, err := resMgr.GetAll()
			if err != nil {
				resMgr.GetLogger().Error(err)
			}
			res[resMgr.GetName()] = resMgrResources
		}(resources, resMgr)
	}
	wg.Wait()
	return resources
}

// GetDestroyableResources : Return the resources which can be destroyed
func (p *Provider) GetDestroyableResources(resources Resources) Resources {
	log.Debug("Getting Destroyable Resources")
	destroyableResources := make(Resources)
	for mgrName, resList := range resources {
		var destroyableResList []*Resource
		for _, resource := range resList {
			if ShouldInitiateDestruction(resource.Tags) {
				destroyableResList = append(destroyableResList, resource)
			}
		}
		destroyableResources[mgrName] = destroyableResList
	}
	return destroyableResources
}

// GetStoppableResources : Return the resources which can be stopped
func (p *Provider) GetStoppableResources(resources Resources) Resources {
	log.Debug("Getting Stoppable Resources")
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
func (p *Provider) GetResumableResources(resources Resources) Resources {
	log.Debug("Getting resumable Resources")
	resumableResource := make(Resources)
	for mgrName, resList := range resources {
		var resumableResList []*Resource
		for _, resource := range resList {
			if resource.IsStopped() && ShouldInitiateResumption(resource.Tags) {
				resumableResList = append(resumableResList, resource)
			}
		}
		resumableResource[mgrName] = resumableResList
	}
	return resumableResource
}

// GetUnusedResources : Return the resources which can are not currently in use and can be destroyed
func (p *Provider) GetUnusedResources(resources Resources) Resources {
	return Resources{}
}

// DestroyResources : Return the resources which can be destroyed
func (p *Provider) DestroyResources(resources Resources) map[string]string {
	errs := make(map[string]string)
	log.Debugf("Destroying Resources...")

	for mgrName, res := range resources {
		mgr := p.getManager(mgrName)
		mgr.GetLogger().Debugf("Destroying %s ", mgrName)
		if err := mgr.Destroy(res); err != nil {
			errs[mgrName] = err.Error()
		}
	}
	return errs
}

// StopResources : Return the resources which can be destroyed
func (p *Provider) StopResources(resources Resources) map[string]string {
	errs := make(map[string]string)
	log.Info("Stopping Resources...")
	var wg sync.WaitGroup

	for mgrName, res := range resources {
		wg.Add(1)

		go func(mgrName string, res []*Resource) {
			defer wg.Done()
			mgr := p.getManager(mgrName)
			if m, ok := mgr.(StopperResumer); ok {
				mgr.GetLogger().Debugf("Stopping %d %s Resources", len(res), mgrName)
				if err := m.Stop(res); err != nil {
					errs[mgrName] = err.Error()
				}
			}
		}(mgrName, res)
	}
	wg.Wait()
	return errs
}

// ResumeResources : Return the resources which can be destroyed
func (p *Provider) ResumeResources(resources Resources) map[string]string {
	var wg sync.WaitGroup
	errs := make(map[string]string)
	log.Info("Resuming Resources...")

	for mgrName, res := range resources {
		wg.Add(1)
		go func(mgrName string, res []*Resource) {
			defer wg.Done()
			mgr := p.getManager(mgrName)
			if m, ok := mgr.(StopperResumer); ok {
				mgr.GetLogger().Debugf("Resuming %d %s Resources", len(res), mgrName)
				if err := m.Resume(res); err != nil {
					errs[mgrName] = err.Error()
				}
			}
		}(mgrName, res)
	}
	wg.Wait()
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

// GetResourceNames Get a array of resource names
func (p *Provider) GetResourceNames() []string {
	var arr []string
	for resMgr := range p.ResourceManagers {
		arr = append(arr, resMgr)
	}
	return arr
}
