package utils

import (
	"github.com/mensaah/reka/resource"
)

// GetResourceState Get the current status of Resource: Pending, Running, ... Stopped
func GetResourceState(s int32) resource.State {
	switch s {
	case 0:
		return resource.Pending
	case 16:
		return resource.Running
	case 32:
		return resource.ShuttingDown
	case 48:
		return resource.Destroyed
	case 64:
		return resource.Stopping
	case 80:
		return resource.Stopped
	default:
		return resource.Stopped
	}
}
