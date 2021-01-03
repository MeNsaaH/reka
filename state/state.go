package state

import (
	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider"
)

var backend Backender
var cfg *config.Config

// State object: Represents a reka state consisting of resources
// aws: {
// 		s3: [ Resource1, Resource2 ]
// }
type State map[string]provider.Resources

func (s *State) diff(st *State) State {
	return Diff(s, st)
}

func (s State) empty() bool {
	return len(s) == 0
}

// Diff returns the difference in between two states
func Diff(previous, current *State) State {
	return State{}
}
