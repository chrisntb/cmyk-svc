package mock

import (
	"cmyk/internal/models"

	"encoding/json"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ListNodes reads and parses the mock nodes data from JSON file
func (c Client) ListNodes() ([]models.Node, error) {
	data, err := os.ReadFile("./internal/clients/mock/nodes.json")
	if err != nil {
		return nil, err
	}

	var nodes models.NodeList
	if err := json.Unmarshal(data, &nodes); err != nil {
		return nil, err
	}

	var result []models.Node
	for _, n := range nodes.Items {
		var ip string
		for _, addr := range n.Status.Addresses {
			if addr.Type == "InternalIP" {
				ip = addr.Address
				break
			}
		}

		var ready bool
		for _, cond := range n.Status.Conditions {
			if cond.Type == "Ready" && cond.Status == "True" {
				ready = true
				break
			}
		}

		var roles []string
		for label := range n.Metadata.Labels {
			if strings.HasPrefix(label, "node-role.kubernetes.io/") {
				role := strings.TrimPrefix(label, "node-role.kubernetes.io/")
				if role != "" {
					roles = append(roles, role)
				}
			}
		}
		roleStr := strings.Join(roles, ", ")
		if roleStr == "" {
			roleStr = "<none>"
		}

		result = append(result, models.Node{
			CPU:            n.Status.Capacity["cpu"],
			IP:             ip,
			KubeletVersion: n.Status.NodeInfo.KubeletVersion,
			Memory:         n.Status.Capacity["memory"],
			Name:           n.Metadata.Name,
			Ready:          ready,
			Roles:          roleStr,
		})
	}

	return result, nil
}

// GetNode reads and parses a mock node data from JSON file
func (c Client) GetNode(name string) (*models.NodeDetail, error) {
	data, err := os.ReadFile("./internal/clients/mock/nodes.json")
	if err != nil {
		return nil, err
	}

	var nodes models.NodeList
	if err := json.Unmarshal(data, &nodes); err != nil {
		return nil, err
	}

	for _, n := range nodes.Items {
		if n.Metadata.Name != name {
			continue
		}

		var roles []string
		for label := range n.Metadata.Labels {
			if strings.HasPrefix(label, "node-role.kubernetes.io/") {
				role := strings.TrimPrefix(label, "node-role.kubernetes.io/")
				if role != "" {
					roles = append(roles, role)
				}
			}
		}
		roleStr := strings.Join(roles, ", ")
		if roleStr == "" {
			roleStr = "<none>"
		}

		var ready bool
		for _, cond := range n.Status.Conditions {
			if cond.Type == "Ready" && cond.Status == "True" {
				ready = true
				break
			}
		}

		var taints []models.NodeTaint
		for _, t := range n.Spec.Taints {
			taints = append(taints, models.NodeTaint{
				Key:    t.Key,
				Value:  t.Value,
				Effect: t.Effect,
			})
		}

		var addresses []models.NodeAddress
		for _, a := range n.Status.Addresses {
			addresses = append(addresses, models.NodeAddress{
				Type:    a.Type,
				Address: a.Address,
			})
		}

		var conditions []models.NodeCondition
		for _, c := range n.Status.Conditions {
			conditions = append(conditions, models.NodeCondition{
				Type:               c.Type,
				Status:             c.Status,
				Reason:             c.Reason,
				Message:            c.Message,
				LastHeartbeatTime:  c.LastHeartbeatTime,
				LastTransitionTime: c.LastTransitionTime,
			})
		}

		var images []models.NodeImage
		for _, img := range n.Status.Images {
			images = append(images, models.NodeImage{
				Names:     img.Names,
				SizeBytes: img.SizeBytes,
			})
		}

		return &models.NodeDetail{
			Addresses:         addresses,
			Allocatable:       models.NodeResources(n.Status.Allocatable),
			Annotations:       n.Metadata.Annotations,
			Capacity:          models.NodeResources(n.Status.Capacity),
			Conditions:        conditions,
			CreationTimestamp: n.Metadata.CreationTimestamp,
			Images:            images,
			Labels:            n.Metadata.Labels,
			Name:              n.Metadata.Name,
			NodeInfo: models.NodeSystemInfo{
				Architecture:            n.Status.NodeInfo.Architecture,
				BootID:                  n.Status.NodeInfo.BootID,
				ContainerRuntimeVersion: n.Status.NodeInfo.ContainerRuntimeVersion,
				KernelVersion:           n.Status.NodeInfo.KernelVersion,
				KubeProxyVersion:        n.Status.NodeInfo.KubeProxyVersion,
				KubeletVersion:          n.Status.NodeInfo.KubeletVersion,
				MachineID:               n.Status.NodeInfo.MachineID,
				OperatingSystem:         n.Status.NodeInfo.OperatingSystem,
				OSImage:                 n.Status.NodeInfo.OSImage,
				SystemUUID:              n.Status.NodeInfo.SystemUUID,
			},
			PodCIDR: n.Spec.PodCIDR,
			Ready:   ready,
			Roles:   roleStr,
			Taints:  taints,
			UID:     n.Metadata.UID,
		}, nil
	}

	return nil, fiber.ErrNotFound
}
