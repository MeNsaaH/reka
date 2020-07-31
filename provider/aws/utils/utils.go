package utils

import (
	"github.com/mensaah/reka/provider"
)

// GetResourceState Get the current status of Resource: Pending, Running, ... Stopped
func GetResourceState(s int64) provider.State {
	switch s {
	case 0:
		return provider.Pending
	case 16:
		return provider.Running
	case 32:
		return provider.ShuttingDown
	case 48:
		return provider.Destroyed
	case 64:
		return provider.Stopping
	case 80:
		return provider.Stopped
	default:
		return provider.Stopped
	}
}

type AWSTag struct {
	Key   *string
	Value *string
}

// ParseResourceTags Returns a valid ResourceTags type from AWS Instance Tags
func ParseResourceTags(tags []AWSTag) provider.ResourceTags {
	t := make(provider.ResourceTags)

	for _, tag := range tags {
		t[*tag.Key] = *tag.Value
	}
	return t
}
