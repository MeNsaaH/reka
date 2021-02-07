package gcp

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	gke "google.golang.org/api/container/v1"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider/gcp/utils"
	"github.com/mensaah/reka/resource"
)

func resizeNodePool(svc *gke.ProjectsLocationsClustersNodePoolsService, project string, cluster string, np string, location string, size int64) error {
	name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s/nodePools/%s", project, location, cluster, np)
	sizeReq := gke.SetNodePoolSizeRequest{
		ClusterId: cluster,
		Name:      name,
		NodeCount: size,
	}
	_, err := svc.SetSize(np, &sizeReq).Do()
	if err != nil {
		return err
	}
	return nil
}

func getNodePoolsDetails(svc *gke.ProjectsLocationsClustersService, cluster *gke.Cluster) ([]*resource.Resource, error) {
	var nodePools []*resource.Resource

	for _, i := range cluster.NodePools {
		np := NewResource(fmt.Sprint(i.Name), nodePoolName)
		np.Tags = i.Config.Labels
		np.Status = resource.Running
		np.Attributes["ActualNodeCount"] = i.InitialNodeCount
		nodePools = append(nodePools, np)
	}

	return nodePools, nil
}

func getGkeClusters(svc *gke.ProjectsLocationsClustersService, projectId string) ([]*resource.Resource, error) {
	var gkeClusters []*resource.Resource
	parent := fmt.Sprintf("projects/%s/locations/-", projectId)
	clusters, err := svc.List(parent).Do()
	if err != nil {
		return []*resource.Resource{}, err
	}

	for _, i := range clusters.Clusters {
		cluster := NewResource(fmt.Sprint(i.Name), gkeName)

		cluster.Location = i.Location
		cluster.Status = utils.GetComputeInstanceStatus(i.Status)
		creationDate, err := time.Parse(time.RFC3339, i.CreateTime)
		if err != nil {
			gkeLogger.Error("Could not parse creation time for cluster %s, value %s", i.Name, i.CreateTime)
		}
		cluster.CreationDate = creationDate
		cluster.Tags = i.ResourceLabels
		if cluster.Status != resource.Error {
			// Add Node Pool Data to Cluster
			cluster.SubResources = make(map[string][]*resource.Resource)

			nodePools, err := getNodePoolsDetails(svc, i)
			if err != nil {
				gkeLogger.Error(err)
				continue
			}
			cluster.SubResources[nodePoolName] = nodePools
		}

		gkeClusters = append(gkeClusters, cluster)
	}
	log.Debugf("Found %d gke Instances", len(gkeClusters))
	return gkeClusters, nil
}

func getAllGkeClusters(cfg *config.Gcp) ([]*resource.Resource, error) {
	var gkeClusters []*resource.Resource
	gkeLogger.Debug("Fetching GKE Clusters")
	ctx := context.Background()
	svc, err := gke.NewService(ctx)
	if err != nil {
		return []*resource.Resource{}, err
	}
	client := gke.NewProjectsLocationsClustersService(svc)
	clusters, err := getGkeClusters(client, cfg.ProjectId)
	if err != nil {
		return []*resource.Resource{}, err
	}
	gkeClusters = append(gkeClusters, clusters...)
	return gkeClusters, nil
}

// https://cloud.google.com/gke/docs/reference/latest/clusters/stop
func stopGkeClusters(cfg *config.Gcp, clusters []*resource.Resource) error {
	log.Debug("Fetching  cluster")
	// TODO Allow configurable options for setting stop cluster size
	stopSize := int64(0)

	var selectedClusters []*resource.Resource
	for _, cluster := range clusters {
		if cluster.IsActive() {
			selectedClusters = append(selectedClusters, cluster)
		}
	}

	if len(clusters) <= 0 {
		return nil
	}

	ctx := context.Background()
	svc, err := gke.NewService(ctx)
	if err != nil {
		gkeLogger.Error(err)
	}
	client := gke.NewProjectsLocationsClustersNodePoolsService(svc)

	for _, cluster := range clusters {
		// TODO Add operation Waiter to check status of stop operation
		// TODO [GKE] Also allow users to be able to specify the type of stop operation to perform (Suspend/Stop)
		// https://cloud.google.com/gke/docs/clusters/cluster-life-cycle#comparison_table
		for _, np := range cluster.SubResources[nodePoolName] {
			err := resizeNodePool(client, cfg.ProjectId, cluster.UUID, np.UUID, cluster.Location, stopSize)
			if err != nil {
				gkeLogger.Error(err)
			}
		}
	}
	return nil
}

func startGkeClusters(cfg *config.Gcp, clusters []*resource.Resource) error {
	log.Debug("Fetching cluster")
	// TODO Allow configurable options for setting stop cluster size

	var selectedClusters []*resource.Resource
	for _, cluster := range clusters {
		if cluster.IsActive() {
			selectedClusters = append(selectedClusters, cluster)
		}
	}

	if len(selectedClusters) <= 0 {
		return nil
	}

	ctx := context.Background()
	svc, err := gke.NewService(ctx)
	if err != nil {
		gkeLogger.Error(err)
	}
	client := gke.NewProjectsLocationsClustersNodePoolsService(svc)

	for _, cluster := range selectedClusters {
		desired, err := utils.GetResourceFromDesiredState(providerName, gkeName, cluster.UUID)
		if err != nil {
			gkeLogger.Error(err)
			continue
		}
		// TODO Add operation Waiter to check status of stop operation
		// TODO [GKE] Also allow users to be able to specify the type of stop operation to perform (Suspend/Stop)
		// https://cloud.google.com/gke/docs/clusters/cluster-life-cycle#comparison_table
		for _, np := range desired.SubResources[nodePoolName] {
			desiredSize, _ := np.Attributes["ActualNodeCount"].(int64)
			err := resizeNodePool(client, cfg.ProjectId, cluster.UUID, np.UUID, cluster.Region, desiredSize)
			if err != nil {
				gkeLogger.Error(err)
			}
		}
	}
	return nil
}

func destroyGkeClusters(cfg *config.Gcp, clusters []*resource.Resource) error {
	var selectedClusters []*resource.Resource
	for _, cluster := range clusters {
		if cluster.IsActive() {
			selectedClusters = append(selectedClusters, cluster)
		}
	}

	if len(selectedClusters) <= 0 {
		return nil
	}

	log.Debugf("Destroying %d clusters", len(selectedClusters))

	ctx := context.Background()
	svc, err := gke.NewService(ctx)
	if err != nil {
		gkeLogger.Error(err)
	}
	client := gke.NewProjectsLocationsClustersService(svc)

	for _, cluster := range selectedClusters {
		// TODO Add operation Waiter to check status of delete operation
		name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", cfg.ProjectId, cluster.Location, cluster.UUID)
		_, err := client.Delete(name).Do()
		if err != nil {
			gkeLogger.Error(err)
		}
	}
	return nil
}
