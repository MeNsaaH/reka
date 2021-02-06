package gcp

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	compute "google.golang.org/api/compute/v1"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider/gcp/utils"
	"github.com/mensaah/reka/resource"
)

func getComputeInstancesInZone(svc *compute.InstancesService, projectId string, zone string) ([]*resource.Resource, error) {
	var computeInstances []*resource.Resource
	instances, err := svc.List(projectId, zone).Do()
	if err != nil {
		return []*resource.Resource{}, err
	}

	for _, i := range instances.Items {
		computeInstance := NewResource(fmt.Sprint(i.Id), computeInstanceName)
		computeInstance.Zone = zone
		computeInstance.Status = utils.GetComputeInstanceStatus(i.Status)
		creationDate, err := time.Parse(time.RFC3339, i.CreationTimestamp)
		if err != nil {
			computeLogger.Error("Could not parse creation time for instance %s, value %s", i.Id, i.CreationTimestamp)
		}
		computeInstance.CreationDate = creationDate
		computeInstance.Tags = i.Labels
		computeInstances = append(computeInstances, computeInstance)
	}
	log.Debugf("Found %d compute Instances in %s zone", len(computeInstances), zone)
	return computeInstances, nil
}

func getAllComputeInstances(cfg *config.Gcp) ([]*resource.Resource, error) {
	var computeInstances []*resource.Resource
	computeLogger.Debug("Fetching compute Instances")
	ctx := context.Background()
	svc, err := compute.NewService(ctx)
	if err != nil {
		return []*resource.Resource{}, err
	}
	client := compute.NewInstancesService(svc)
	zone := "us-east1-b"
	zoneInstances, err := getComputeInstancesInZone(client, cfg.ProjectId, zone)
	if err != nil {
		return []*resource.Resource{}, err
	}
	computeInstances = append(computeInstances, zoneInstances...)
	return computeInstances, nil
}

func stopComputeInstances(cfg *config.Gcp, instances []*resource.Resource) error {
	log.Debug("Fetching Cloud storage computeInstance")
	ctx := context.Background()
	svc, err := compute.NewService(ctx)
	if err != nil {
		computeLogger.Error(err)
	}
	client := compute.NewInstancesService(svc)

	for _, instance := range instances {
		// TODO Add operation Waiter to check status of stop operation
		// TODO Also allow users to be able to specify the type of stop operation to perform (Suspend/Stop)
		// https://cloud.google.com/compute/docs/instances/instance-life-cycle#comparison_table
		_, err := client.Stop(cfg.ProjectId, instance.Zone, instance.UUID).Do()
		if err != nil {
			computeLogger.Error(err)
		}
	}
	return nil
}

func startComputeInstances(cfg *config.Gcp, instances []*resource.Resource) error {
	log.Debug("Fetching Cloud storage computeInstance")
	ctx := context.Background()
	svc, err := compute.NewService(ctx)
	if err != nil {
		computeLogger.Error(err)
	}
	client := compute.NewInstancesService(svc)

	for _, instance := range instances {
		// TODO Add operation Waiter to check status of start operation
		_, err := client.Start(cfg.ProjectId, instance.Zone, instance.UUID).Do()
		if err != nil {
			computeLogger.Error(err)
		}
	}
	return nil
}

func destroyComputeInstances(cfg *config.Gcp, instances []*resource.Resource) error {
	log.Debug("Fetching Cloud storage computeInstance")
	ctx := context.Background()
	svc, err := compute.NewService(ctx)
	if err != nil {
		computeLogger.Error(err)
	}
	client := compute.NewInstancesService(svc)

	for _, instance := range instances {
		// TODO Add operation Waiter to check status of delete operation
		_, err := client.Delete(cfg.ProjectId, instance.Zone, instance.UUID).Do()
		if err != nil {
			computeLogger.Error(err)
		}
	}
	return nil
}
