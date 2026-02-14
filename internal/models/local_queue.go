package models

// Condition represents a standard Kubernetes condition
type Condition struct {
	LastTransitionTime string `json:"lastTransitionTime"`
	Message            string `json:"message"`
	ObservedGeneration int64  `json:"observedGeneration,omitempty"`
	Reason             string `json:"reason"`
	Status             string `json:"status"`
	Type               string `json:"type"`
}

// FlavorUsage represents the usage of a flavor in a LocalQueue
type FlavorUsage struct {
	Name      string          `json:"name"`
	Resources []ResourceUsage `json:"resources"`
}

// LocalQueue represents a Kueue LocalQueue
type LocalQueue struct {
	AdmittedWorkloads  int32         `json:"admittedWorkloads"`
	ClusterQueue       string        `json:"clusterQueue"`
	Conditions         []Condition   `json:"conditions,omitempty"`
	FlavorsReservation []FlavorUsage `json:"flavorsReservation,omitempty"`
	FlavorsUsage       []FlavorUsage `json:"flavorsUsage,omitempty"`
	Name               string        `json:"name"`
	Namespace          string        `json:"namespace"`
	PendingWorkloads   int32         `json:"pendingWorkloads"`
	ReservingWorkloads int32         `json:"reservingWorkloads"`
	StopPolicy         string        `json:"stopPolicy,omitempty"`
}

// ResourceUsage represents the usage of a specific resource
type ResourceUsage struct {
	Name  string `json:"name"`
	Total string `json:"total,omitempty"`
}
