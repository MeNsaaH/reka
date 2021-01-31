package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/mensaah/reka/provider/aws/utils"

	"github.com/mensaah/reka/resource"
)

// returns only instance IDs of unprotected rds instances
func getrdsInstanceDetails(svc *rds.Client, output *rds.DescribeDBClustersOutput, region string) ([]*resource.Resource, error) {
	var rdsInstances []*resource.Resource
	rdsLogger.Debug("Fetching RDS Details")
	for _, instance := range output.DBClusters {
		tags := make(resource.Tags)
		for _, t := range instance.TagList {
			tags[*t.Key] = *t.Value
		}
		tags["creation-date"] = (*instance.ClusterCreateTime).String()
		rds := NewResource(*instance.DBClusterIdentifier, rdsName)
		rds.Region = region
		// Get CreationDate by getting LaunchTime of attached Volume
		rds.CreationDate = *instance.ClusterCreateTime
		rds.Tags = tags
		rds.Status = utils.GetRDSStatus(*instance.Status)
		rdsInstances = append(rdsInstances, rds)
	}

	return rdsInstances, nil
}

// GetAllRDSInstances Get all instances
func GetAllRDSInstances(cfg aws.Config) ([]*resource.Resource, error) {
	rdsLogger.Debug("Fetching RDS Clusters")

	svc := rds.NewFromConfig(cfg)
	params := &rds.DescribeDBClustersInput{}

	// Build the request with its input parameters
	resp, err := svc.DescribeDBClusters(context.TODO(), params)

	if err != nil {
		return nil, err
	}
	instances, err := getrdsInstanceDetails(svc, resp, cfg.Region)
	if err != nil {
		return nil, err
	}
	rdsLogger.Debugf("Found %d RDS clusters", len(instances))
	return instances, nil
}

// StopRDSInstances Stop Running Instances
func StopRDSInstances(cfg aws.Config, instances []*resource.Resource) error {
	svc := rds.NewFromConfig(cfg)
	var instanceIds []rds.StopDBClusterInput

	for _, instance := range instances {
		if instance.IsActive() {
			instanceIds = append(instanceIds, rds.StopDBClusterInput{
				DBClusterIdentifier: &instance.UUID,
			})
		}
	}

	if len(instanceIds) <= 0 {
		return nil
	}

	for _, instance := range instanceIds {
		rdsLogger.Debug("Stopping RDS Clusters ", instance.DBClusterIdentifier, " ...")
		resp, err := svc.StopDBCluster(context.TODO(), &instance)
		// TODO Attach error to specific instance where the error occurred if possible
		if err != nil {
			fmt.Println(resp, err)
		}
		return err
	}
	return nil
}

// ResumeRDSInstances Resume Stopped instances
func ResumeRDSInstances(cfg aws.Config, instances []*resource.Resource) error {
	svc := rds.NewFromConfig(cfg)
	var instanceIds []*rds.StartDBClusterInput

	for _, instance := range instances {
		if instance.IsStopped() {
			instanceIds = append(instanceIds, &rds.StartDBClusterInput{
				DBClusterIdentifier: &instance.UUID,
			})
		}
	}

	if len(instanceIds) <= 0 {
		return nil
	}

	for _, instance := range instanceIds {
		rdsLogger.Debug("Starting RDS Cluster ", instance.DBClusterIdentifier, " ...")
		resp, err := svc.StartDBCluster(context.TODO(), instance)
		// TODO Attach error to specific instance where the error occurred if possible
		if err != nil {
			fmt.Println(resp, err)
		}
		return err
	}
	return nil
}

// TerminateRDSInstances Shutdown instances
func TerminateRDSInstances(cfg aws.Config, instances []*resource.Resource) error {
	svc := rds.NewFromConfig(cfg)
	var instanceIds []*rds.DeleteDBClusterInput

	for _, instance := range instances {
		if instance.IsStopped() || instance.IsActive() {
			instanceIds = append(instanceIds, &rds.DeleteDBClusterInput{
				DBClusterIdentifier: &instance.UUID,
			})
		}
	}

	if len(instanceIds) <= 0 {
		return nil
	}

	for _, instance := range instanceIds {
		rdsLogger.Debug("Terminating RDS Cluster ", instance.DBClusterIdentifier, " ...")
		resp, err := svc.DeleteDBCluster(context.TODO(), instance)
		// TODO Attach error to specific instance where the error occurred if possible
		if err != nil {
			fmt.Println(resp, err)
		}
		return err
	}
	return nil
}
