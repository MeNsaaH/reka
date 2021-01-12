package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/mensaah/reka/provider/aws/utils"
	"github.com/mensaah/reka/resource"
)

func getNodegroupDetails(svc *eks.Client, clusterName string) ([]*resource.Resource, error) {
	var nodegroups []*resource.Resource
	params := &eks.ListNodegroupsInput{ClusterName: &clusterName}
	ngList, err := svc.ListNodegroups(context.TODO(), params)
	if err != nil {
		return []*resource.Resource{}, err
	}
	for _, i := range ngList.Nodegroups {
		params := &eks.DescribeNodegroupInput{ClusterName: &clusterName, NodegroupName: &i}
		resp, err := svc.DescribeNodegroup(context.TODO(), params)
		if err != nil {
			return []*resource.Resource{}, err
		}
		// https://stackoverflow.com/a/48554123/7167357
		tags := resource.Tags(resp.Nodegroup.Tags)
		tags["creation-date"] = (*resp.Nodegroup.CreatedAt).String()

		ngResource := NewResource(i, nodegroupName)
		ngResource.CreationDate = *resp.Nodegroup.CreatedAt
		ngResource.Attributes["DesiredSize"] = resp.Nodegroup.ScalingConfig.DesiredSize
		ngResource.Status = utils.GetEksResourceStatus(string(resp.Nodegroup.Status))
		nodegroups = append(nodegroups, ngResource)
	}
	return nodegroups, nil
}

func resizeNodeGroup(svc *eks.Client, clusterName string, ngName string, size int32) error {
	params := &eks.UpdateNodegroupConfigInput{
		ClusterName:   &clusterName,
		NodegroupName: &ngName,
		ScalingConfig: &types.NodegroupScalingConfig{DesiredSize: &size},
	}
	_, err := svc.UpdateNodegroupConfig(context.TODO(), params)
	return err
}

func getClusterDetails(svc *eks.Client, output *eks.ListClustersOutput) []*resource.Resource {
	var eksClusters []*resource.Resource
	eksLogger.Debug("Fetching EKS Details")
	for _, c := range output.Clusters {
		clusterDetailsInput := &eks.DescribeClusterInput{Name: &c}
		// Build the request with its input parameters
		resp, err := svc.DescribeCluster(context.TODO(), clusterDetailsInput)
		if err != nil {
			eksLogger.Errorf("Failed to get details for cluster %s: %s", c, err)
			continue
		}
		cluster := resp.Cluster
		nodeGroups, err := getNodegroupDetails(svc, *cluster.Name)
		if err != nil {
			eksLogger.Errorf("Failed to get nodegroup for cluster %s: %s", c, err)
			continue
		}

		// https://stackoverflow.com/a/48554123/7167357
		tags := resource.Tags(cluster.Tags)
		tags["creation-date"] = (*cluster.CreatedAt).String()

		eksResource := NewResource(*cluster.Name, eksName)
		// Get CreationDate by getting LaunchTime of attached Volume
		eksResource.CreationDate = *cluster.CreatedAt
		eksResource.SubResources = make(map[string][]*resource.Resource)
		eksResource.SubResources[nodegroupName] = nodeGroups
		eksResource.Tags = tags
		eksResource.Status = utils.GetEksResourceStatus(string(cluster.Status))
		eksClusters = append(eksClusters, eksResource)
	}

	return eksClusters
}

// GetAllEKSClusters Get all clusters
func GetAllEKSClusters(cfg aws.Config) ([]*resource.Resource, error) {
	eksLogger.Debug("Fetching EKS Clusters")

	svc := eks.NewFromConfig(cfg)
	params := &eks.ListClustersInput{}

	// Build the request with its input parameters
	resp, err := svc.ListClusters(context.TODO(), params)

	if err != nil {
		return nil, err
	}
	clusters := getClusterDetails(svc, resp)
	eksLogger.Debugf("Found %d EKS clusters", len(clusters))
	return clusters, nil
}

// StopEKSClusters Stop Running Clusters
func StopEKSClusters(cfg aws.Config, clusters []*resource.Resource) error {
	svc := eks.NewFromConfig(cfg)
	var stoppableClusters []*resource.Resource

	for _, cluster := range clusters {
		if cluster.IsActive() {
			stoppableClusters = append(stoppableClusters, cluster)
		}
	}

	if len(clusters) <= 0 {
		return nil
	}

	for _, clsr := range stoppableClusters {
		eksLogger.Debugf("Stopping EKS Clusters %s ...", clsr)
		for _, ng := range clsr.SubResources[nodegroupName] {
			err := resizeNodeGroup(svc, clsr.UUID, ng.UUID, 0)
			if err != nil {
				eksLogger.Errorf("Failed Stopping Nodegroup %s in cluster %s: %s", ng, clsr, err)
			}
		}
	}
	return nil
}

// ResumeEKSClusters Resume Stopped clusters
func ResumeEKSClusters(cfg aws.Config, clusters []*resource.Resource) error {
	svc := eks.NewFromConfig(cfg)

	var resumableClusters []*resource.Resource
	for _, cluster := range clusters {
		if cluster.IsStopped() {
			resumableClusters = append(resumableClusters, cluster)
		}
	}

	if len(clusters) <= 0 {
		return nil
	}
	for _, clsr := range clusters {
		desired, err := utils.GetResourceFromDesiredState(providerName, eksName, clsr.UUID)
		if err != nil {
			eksLogger.Error(err.Error())
			continue
		}
		eksLogger.Debugf("Stopping EKS Clusters %s ...", clsr)
		for _, ng := range desired.SubResources[nodegroupName] {
			desiredSize, _ := ng.Attributes["DesiredSize"].(int32)
			err := resizeNodeGroup(svc, clsr.UUID, ng.UUID, desiredSize)
			if err != nil {
				eksLogger.Errorf("Failed Stopping Nodegroup %s in cluster %s: %s", ng, clsr, err)
			}
		}
	}
	return nil
}

// TerminateEKSClusters Shutdown clusters
func TerminateEKSClusters(cfg aws.Config, clusters []*resource.Resource) error {
	svc := eks.NewFromConfig(cfg)

	var targetClusters []*resource.Resource
	for _, cluster := range clusters {
		if cluster.IsStopped() || cluster.IsActive() {
			targetClusters = append(targetClusters, cluster)
		}
	}

	if len(clusters) <= 0 {
		return nil
	}
	for _, cluster := range targetClusters {
		params := &eks.DeleteClusterInput{Name: &cluster.UUID}
		_, err := svc.DeleteCluster(context.TODO(), params)
		if err != nil {
			eksLogger.Errorf("Error Deleting Cluster %s: %s", cluster, err)
		}
	}

	eksLogger.Debug("Terminating EKS Clusters ", clusters, " ...")

	return nil
}
