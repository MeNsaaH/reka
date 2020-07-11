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

type ResourceInterface interface {
	String() string
	isActive() bool
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
