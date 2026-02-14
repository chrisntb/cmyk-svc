package models

// ResourceFlavor represents a Kueue ResourceFlavor
type ResourceFlavor struct {
	Name         string            `json:"name"`
	NodeLabels   map[string]string `json:"nodeLabels,omitempty"`
	NodeTaints   []NodeTaint       `json:"nodeTaints,omitempty"`
	Tolerations  []Toleration      `json:"tolerations,omitempty"`
	TopologyName string            `json:"topologyName,omitempty"`
}

// Toleration represents a Kubernetes toleration
type Toleration struct {
	Effect            string `json:"effect,omitempty"`
	Key               string `json:"key,omitempty"`
	Operator          string `json:"operator,omitempty"`
	TolerationSeconds *int64 `json:"tolerationSeconds,omitempty"`
	Value             string `json:"value,omitempty"`
}
