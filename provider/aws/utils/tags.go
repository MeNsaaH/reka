package utils

import (
	"github.com/mensaah/reka/types"
)

// AWSTag type to convert all aws tags to before converting to the default `types.ResourceTags`
type AWSTag struct {
	Key   *string
	Value *string
}

// ParseResourceTags Returns a valid ResourceTags type from AWS Instance Tags
func ParseResourceTags(tags []AWSTag) types.ResourceTags {
	t := make(types.ResourceTags)

	for _, tag := range tags {
		t[*tag.Key] = *tag.Value
	}
	return t
}
