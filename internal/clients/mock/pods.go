package mock

import (
	"cmyk/internal/models"
	"fmt"

	"encoding/json"
	"os"
	"strings"
)

// ListPods reads and parses the mock pods data from JSON file
func (c Client) ListPods() ([]models.Pod, error) {
	data, err := os.ReadFile("./internal/clients/mock/pods.json")
	if err != nil {
		return nil, err
	}

	var pods models.PodList
	if err := json.Unmarshal(data, &pods); err != nil {
		return nil, err
	}

	var result []models.Pod
	for _, p := range pods.Items {
		var restarts int
		for _, cs := range p.Status.ContainerStatuses {
			restarts += cs.RestartCount
		}

		statusClass := "pending"
		switch strings.ToLower(p.Status.Phase) {
		case "running":
			statusClass = "running"
		case "succeeded":
			statusClass = "ready"
		case "failed":
			statusClass = "notready"
		}

		result = append(result, models.Pod{
			Name:        p.Metadata.Name,
			Namespace:   p.Metadata.Namespace,
			Status:      p.Status.Phase,
			StatusClass: statusClass,
			Node:        p.Spec.NodeName,
			PodIP:       p.Status.PodIP,
			Restarts:    restarts,
		})
	}

	return result, nil
}

// GetPod reads and parses a single pod from mock data
//
//revive:disable:cyclomatic
func (c Client) GetPod(namespace, name string) (*models.PodDetail, error) {
	data, err := os.ReadFile("./internal/clients/mock/pods.json")
	if err != nil {
		return nil, err
	}

	var pods models.PodList
	if err := json.Unmarshal(data, &pods); err != nil {
		return nil, err
	}

	for _, p := range pods.Items {
		if p.Metadata.Name != name || p.Metadata.Namespace != namespace {
			continue
		}

		statusClass := "pending"
		switch strings.ToLower(p.Status.Phase) {
		case "running":
			statusClass = "running"
		case "succeeded":
			statusClass = "ready"
		case "failed":
			statusClass = "notready"
		}

		// Build containers list
		var containers []models.PodContainer
		for _, c := range p.Spec.Containers {
			container := models.PodContainer{
				Name:  c.Name,
				Image: c.Image,
			}

			// Ports
			for _, port := range c.Ports {
				container.Ports = append(container.Ports, models.ContainerPort{
					Name:          port.Name,
					ContainerPort: port.ContainerPort,
					Protocol:      port.Protocol,
				})
			}

			// Environment variables
			for _, e := range c.Env {
				container.Env = append(container.Env, models.EnvVar{
					Name:  e.Name,
					Value: e.Value,
				})
			}

			// Volume mounts
			for _, vm := range c.VolumeMounts {
				container.VolumeMounts = append(container.VolumeMounts, models.VolumeMount{
					Name:      vm.Name,
					MountPath: vm.MountPath,
					ReadOnly:  vm.ReadOnly,
				})
			}

			// Resources
			container.Resources = models.ResourceRequirements{
				Limits:   c.Resources.Limits,
				Requests: c.Resources.Requests,
			}

			// Get status from container statuses
			for _, cs := range p.Status.ContainerStatuses {
				if cs.Name == c.Name {
					container.Ready = cs.Ready
					container.RestartCount = cs.RestartCount
					if cs.State.Running != nil {
						container.State = "Running"
					} else if cs.State.Waiting != nil {
						container.State = "Waiting"
						container.StateReason = cs.State.Waiting.Reason
					} else if cs.State.Terminated != nil {
						container.State = "Terminated"
						container.StateReason = cs.State.Terminated.Reason
					}
					break
				}
			}

			containers = append(containers, container)
		}

		// Build init containers list
		var initContainers []models.PodContainer
		for _, c := range p.Spec.InitContainers {
			container := models.PodContainer{
				Name:  c.Name,
				Image: c.Image,
			}

			for _, cs := range p.Status.InitContainerStatuses {
				if cs.Name == c.Name {
					container.Ready = cs.Ready
					container.RestartCount = cs.RestartCount
					if cs.State.Running != nil {
						container.State = "Running"
					} else if cs.State.Waiting != nil {
						container.State = "Waiting"
						container.StateReason = cs.State.Waiting.Reason
					} else if cs.State.Terminated != nil {
						container.State = "Terminated"
						container.StateReason = cs.State.Terminated.Reason
					}
					break
				}
			}

			initContainers = append(initContainers, container)
		}

		// Build conditions list
		var conditions []models.PodCondition
		for _, cond := range p.Status.Conditions {
			conditions = append(conditions, models.PodCondition{
				Type:               cond.Type,
				Status:             cond.Status,
				Reason:             cond.Reason,
				Message:            cond.Message,
				LastTransitionTime: cond.LastTransitionTime,
			})
		}

		// Build volumes list
		var volumes []models.PodVolume
		for _, v := range p.Spec.Volumes {
			vol := models.PodVolume{Name: v.Name}
			switch {
			case v.ConfigMap != nil:
				vol.Type = "ConfigMap"
				vol.Source = v.ConfigMap.Name
			case v.Secret != nil:
				vol.Type = "Secret"
				vol.Source = v.Secret.SecretName
			case v.PersistentVolumeClaim != nil:
				vol.Type = "PVC"
				vol.Source = v.PersistentVolumeClaim.ClaimName
			case v.HostPath != nil:
				vol.Type = "HostPath"
				vol.Source = v.HostPath.Path
			case v.EmptyDir != nil:
				vol.Type = "EmptyDir"
				vol.Source = "-"
			case v.Projected != nil:
				vol.Type = "Projected"
				vol.Source = "-"
			default:
				vol.Type = "Other"
				vol.Source = "-"
			}
			volumes = append(volumes, vol)
		}

		// Build owner references
		var ownerRefs []models.OwnerReference
		for _, or := range p.Metadata.OwnerReferences {
			ownerRefs = append(ownerRefs, models.OwnerReference{
				APIVersion: or.APIVersion,
				Kind:       or.Kind,
				Name:       or.Name,
				UID:        or.UID,
			})
		}

		return &models.PodDetail{
			Name:              p.Metadata.Name,
			Namespace:         p.Metadata.Namespace,
			UID:               p.Metadata.UID,
			CreationTimestamp: p.Metadata.CreationTimestamp,
			Labels:            p.Metadata.Labels,
			Annotations:       p.Metadata.Annotations,
			Status:            p.Status.Phase,
			StatusClass:       statusClass,
			Node:              p.Spec.NodeName,
			PodIP:             p.Status.PodIP,
			HostIP:            p.Status.HostIP,
			QOSClass:          p.Status.QOSClass,
			ServiceAccount:    p.Spec.ServiceAccountName,
			Containers:        containers,
			InitContainers:    initContainers,
			Conditions:        conditions,
			Volumes:           volumes,
			OwnerReferences:   ownerRefs,
		}, nil
	}

	return nil, fmt.Errorf("pod %s/%s not found", namespace, name)
	//revive:enable:cyclomatic
}
