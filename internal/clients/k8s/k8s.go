package k8s

import (
	"cmyk/internal/models"
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	kueueversioned "sigs.k8s.io/kueue/client-go/clientset/versioned"
)

type Client struct {
	Clientset      *kubernetes.Clientset
	KueueClientset kueueversioned.Interface
}

func New(kubeconfig string) (*Client, error) {
	// Create the client configuration from the kubeconfig file.
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed building kube config: %w", err)
	}

	// Configure client-side rate limiting.
	config.QPS = 50
	config.Burst = 100

	// A clientset contains clients for all the API groups and versions supported by the cluster.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed creating k8s clientset: %w", err)
	}

	kueueClientset, err := kueueversioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed creating kueue clientset: %w", err)
	}

	return &Client{Clientset: clientset, KueueClientset: kueueClientset}, nil
}

func (c Client) PodCountInDefaultNamespace() (int, error) {
	count := 0
	pods, err := c.Clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return count, fmt.Errorf("failed listing pods: %w", err)
	}
	count = len(pods.Items)
	return count, nil
}

func (c Client) ListNodes() ([]models.Node, error) {
	nodeList, err := c.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed listing nodes: %w", err)
	}

	var result []models.Node
	for _, n := range nodeList.Items {
		var ip string
		for _, addr := range n.Status.Addresses {
			if addr.Type == "InternalIP" {
				ip = string(addr.Address)
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
		for label := range n.Labels {
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
			Name:           n.Name,
			Ready:          ready,
			Roles:          roleStr,
			IP:             ip,
			CPU:            n.Status.Capacity.Cpu().String(),
			Memory:         n.Status.Capacity.Memory().String(),
			KubeletVersion: n.Status.NodeInfo.KubeletVersion,
		})
	}

	return result, nil
}

func (c Client) GetNode(name string) (*models.NodeDetail, error) {
	node, err := c.Clientset.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed getting node: %w", err)
	}

	var roles []string
	for label := range node.Labels {
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
	for _, cond := range node.Status.Conditions {
		if cond.Type == "Ready" && cond.Status == "True" {
			ready = true
			break
		}
	}

	var taints []models.NodeTaint
	for _, t := range node.Spec.Taints {
		taints = append(taints, models.NodeTaint{
			Key:    t.Key,
			Value:  t.Value,
			Effect: string(t.Effect),
		})
	}

	var addresses []models.NodeAddress
	for _, a := range node.Status.Addresses {
		addresses = append(addresses, models.NodeAddress{
			Type:    string(a.Type),
			Address: a.Address,
		})
	}

	var conditions []models.NodeCondition
	for _, cond := range node.Status.Conditions {
		conditions = append(conditions, models.NodeCondition{
			Type:               string(cond.Type),
			Status:             string(cond.Status),
			Reason:             cond.Reason,
			Message:            cond.Message,
			LastHeartbeatTime:  cond.LastHeartbeatTime.Format("2006-01-02T15:04:05Z"),
			LastTransitionTime: cond.LastTransitionTime.Format("2006-01-02T15:04:05Z"),
		})
	}

	var images []models.NodeImage
	for _, img := range node.Status.Images {
		images = append(images, models.NodeImage{
			Names:     img.Names,
			SizeBytes: img.SizeBytes,
		})
	}

	capacity := make(models.NodeResources)
	for name, quantity := range node.Status.Capacity {
		capacity[string(name)] = quantity.String()
	}

	allocatable := make(models.NodeResources)
	for name, quantity := range node.Status.Allocatable {
		allocatable[string(name)] = quantity.String()
	}

	return &models.NodeDetail{
		Addresses:         addresses,
		Allocatable:       allocatable,
		Annotations:       node.Annotations,
		Capacity:          capacity,
		Conditions:        conditions,
		CreationTimestamp: node.CreationTimestamp.Format("2006-01-02T15:04:05Z"),
		Images:            images,
		Labels:            node.Labels,
		Name:              node.Name,
		NodeInfo: models.NodeSystemInfo{
			Architecture:            node.Status.NodeInfo.Architecture,
			BootID:                  node.Status.NodeInfo.BootID,
			ContainerRuntimeVersion: node.Status.NodeInfo.ContainerRuntimeVersion,
			KernelVersion:           node.Status.NodeInfo.KernelVersion,
			KubeProxyVersion:        node.Status.NodeInfo.KubeProxyVersion,
			KubeletVersion:          node.Status.NodeInfo.KubeletVersion,
			MachineID:               node.Status.NodeInfo.MachineID,
			OperatingSystem:         node.Status.NodeInfo.OperatingSystem,
			OSImage:                 node.Status.NodeInfo.OSImage,
			SystemUUID:              node.Status.NodeInfo.SystemUUID,
		},
		PodCIDR: node.Spec.PodCIDR,
		Ready:   ready,
		Roles:   roleStr,
		Taints:  taints,
		UID:     string(node.UID),
	}, nil
}

func (c Client) ListPods() ([]models.Pod, error) {
	podList, err := c.Clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed listing pods: %w", err)
	}

	var result []models.Pod
	for _, p := range podList.Items {
		var restarts int
		for _, cs := range p.Status.ContainerStatuses {
			restarts += int(cs.RestartCount)
		}

		statusClass := "pending"
		switch strings.ToLower(string(p.Status.Phase)) {
		case "running":
			statusClass = "running"
		case "succeeded":
			statusClass = "ready"
		case "failed":
			statusClass = "notready"
		}

		result = append(result, models.Pod{
			Name:        p.Name,
			Namespace:   p.Namespace,
			Status:      string(p.Status.Phase),
			StatusClass: statusClass,
			Node:        p.Spec.NodeName,
			PodIP:       p.Status.PodIP,
			Restarts:    restarts,
		})
	}

	return result, nil
}

//revive:disable:cyclomatic
func (c Client) GetPod(namespace, name string) (*models.PodDetail, error) {
	pod, err := c.Clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed getting pod: %w", err)
	}

	statusClass := "pending"
	switch strings.ToLower(string(pod.Status.Phase)) {
	case "running":
		statusClass = "running"
	case "succeeded":
		statusClass = "ready"
	case "failed":
		statusClass = "notready"
	}

	// Build containers list
	var containers []models.PodContainer
	for _, c := range pod.Spec.Containers {
		container := models.PodContainer{
			Name:  c.Name,
			Image: c.Image,
		}

		// Ports
		for _, p := range c.Ports {
			container.Ports = append(container.Ports, models.ContainerPort{
				Name:          p.Name,
				ContainerPort: p.ContainerPort,
				Protocol:      string(p.Protocol),
			})
		}

		// Environment variables (only non-secret ones)
		for _, e := range c.Env {
			if e.ValueFrom == nil {
				container.Env = append(container.Env, models.EnvVar{
					Name:  e.Name,
					Value: e.Value,
				})
			}
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
			Limits: map[string]string{
				"cpu":    c.Resources.Limits.Cpu().String(),
				"memory": c.Resources.Limits.Memory().String(),
			},
			Requests: map[string]string{
				"cpu":    c.Resources.Requests.Cpu().String(),
				"memory": c.Resources.Requests.Memory().String(),
			},
		}

		// Get status from container statuses
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.Name == c.Name {
				container.Ready = cs.Ready
				container.RestartCount = int(cs.RestartCount)
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
	for _, c := range pod.Spec.InitContainers {
		container := models.PodContainer{
			Name:  c.Name,
			Image: c.Image,
		}

		for _, cs := range pod.Status.InitContainerStatuses {
			if cs.Name == c.Name {
				container.Ready = cs.Ready
				container.RestartCount = int(cs.RestartCount)
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
	for _, cond := range pod.Status.Conditions {
		conditions = append(conditions, models.PodCondition{
			Type:               string(cond.Type),
			Status:             string(cond.Status),
			Reason:             cond.Reason,
			Message:            cond.Message,
			LastTransitionTime: cond.LastTransitionTime.Format("2006-01-02T15:04:05Z"),
		})
	}

	// Build volumes list
	var volumes []models.PodVolume
	for _, v := range pod.Spec.Volumes {
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
	for _, or := range pod.OwnerReferences {
		ownerRefs = append(ownerRefs, models.OwnerReference{
			APIVersion: or.APIVersion,
			Kind:       or.Kind,
			Name:       or.Name,
			UID:        string(or.UID),
		})
	}

	return &models.PodDetail{
		Name:              pod.Name,
		Namespace:         pod.Namespace,
		UID:               string(pod.UID),
		CreationTimestamp: pod.CreationTimestamp.Format("2006-01-02T15:04:05Z"),
		Labels:            pod.Labels,
		Annotations:       pod.Annotations,
		Status:            string(pod.Status.Phase),
		StatusClass:       statusClass,
		Node:              pod.Spec.NodeName,
		PodIP:             pod.Status.PodIP,
		HostIP:            pod.Status.HostIP,
		QOSClass:          string(pod.Status.QOSClass),
		ServiceAccount:    pod.Spec.ServiceAccountName,
		Containers:        containers,
		InitContainers:    initContainers,
		Conditions:        conditions,
		Volumes:           volumes,
		OwnerReferences:   ownerRefs,
	}, nil
	//revive:enable:cyclomatic
}
