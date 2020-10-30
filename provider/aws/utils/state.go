package utils

import (
	"github.com/mensaah/reka/types"
)

// GetResourceState Get the current status of Resource: Pending, Running, ... Stopped
func GetResourceState(s int32) types.State {
	switch s {
	case 0:
		return types.Pending
	case 16:
		return types.Running
	case 32:
		return types.ShuttingDown
	case 48:
		return types.Destroyed
	case 64:
		return types.Stopping
	case 80:
		return types.Stopped
	default:
		return types.Stopped
	}
}
