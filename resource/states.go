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
