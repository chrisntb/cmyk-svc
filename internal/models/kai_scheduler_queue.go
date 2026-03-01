package models

// KaiSchedulerQueueResource represents the resource allocation for a single resource type
type KaiSchedulerQueueResource struct {
	Limit           float64 `json:"limit"`
	OverQuotaWeight float64 `json:"over_quota_weight,omitempty"`
	Quota           float64 `json:"quota"`
}

// KaiSchedulerQueueResources represents the resource allocations for a queue
type KaiSchedulerQueueResources struct {
	Cpu    KaiSchedulerQueueResource  `json:"cpu"`
	Gpu    KaiSchedulerQueueResource  `json:"gpu"`
	Memory *KaiSchedulerQueueResource `json:"memory,omitempty"`
}

// KaiSchedulerChildQueue represents a child queue in the kai scheduler hierarchy
type KaiSchedulerChildQueue struct {
	Name      string                     `json:"name"`
	Parent    string                     `json:"parent"`
	Resources KaiSchedulerQueueResources `json:"resources"`
}

// KaiSchedulerParentQueue represents a parent queue in the kai scheduler hierarchy
type KaiSchedulerParentQueue struct {
	ChildQueues []string                   `json:"child_queues,omitempty"`
	Name        string                     `json:"name"`
	Resources   KaiSchedulerQueueResources `json:"resources"`
}
