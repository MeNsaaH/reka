package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider/aws/ec2"
	"github.com/mensaah/reka/provider/aws/s3"
	"github.com/mensaah/reka/types"
)

const (
	ProviderName = "aws"
)

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

// NewProvider : Creates a New AWS Provider
func NewProvider() (*types.Provider, error) {

	aws := types.Provider{}
	aws.Name = ProviderName

	// Get and Load AWS Config
	awsConfig := config.GetAWS()
	// Set AWS Config
	awsConfig.Config = GetConfig()

	ec2Manager := ec2.InitManager()
	s3Manager := s3.InitManager()

	resMgrs := map[string]types.ResourceManager{
		ec2.Name: &ec2Manager,
		s3.Name:  &s3Manager,
	}
	aws.ResourceManagers = resMgrs
	return &aws, nil
}
