package types

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Resource : The Provider Interface
// fields with `gorm:"-"` are ignored
type Resource struct {
	ID       string
	Name     string
	LongName string
	ImageURL string // A link to an image representing the resource

	Provider string // The provider it belongs to GCP, AWS ...
	Region   string // Region of Resource

	// The current state of the instance; stopped, running, pending
	State        State
	CreationDate time.Time

	// Error thrown during Fetching resource related data
	FetchError error `gorm:"-"`
	// Error thrown during Destroying the resource
	DestroyError error `gorm:"-"`
	// Error thrown when stopping/hibernating/pausing/shuttingdown the instance
	StopError error `gorm:"-"`
	// Error thrown when resuming the instance
	ResumeError error `gorm:"-"`

	// Tags are for AWS Instances
	Tags ResourceTags `gorm:"-"`
}

func (r Resource) String() string {
	return fmt.Sprintf("<%s:%s>", r.Name, r.ID)
}

// IsActive return if resource is currently running
func (r Resource) IsActive() bool {
	return r.State == Running
}

// IsStopped return whether resource is currently stopped not destroyed
func (r Resource) IsStopped() bool {
	return r.State == Stopped
}

// StopperResumer : Interface implemented by resources that can be stopped and resumed
type StopperResumer interface {
	Stop([]*Resource) error
	Resume([]*Resource) error
}

// ResourceManager : Implements all methods to manage resource state
type ResourceManager interface {
	Destroy([]*Resource) error
	GetAll() ([]*Resource, error)
	GetReapable() ([]*Resource, error)
	GetName() string
	GetLogger() *log.Entry
}

// DefaultResourceManager : Base Resource Manager to be embedded in other structs
type DefaultResourceManager struct {
	Name     string // Short Name of the Resource Manager
	LongName string // A More Elaborate name for the manager
	Provider string // The supported Provider
	Logger   *log.Entry
}

//// Default Functions for Resource Managers
//func (mgr *DefaultResourceManager) Stop([]*Resource) error {
//  mgr.Logger.Debug("Stop not supported for this Resource")
//  return nil
//}

//func (mgr *DefaultResourceManager) Resume([]*Resource) error {
//  mgr.Logger.Debug("`Resume` not supported for this Resource")
//  return nil
//}

//func (mgr *DefaultResourceManager) Destroy([]*Resource) error {
//  mgr.Logger.Debug("`Destroy` not implemented for this Resource")
//  return nil
//}
