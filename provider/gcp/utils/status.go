package utils

import (
	"fmt"

	"github.com/mensaah/reka/resource"
	"github.com/mensaah/reka/state"
)

// GetResourceStatus Get the current status of Resource: Pending, Running, ... Stopped
func GetComputeInstanceStatus(s string) resource.Status {
	switch s {
	case "PROVISIONING", "REPARING", "STAGING":
		return resource.Pending
	case "RUNNING":
		return resource.Running
	case "DEPROVISIONING":
		return resource.ShuttingDown
	case "STOPPING", "SUSPENDING":
		return resource.Stopping
	case "STOPPED", "SUSPENDED", "TERMINATED":
		return resource.Stopped
	case "ERROR", "DEGRADED":
		return resource.Error
	default:
		return resource.Destroyed
	}
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
