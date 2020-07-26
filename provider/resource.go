package provider

import (
	"fmt"
	"time"
)

// ResourceDestroyer implements methods to completely Destroy a resource
type ResourceDestroyer interface {
	Destroy([]*Resource) error
}

// ResourceStopperResumer implements methods to stop and resume a resource when condition is met
// Resources that implement this interface can be stopped and resumed e.g VMs
type ResourceStopperResumer interface {
	Stop([]*Resource) error
	Resume([]*Resource) error
}

// ResourceManager : Implements all methods to manage resource state
type ResourceManager interface {
	ResourceDestroyer
	ResourceStopperResumer
	GetAll() ([]*Resource, error)
	GetReapable(Config) ([]*Resource, error)
	GetName() string
}

// Resource : The Provider Interface
type Resource struct {
	ID       string
	Name     string
	LongName string

	// The current state of the instance; stopped, running, pending
	State        State
	CreationDate time.Time

	// Resources that need to be deleted or destroyed before this instance can be destroyed
	Dependents []Resource
	// Error thrown during Fetching resource related data
	FetchError error
	// Error thrown during Destroying the resource
	DestroyError error
	// Error thrown when stopping/hibernating/pausing/shuttingdown the instance
	StopError error

	Region string       // Region of Resource for AWS Instances
	Tags   ResourceTags // Tags are for AWS Instances
}

func (r Resource) String() string {
	return fmt.Sprintf("<%s:%s>", r.Name, r.ID)
}

func (r Resource) IsActive() bool {
	return r.State == Running
}

func (r Resource) IsStopped() bool {
	return r.State == Stopped
}
