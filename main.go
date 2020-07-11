package main

import (
	"fmt"

	//  "github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/mensaah/reka/provider/aws"
)

func main() {

	cfg := aws.GetConfig()

	instances, err := aws.GetAllEC2Instances(cfg, "us-east-2")
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}

	for _, i := range instances {
		fmt.Println(i.CreationDate)
	}

}

// TODO
// - Authenticate User
// - List all resources with specified tags and dependencies
// - Destroy Resources with specified Tags
