package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/mensaah/reka/provider"
)

func parseTags(tags interface{}) provider.ResourceTags {
	t := make(provider.ResourceTags)

	switch v := tags.(type) {
	case []ec2.Tag:
		for _, tag := range v {
			t[*tag.Key] = *tag.Value
		}
	}
	return t
}
