package state

import (
	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider"
)

var backend Backender
var cfg *config.Config

// ProvidersState represents state for providers
//   {
//    aws: {
// 		s3: [ Resource1, Resource2 ]
//    }
//  }
type ProvidersState map[string]provider.Resources

// State object: Represents a reka state consisting of current and desired state
// Desired state is used for resumption of resources. It stores the attributes of the resource to resume to
// for resources that need extra info like size of node pool to resize to
// current: {
//    aws: {
// 		s3: [ Resource1, Resource2 ]
//    }
// }
// desired: {
//    aws: {
// 		s3: [ Resource1, Resource2 ]
//    }
// }
// State object stores state for both current and desired states
// of resources
type State struct {
	Current ProvidersState
	Desired ProvidersState
}

func (s *State) diff() State {
	return Diff(s.Desired, s.Current)
}

// Empty checks if state is empty
func (s State) Empty() bool {
	return len(s.Current) == 0 && len(s.Desired) == 0
}

// Diff returns the difference in between two states
func Diff(previous, current ProvidersState) State {
	return State{}
}

// NewEmptyState gets a new empty state
func NewEmptyState() State {
	return State{Desired: make(ProvidersState), Current: make(ProvidersState)}
}
