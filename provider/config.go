package provider

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

// Config : The Config values passed to application
type Config struct {

	// AWS Configs
	AwsConfig aws.Config
	AwsRegion string
}
