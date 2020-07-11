package aws

import ()

// EC2 The EC2 Resource
type EC2 struct {
	AWSResource
}

func (e *EC2) destroy() (string, error) {
	return "", nil
}

// NewEC2 Returns a new EC2 Resource object
func NewEC2(id string) EC2 {
	resource := EC2{}
	resource.ID = id
	resource.Name = "EC2"
	resource.LongName = "Elastic Compute Cloud"

	return resource
}
