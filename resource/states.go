//go:generate enumer -type=State ./provider/states.go
package resource

// State :The Current State of the Resource
type State int

const (
	Pending State = iota
	Running
	ShuttingDown
	Destroyed
	Stopping
	Stopped
)

// StyleClass : The css class to represent the state with
func (s State) StyleClass() string {
	if s == Running {
		return "success"
	} else if s == Pending || s == ShuttingDown || s == Stopping {
		return "info"
	} else {
		return "danger"
	}
}
