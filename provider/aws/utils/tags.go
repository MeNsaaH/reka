package utils

import (
	"github.com/mensaah/reka/resource"
)

// AWSTag type to convert all aws tags to before converting to the default `resource.Tags`
type AWSTag struct {
	Key   *string
	Value *string
}

// ParseTags Returns a valid Tags type from AWS Instance Tags
func ParseTags(tags []AWSTag) resource.Tags {
	t := make(resource.Tags)

	for _, tag := range tags {
		t[*tag.Key] = *tag.Value
	}
	return t
}
