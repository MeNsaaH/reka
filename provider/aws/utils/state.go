package utils

import (
	"fmt"

	"github.com/mensaah/reka/resource"
	"github.com/mensaah/reka/state"
)

// GetResourceStatus Get the current status of Resource: Pending, Running, ... Stopped
func GetResourceStatus(s int32) resource.Status {
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

// GetEksResourceStatus Get the current status of EKS Resource: Pending, Running, ... Stopped
func GetEksResourceStatus(s string) resource.Status {
	switch s {
	case "CREATING":
	case "UPDATING":
		return resource.Pending
	case "ACTIVE":
		return resource.Running
	case "DELETING":
		return resource.ShuttingDown
	case "FAILED":
	case "CREATE_FAILED":
	case "DELETE_FAILED":
		return resource.Error
	default:
		return resource.Stopped
	}
	return resource.Stopped
}

func GetResourceFromDesiredState(providerName, resMgr, uid string) (*resource.Resource, error) {
	activeState := (state.GetBackend()).GetState()

	for _, w := range activeState.Desired[providerName][resMgr] {
		if w.UUID == uid {
			return w, nil
		}
	}
	return &resource.Resource{}, fmt.Errorf("%s Resource %s not found in state", resMgr, uid)
}
