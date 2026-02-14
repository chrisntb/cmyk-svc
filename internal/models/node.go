package models

// Node represents a Kubernetes node
type Node struct {
	CPU            string `json:"cpu"`
	IP             string `json:"ip"`
	KubeletVersion string `json:"kubeletVersion"`
	Memory         string `json:"memory"`
	Name           string `json:"name"`
	Ready          bool   `json:"ready"`
	Roles          string `json:"roles"`
}

// NodeAddress represents a node address
type NodeAddress struct {
	Address string `json:"address"`
	Type    string `json:"type"`
}

// NodeCondition represents a node condition
type NodeCondition struct {
	LastHeartbeatTime  string `json:"lastHeartbeatTime,omitempty"`
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	Message            string `json:"message,omitempty"`
	Reason             string `json:"reason,omitempty"`
	Status             string `json:"status"`
	Type               string `json:"type"`
}

// NodeDetail represents detailed Kubernetes node information
type NodeDetail struct {
	Addresses         []NodeAddress     `json:"addresses,omitempty"`
	Allocatable       NodeResources     `json:"allocatable,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	Capacity          NodeResources     `json:"capacity,omitempty"`
	Conditions        []NodeCondition   `json:"conditions,omitempty"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Images            []NodeImage       `json:"images,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Name              string            `json:"name"`
	NodeInfo          NodeSystemInfo    `json:"nodeInfo"`
	PodCIDR           string            `json:"podCIDR,omitempty"`
	Ready             bool              `json:"ready"`
	Roles             string            `json:"roles"`
	Taints            []NodeTaint       `json:"taints,omitempty"`
	UID               string            `json:"uid"`
}

// NodeImage represents a container image on a node
type NodeImage struct {
	Names     []string `json:"names,omitempty"`
	SizeBytes int64    `json:"sizeBytes,omitempty"`
}

// NodeResources represents resource capacity as a map of resource name to quantity
type NodeResources map[string]string

// NodeSystemInfo represents system information about a node
type NodeSystemInfo struct {
	Architecture            string `json:"architecture"`
	BootID                  string `json:"bootID"`
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
	KernelVersion           string `json:"kernelVersion"`
	KubeProxyVersion        string `json:"kubeProxyVersion"`
	KubeletVersion          string `json:"kubeletVersion"`
	MachineID               string `json:"machineID"`
	OperatingSystem         string `json:"operatingSystem"`
	OSImage                 string `json:"osImage"`
	SystemUUID              string `json:"systemUUID"`
}

// NodeTaint represents a taint on a node
type NodeTaint struct {
	Effect    string `json:"effect"`
	Key       string `json:"key"`
	TimeAdded string `json:"timeAdded,omitempty"`
	Value     string `json:"value,omitempty"`
}

// NodeList represents a list of K8s nodes from the API
type NodeList struct {
	Items []NodeItem `json:"items"`
}

// NodeItem represents a K8s node from the API
type NodeItem struct {
	Metadata struct {
		Annotations       map[string]string `json:"annotations,omitempty"`
		CreationTimestamp string            `json:"creationTimestamp"`
		Labels            map[string]string `json:"labels,omitempty"`
		Name              string            `json:"name"`
		UID               string            `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		PodCIDR string `json:"podCIDR,omitempty"`
		Taints  []struct {
			Effect    string `json:"effect"`
			Key       string `json:"key"`
			TimeAdded string `json:"timeAdded,omitempty"`
			Value     string `json:"value,omitempty"`
		} `json:"taints,omitempty"`
	} `json:"spec"`
	Status struct {
		Addresses []struct {
			Address string `json:"address"`
			Type    string `json:"type"`
		} `json:"addresses,omitempty"`
		Allocatable map[string]string `json:"allocatable,omitempty"`
		Capacity    map[string]string `json:"capacity,omitempty"`
		Conditions  []struct {
			LastHeartbeatTime  string `json:"lastHeartbeatTime,omitempty"`
			LastTransitionTime string `json:"lastTransitionTime,omitempty"`
			Message            string `json:"message,omitempty"`
			Reason             string `json:"reason,omitempty"`
			Status             string `json:"status"`
			Type               string `json:"type"`
		} `json:"conditions,omitempty"`
		Images []struct {
			Names     []string `json:"names,omitempty"`
			SizeBytes int64    `json:"sizeBytes,omitempty"`
		} `json:"images,omitempty"`
		NodeInfo struct {
			Architecture            string `json:"architecture"`
			BootID                  string `json:"bootID"`
			ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
			KernelVersion           string `json:"kernelVersion"`
			KubeProxyVersion        string `json:"kubeProxyVersion"`
			KubeletVersion          string `json:"kubeletVersion"`
			MachineID               string `json:"machineID"`
			OperatingSystem         string `json:"operatingSystem"`
			OSImage                 string `json:"osImage"`
			SystemUUID              string `json:"systemUUID"`
		} `json:"nodeInfo"`
	} `json:"status"`
}
