package provider

import ()

// All Type and Interface Definitions to be used by all Providers (GCP, AWS, Azure) etc

// Provider : The Provider Interface
type Provider interface {
	getAllResources() []interface{}
	getReapableResources() []interface{}
	destroyResources([]interface{}) (string, error)
	Nuke([]interface{}) (string, error)
}
