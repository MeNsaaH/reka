package types

import (
	"fmt"
	"time"

	"github.com/mensaah/reka/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ResourceManager : S3, CloudStorage etc
type ResourceManager struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	LongName string
	ImageURL string         // A link to an image/logo representing the resource
	Config   *config.Config `gorm:"-"`

	Logger *log.Entry `gorm:"-"`

	// Methods Implemented By Resource Manager
	GetAll func() ([]*Resource, error) `gorm:"-"` // Required

	Destroy func([]*Resource) error `gorm:"-"` // Required
	Stop    func([]*Resource) error `gorm:"-"`
	Resume  func([]*Resource) error `gorm:"-"`
}

func (mgr ResourceManager) String() string {
	return mgr.Name
}

// Resource : The Provider Interface
// fields with `gorm:"-"` are ignored in database columns
type Resource struct {
	gorm.Model
	UUID        string `gorm:"unique;not null"`
	ManagerName string
	Manager     *ResourceManager `gorm:"foreignKey:ManagerName;references:Name"`

	Region string // Region of Resource

	// The current state of the instance; stopped, running, pending
	State State
	// The time the instance was created on the Provider
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
	return fmt.Sprintf("<%s:%s>", r.Manager, r.UUID)
}

// IsActive return if resource is currently running
func (r Resource) IsActive() bool {
	return r.State == Running
}

// IsStopped return whether resource is currently stopped not destroyed
func (r Resource) IsStopped() bool {
	return r.State == Stopped
}
