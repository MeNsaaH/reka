package provider

import (
	"fmt"
	"time"
)

// State :The Current State of the Resource
type State int

const (
	Pending State = iota
	Running
	ShuttingDown
	Terminated
	Stopping
	Stopped
)

type ResourceDestroyer interface {
	Destroy([]*Resource) error
}

type ResourceStopResumer interface {
	Stop([]*Resource) error
	Resume([]*Resource) error
}

type ResourceManager interface {
	ResourceDestroyer
	ResourceStopResumer
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
	return fmt.Sprintf("%s - %s", r.Name, r.ID)
}

func (r Resource) IsActive() bool {
	return r.State == Running
}

func (r Resource) IsStopped() bool {
	return r.State == Stopped
}
