package utils

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

func GetServiceListWithStatus(ctx context.Context, c client.APIClient) ([]swarm.Service, error) {
	servicesList, err := c.ServiceList(ctx, types.ServiceListOptions{Status: true})
	if err != nil {
		return nil, err
	}

	return AppendServiceStatus(ctx, c, servicesList)
}

// older apis...
// https://github.com/docker/cli/blob/4f058143c705495b0302215a503fc8df82c40565/cli/command/service/list.go#L110
func AppendServiceStatus(ctx context.Context, c client.APIClient, services []swarm.Service) ([]swarm.Service, error) {
	status := map[string]*swarm.ServiceStatus{}
	taskFilter := filters.NewArgs()
	for i, s := range services {
		// there is no need in this switch to check for job modes. jobs are not
		// supported until after ServiceStatus was introduced.
		switch {
		case s.ServiceStatus != nil:
			// Server already returned service-status, so we don't
			// have to look-up tasks for this service.
			continue
		case s.Spec.Mode.Replicated != nil:
			// For replicated services, set the desired number of tasks;
			// that way we can present this information in case we're unable
			// to get a list of tasks from the server.
			services[i].ServiceStatus = &swarm.ServiceStatus{DesiredTasks: *s.Spec.Mode.Replicated.Replicas}
			status[s.ID] = &swarm.ServiceStatus{}
			taskFilter.Add("service", s.ID)
		case s.Spec.Mode.Global != nil:
			// No such thing as number of desired tasks for global services
			services[i].ServiceStatus = &swarm.ServiceStatus{}
			status[s.ID] = &swarm.ServiceStatus{}
			taskFilter.Add("service", s.ID)
		default:
			// Unknown task type
		}
	}
	if len(status) == 0 {
		// All services have their ServiceStatus set, so we're done
		return services, nil
	}

	tasks, err := c.TaskList(ctx, types.TaskListOptions{Filters: taskFilter})
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return services, nil
	}
	activeNodes, err := getActiveNodes(ctx, c)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		if status[task.ServiceID] == nil {
			// This should not happen in practice; either all services have
			// a ServiceStatus set, or none of them.
			continue
		}
		// TODO: this should only be needed for "global" services. Replicated
		// services have `Spec.Mode.Replicated.Replicas`, which should give this value.
		if task.DesiredState != swarm.TaskStateShutdown {
			status[task.ServiceID].DesiredTasks++
		}
		if _, nodeActive := activeNodes[task.NodeID]; nodeActive && task.Status.State == swarm.TaskStateRunning {
			status[task.ServiceID].RunningTasks++
		}
	}

	for i, service := range services {
		if s := status[service.ID]; s != nil {
			services[i].ServiceStatus = s
		}
	}
	return services, nil
}

func getActiveNodes(ctx context.Context, c client.NodeAPIClient) (map[string]struct{}, error) {
	nodes, err := c.NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		return nil, err
	}
	activeNodes := make(map[string]struct{})
	for _, n := range nodes {
		if n.Status.State != swarm.NodeStateDown {
			activeNodes[n.ID] = struct{}{}
		}
	}
	return activeNodes, nil
}
