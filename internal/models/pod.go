package models

// ContainerPort represents a port exposed by a container
type ContainerPort struct {
	ContainerPort int32  `json:"containerPort"`
	HostIP        string `json:"hostIP,omitempty"`
	HostPort      int32  `json:"hostPort,omitempty"`
	Name          string `json:"name,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

// OwnerReference represents an owner reference
type OwnerReference struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	UID        string `json:"uid"`
}

// Pod represents a Kubernetes pod
type Pod struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Node        string `json:"node"`
	PodIP       string `json:"podIP"`
	Restarts    int    `json:"restarts"`
	Status      string `json:"status"`
	StatusClass string `json:"statusClass"`
}

// PodCondition represents a pod condition
type PodCondition struct {
	LastProbeTime      string `json:"lastProbeTime,omitempty"`
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	Message            string `json:"message,omitempty"`
	Reason             string `json:"reason,omitempty"`
	Status             string `json:"status"`
	Type               string `json:"type"`
}

// PodContainer represents a container in a pod
type PodContainer struct {
	Env          []EnvVar             `json:"env,omitempty"`
	Image        string               `json:"image"`
	Name         string               `json:"name"`
	Ports        []ContainerPort      `json:"ports,omitempty"`
	Ready        bool                 `json:"ready"`
	Resources    ResourceRequirements `json:"resources,omitempty"`
	RestartCount int                  `json:"restartCount"`
	State        string               `json:"state,omitempty"`
	StateReason  string               `json:"stateReason,omitempty"`
	VolumeMounts []VolumeMount        `json:"volumeMounts,omitempty"`
}

// PodDetail represents detailed Kubernetes pod information
type PodDetail struct {
	Annotations       map[string]string `json:"annotations,omitempty"`
	Conditions        []PodCondition    `json:"conditions,omitempty"`
	Containers        []PodContainer    `json:"containers,omitempty"`
	CreationTimestamp string            `json:"creationTimestamp"`
	HostIP            string            `json:"hostIP,omitempty"`
	InitContainers    []PodContainer    `json:"initContainers,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Node              string            `json:"node,omitempty"`
	OwnerReferences   []OwnerReference  `json:"ownerReferences,omitempty"`
	PodIP             string            `json:"podIP,omitempty"`
	QOSClass          string            `json:"qosClass,omitempty"`
	ServiceAccount    string            `json:"serviceAccount,omitempty"`
	Status            string            `json:"status"`
	StatusClass       string            `json:"statusClass"`
	UID               string            `json:"uid"`
	Volumes           []PodVolume       `json:"volumes,omitempty"`
}

// PodVolume represents a volume in a pod
type PodVolume struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Type   string `json:"type"`
}

// ResourceRequirements represents resource requests and limits
type ResourceRequirements struct {
	Limits   map[string]string `json:"limits,omitempty"`
	Requests map[string]string `json:"requests,omitempty"`
}

// VolumeMount represents a volume mount in a container
type VolumeMount struct {
	MountPath string `json:"mountPath"`
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
	SubPath   string `json:"subPath,omitempty"`
}

// PodList represents a list of K8s pods from the API
type PodList struct {
	Items []PodItem `json:"items"`
}

// PodItem represents a K8s pod from the API
type PodItem struct {
	Metadata struct {
		Annotations       map[string]string `json:"annotations,omitempty"`
		CreationTimestamp string            `json:"creationTimestamp"`
		Labels            map[string]string `json:"labels,omitempty"`
		Name              string            `json:"name"`
		Namespace         string            `json:"namespace"`
		OwnerReferences   []struct {
			APIVersion string `json:"apiVersion"`
			Kind       string `json:"kind"`
			Name       string `json:"name"`
			UID        string `json:"uid"`
		} `json:"ownerReferences,omitempty"`
		UID string `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		Containers []struct {
			Env []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"env,omitempty"`
			Image string `json:"image"`
			Name  string `json:"name"`
			Ports []struct {
				ContainerPort int32  `json:"containerPort"`
				Name          string `json:"name,omitempty"`
				Protocol      string `json:"protocol,omitempty"`
			} `json:"ports,omitempty"`
			Resources struct {
				Limits   map[string]string `json:"limits,omitempty"`
				Requests map[string]string `json:"requests,omitempty"`
			} `json:"resources,omitempty"`
			VolumeMounts []struct {
				MountPath string `json:"mountPath"`
				Name      string `json:"name"`
				ReadOnly  bool   `json:"readOnly,omitempty"`
			} `json:"volumeMounts,omitempty"`
		} `json:"containers"`
		InitContainers []struct {
			Image string `json:"image"`
			Name  string `json:"name"`
		} `json:"initContainers,omitempty"`
		NodeName           string `json:"nodeName,omitempty"`
		ServiceAccountName string `json:"serviceAccountName,omitempty"`
		Volumes            []struct {
			ConfigMap *struct {
				Name string `json:"name"`
			} `json:"configMap,omitempty"`
			EmptyDir *struct{} `json:"emptyDir,omitempty"`
			HostPath *struct {
				Path string `json:"path"`
			} `json:"hostPath,omitempty"`
			Name                  string `json:"name"`
			PersistentVolumeClaim *struct {
				ClaimName string `json:"claimName"`
			} `json:"persistentVolumeClaim,omitempty"`
			Projected *struct{} `json:"projected,omitempty"`
			Secret    *struct {
				SecretName string `json:"secretName"`
			} `json:"secret,omitempty"`
		} `json:"volumes,omitempty"`
	} `json:"spec"`
	Status struct {
		Conditions []struct {
			LastProbeTime      string `json:"lastProbeTime,omitempty"`
			LastTransitionTime string `json:"lastTransitionTime,omitempty"`
			Message            string `json:"message,omitempty"`
			Reason             string `json:"reason,omitempty"`
			Status             string `json:"status"`
			Type               string `json:"type"`
		} `json:"conditions,omitempty"`
		ContainerStatuses []struct {
			Name         string `json:"name"`
			Ready        bool   `json:"ready"`
			RestartCount int    `json:"restartCount"`
			State        struct {
				Running *struct {
					StartedAt string `json:"startedAt"`
				} `json:"running,omitempty"`
				Terminated *struct {
					Reason string `json:"reason"`
				} `json:"terminated,omitempty"`
				Waiting *struct {
					Reason string `json:"reason"`
				} `json:"waiting,omitempty"`
			} `json:"state"`
		} `json:"containerStatuses,omitempty"`
		HostIP                string `json:"hostIP,omitempty"`
		InitContainerStatuses []struct {
			Name         string `json:"name"`
			Ready        bool   `json:"ready"`
			RestartCount int    `json:"restartCount"`
			State        struct {
				Running *struct {
					StartedAt string `json:"startedAt"`
				} `json:"running,omitempty"`
				Terminated *struct {
					Reason string `json:"reason"`
				} `json:"terminated,omitempty"`
				Waiting *struct {
					Reason string `json:"reason"`
				} `json:"waiting,omitempty"`
			} `json:"state"`
		} `json:"initContainerStatuses,omitempty"`
		Phase    string `json:"phase,omitempty"`
		PodIP    string `json:"podIP,omitempty"`
		QOSClass string `json:"qosClass,omitempty"`
	} `json:"status"`
}
