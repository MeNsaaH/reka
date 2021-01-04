package state

import (
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
)

// Backender interface
// Definitions of set of methods required for all state types
type Backender interface {
	Writer
	Reader
}

// Reader is the interface for state types that can return state resources
type Reader interface {
	GetState() *State
}

// Writer is the interface for state types that can write to state files
type Writer interface {
	WriteState(st *State) error
}

// InitBackend initializes the backend for reading and writing to state
// It checks state configuration defined in the config and returns
// The appropriate backend
func InitBackend() Backender {
	cfg = config.GetConfig()
	switch cfg.StateBackend.Type {
	default:
		log.Debugf("using Local State at %s", cfg.StateBackend.Path)
		s := NewEmptyState()
		backend = LocalBackend{
			Path:  cfg.StateBackend.Path,
			state: &s,
		}
	}
	return backend
}
