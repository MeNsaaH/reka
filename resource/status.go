//go:generate enumer -type=Status ./status.go
package resource

// Status :The State of the Resource
type Status int

const (
	Pending Status = iota
	Running
	ShuttingDown
	Destroyed
	Stopping
	Stopped
	Unused
	Error
)

// StyleClass : The css class to represent the state with
func (s Status) StyleClass() string {
	if s == Running {
		return "success"
	} else if s == Pending || s == ShuttingDown || s == Stopping {
		return "info"
	} else {
		return "danger"
	}
}
