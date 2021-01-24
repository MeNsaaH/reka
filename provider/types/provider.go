package types

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
	"github.com/mensaah/reka/rules"
)

// Resources Contains map where key is the name of the manager and the value is an array of the
// resources managed by the manager
type Resources map[string][]*resource.Resource

type SafeResources struct {
	mu  sync.Mutex
	v   Resources
	err map[string]error
}

// Provider : Provider definition
// Implements all logic for controlling Resource Managers
type Provider struct {
	Name     string
	Logger   *log.Entry
	LogPath  string
	Managers map[string]*resource.Manager // [mgrName: Manager]
}

// SetLogger : Sets Logger properties for Provider
func (p *Provider) SetLogger(id string) {
	logFile := fmt.Sprintf("%s/%s.log", config.GetConfig().LogPath, id)
	// Setup Logger
	p.LogPath = logFile
	p.Logger = config.GetLogger(p.Name, logFile)
}

// GetAllResources : Returns all resources which reka can find
func (p *Provider) GetAllResources() Resources {
	p.Logger.Info("Fetching All Resources")
	var wg sync.WaitGroup
	resources := SafeResources{v: make(Resources)}
	for _, resMgr := range p.Managers {
		wg.Add(1)
		go func(res *SafeResources, resMgr *resource.Manager) {
			defer wg.Done()

			resMgrResources, err := resMgr.GetAll()
			if err != nil {
				p.Logger.Error(err)
			}
			res.mu.Lock()
			defer res.mu.Unlock()
			res.v[resMgr.Name] = resMgrResources
		}(&resources, resMgr)
	}
	wg.Wait()
	return resources.v
}

// GetDestroyableResources : Return the resources which can be destroyed
func (p *Provider) GetDestroyableResources(resources Resources) Resources {
	p.Logger.Debug("Getting Destroyable Resources")
	count := 0
	destroyableResources := make(Resources)
	for mgrName, resList := range resources {
		var destroyableResList []*resource.Resource
		for _, r := range resList {
			// Returns the first Matching Rule Action for a resource
			for _, rule := range rules.GetRules() {
				if action := rule.CheckResource(r); action == rules.Destroy {
					destroyableResList = append(destroyableResList, r)
				}
			}
		}
		count += len(destroyableResList)
		destroyableResources[mgrName] = destroyableResList
	}
	p.Logger.Infof("Found %d resources to be destroyed", count)
	return destroyableResources
}

// GetStoppableResources : Return the resources which can be stopped
func (p *Provider) GetStoppableResources(resources Resources) Resources {
	p.Logger.Debug("Getting Stoppable Resources")
	stoppableResources := make(Resources)
	count := 0
	for mgrName, resList := range resources {
		var stoppableResList []*resource.Resource
		for _, r := range resList {
			for _, rule := range rules.GetRules() {
				if action := rule.CheckResource(r); r.IsActive() && action == rules.Stop {
					stoppableResList = append(stoppableResList, r)
				}
			}
		}
		if p.getManager(mgrName).Stop != nil && p.getManager(mgrName).Resume != nil {
			count += len(stoppableResList)
			stoppableResources[mgrName] = stoppableResList
		}
	}
	p.Logger.Infof("Found %d resources to be stopped", count)
	return stoppableResources
}

// GetResumableResources : Return the resources which can be Resumed
func (p *Provider) GetResumableResources(resources Resources) Resources {
	p.Logger.Debug("Getting resumable Resources")
	resumableResource := make(Resources)
	count := 0
	for mgrName, resList := range resources {
		var resumableResList []*resource.Resource
		for _, r := range resList {
			for _, rule := range rules.GetRules() {
				if action := rule.CheckResource(r); r.IsStopped() && action == rules.Resume {
					resumableResList = append(resumableResList, r)
				}
			}
		}
		if p.getManager(mgrName).Stop != nil && p.getManager(mgrName).Resume != nil {
			resumableResource[mgrName] = resumableResList
			count += len(resumableResList)
		}
	}
	p.Logger.Infof("Found %d resources to be resumed", count)
	return resumableResource
}

// GetUnusedResources : Return the resources which can are not currently in use and can be destroyed
func (p *Provider) GetUnusedResources(resources Resources) Resources {
	unusedResources := make(Resources)
	for mgrName, resList := range resources {
		var unusedResList []*resource.Resource
		for _, r := range resList {
			if r.IsUnused() {
				unusedResList = append(unusedResList, r)
			}
		}
		unusedResources[mgrName] = unusedResList
	}
	return unusedResources
}

// DestroyResources : Return the resources which can be destroyed
func (p *Provider) DestroyResources(resources Resources) map[string]error {
	errs := make(map[string]error)
	p.Logger.Info("Destroying Resources...")
	var wg sync.WaitGroup

	for mgrName, res := range resources {
		if len(res) > 0 {
			wg.Add(1)
			go func(mgrName string, res []*resource.Resource) {
				defer wg.Done()
				mgr := p.getManager(mgrName)
				p.Logger.Debugf("Destroying %s ", mgrName)
				if err := mgr.Destroy(res); err != nil {
					errs[mgrName] = err
				}
			}(mgrName, res)
		}
	}
	wg.Wait()
	return errs
}

// StopResources : Return the resources which can be destroyed
func (p *Provider) StopResources(resources Resources) map[string]error {
	errs := make(map[string]error)
	p.Logger.Info("Stopping Resources...")
	var wg sync.WaitGroup

	for mgrName, res := range resources {
		if len(res) > 0 {
			wg.Add(1)

			go func(mgrName string, res []*resource.Resource) {
				defer wg.Done()
				mgr := p.getManager(mgrName)
				if mgr.Stop != nil {
					p.Logger.Debugf("Stopping %d %s Resources", len(res), mgrName)
					if err := mgr.Stop(res); err != nil {
						errs[mgrName] = err
					}
				}
			}(mgrName, res)
		}
	}
	wg.Wait()
	return errs
}

// ResumeResources : Return the resources which can be destroyed
func (p *Provider) ResumeResources(resources Resources) map[string]error {
	var wg sync.WaitGroup
	errs := make(map[string]error)
	p.Logger.Info("Resuming Resources...")

	for mgrName, res := range resources {
		if len(res) > 0 {
			wg.Add(1)
			go func(mgrName string, res []*resource.Resource) {
				defer wg.Done()
				mgr := p.getManager(mgrName)
				if mgr.Resume != nil {
					p.Logger.Debugf("Resuming %d %s Resources", len(res), mgrName)
					if err := mgr.Resume(res); err != nil {
						errs[mgrName] = err
					}
				}
			}(mgrName, res)
		}
	}
	wg.Wait()
	return errs
}

func (p *Provider) getManager(name string) *resource.Manager {
	return p.Managers[name]
}

// Nuke : POOF !!!
// destroys everything tracked by reka
func (p *Provider) Nuke() (string, error) {
	return "", nil
}

// GetResourceNames Get a array of resource names
func (p *Provider) GetResourceNames() []string {
	var arr []string
	for resMgr := range p.Managers {
		arr = append(arr, resMgr)
	}
	return arr
}
