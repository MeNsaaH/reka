package utils

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
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
		return provider.Terminated
	case 64:
		return provider.Stopping
	case 80:
		return provider.Stopped
	default:
		return provider.Stopped
	}
}

// ParseResourceTags Returns a valid ResourceTags type from AWS Instance Tags
func ParseResourceTags(tags interface{}) provider.ResourceTags {
	t := make(provider.ResourceTags)

	switch v := tags.(type) {
	case []ec2.Tag:
		for _, tag := range v {
			t[*tag.Key] = *tag.Value
		}
	}
	return t
}
