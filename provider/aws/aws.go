package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"

	"github.com/mensaah/reka/provider"
	"github.com/mensaah/reka/provider/aws/ec2"
)

func getState(s int64) provider.State {
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

func NewProvider() provider.Provider {
	aws := provider.Provider{}
	aws.Name = "AWS"
	aws.Config = &provider.Config{}
	config := GetConfig()
	ec2Manager := ec2.InitManager(config)

	resources := map[string]provider.ResourceManager{
		ec2.Name: &ec2Manager,
	}
	aws.ResourceManagers = resources
	return aws
}
