package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/mensaah/reka/provider"
	//  "github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

// AWSResource : An AWS Specific Resource
type AWSResource struct {
	provider.Resource
	// Specific Details to AWS Instances
	Region string
	Tags   provider.ResourceTags
}

func getState(s int64) provider.State {
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

func GetConfig() aws.Config {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region that the service clients should use
	//  cfg.Region = endpoints.UsEast2RegionID
	cfg.Region = "us-east-2"
	return cfg
}
