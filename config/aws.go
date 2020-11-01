package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/credentials"
	log "github.com/sirupsen/logrus"
)

func loadAwsConfig(accessKeyID, secretAccessKey, defaultRegion string) aws.Config {
	var (
		err error
		cfg aws.Config
	)

	if accessKeyID != "" && secretAccessKey != "" {
		cfg, err = awsCfg.LoadDefaultConfig(
			awsCfg.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID: accessKeyID, SecretAccessKey: secretAccessKey,
					Source: "Reka Variables",
				},
			}), awsCfg.WithRegion(defaultRegion))
	} else {
		cfg, err = awsCfg.LoadDefaultConfig(awsCfg.WithRegion(defaultRegion))
	}

	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
