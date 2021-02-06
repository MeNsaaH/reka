package resource

import (
	"fmt"
	"strings"
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

func (mgr Manager) IsStoppable() bool {
	return mgr.Stop != nil && mgr.Resume != nil
}

// Resource : The Provider Interface
// fields with `gorm:"-"` are ignored in database columns
// TODO create dependent Resource field. And ensure dependent resources are first destroyed before
// destruction of resource happens
type Resource struct {
	// Add ID, CreatedAt, UpdatedAt and DeletedAt fields
	gorm.Model `json:"-"`
	// UUID defines any unique field use to identify a resource. For some resources its their IDs, some their names.
	// Not using ID because gorm.Model defines an ID field already
	UUID         string   `gorm:"unique;not null"`
	Manager      *Manager `gorm:"foreignKey:ManagerName;references:Name" json:"-"`
	ProviderName string

	Region string // Region/Location of Resource
	Zone   string // Zone of Resource

	// The state of the instance; stopped, running, pending
	Status Status
	// The time the instance was created on the Provider
	CreationDate time.Time

	// Attributes that the resource possess e.g. EKS Clusters have nodegroups
	Attributes map[string]interface{}

	// Tags are for AWS Instances
	Tags Tags `gorm:"-"`

	// SubResources are other resources that are part of the current resource but cannot stand as
	// an independent resource e.g Nodegroups are subresource of EKS Clusters. Destroying an EKS Cluster
	// destroys all nodegroups associated with it
	SubResources map[string][]*Resource
}

func (r Resource) String() string {
	return fmt.Sprintf("<%s:%s>", r.Manager, r.UUID)
}

// IsActive return if resource is currently running
func (r Resource) IsActive() bool {
	return r.Status == Running
}

// IsStopped return whether resource is currently stopped not destroyed
func (r Resource) IsStopped() bool {
	return r.Status == Stopped
}

func (r Resource) IsUnused() bool {
	return r.Status == Unused
}

// Uri a simple uri of the resource in the form provider.resource_type for example ec2 instances will have the
// url aws.ec2
func (r Resource) Uri() string {
	return fmt.Sprintf("%s.%s", strings.ToLower(r.ProviderName), strings.ToLower(r.Manager.Name))

}
