package resource

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/mensaah/reka/config"
)

// Manager : S3, CloudStorage etc
type Manager struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	LongName string
	ImageURL string         // A link to an image/logo representing the resource
	Config   *config.Config `gorm:"-"`

	Logger *log.Entry `gorm:"-"`

	// Methods Implemented By Resource Manager
	GetAll func() ([]*Resource, error) `gorm:"-" json:"-"` // Required

	Destroy func([]*Resource) error `gorm:"-" json:"-"` // Required
	Stop    func([]*Resource) error `gorm:"-" json:"-"`
	Resume  func([]*Resource) error `gorm:"-" json:"-"`
}

func (mgr Manager) String() string {
	return mgr.Name
}

// Resource : The Provider Interface
// fields with `gorm:"-"` are ignored in database columns
type Resource struct {
	gorm.Model
	UUID        string `gorm:"unique;not null"`
	ManagerName string
	Manager     *Manager `gorm:"foreignKey:ManagerName;references:Name"`

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
	Tags Tags `gorm:"-"`
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
