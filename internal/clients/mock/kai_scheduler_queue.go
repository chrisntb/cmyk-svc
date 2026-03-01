package mock

import (
	"cmyk/internal/models"
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
)

type rawKaiSchedulerResource struct {
	Limit           float64 `json:"limit"`
	OverQuotaWeight float64 `json:"overQuotaWeight"`
	Quota           float64 `json:"quota"`
}

type rawKaiSchedulerQueue struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Spec struct {
		ParentQueue string `json:"parentQueue"`
		Resources   struct {
			Cpu    rawKaiSchedulerResource  `json:"cpu"`
			Gpu    rawKaiSchedulerResource  `json:"gpu"`
			Memory *rawKaiSchedulerResource `json:"memory"`
		} `json:"resources"`
	} `json:"spec"`
	Status struct {
		ChildQueues []string `json:"childQueues"`
	} `json:"status"`
}

type rawKaiSchedulerQueueList struct {
	Items []rawKaiSchedulerQueue `json:"items"`
}

func toQueueResource(r rawKaiSchedulerResource) models.KaiSchedulerQueueResource {
	return models.KaiSchedulerQueueResource{
		Limit:           r.Limit,
		OverQuotaWeight: r.OverQuotaWeight,
		Quota:           r.Quota,
	}
}

func toQueueResources(r rawKaiSchedulerQueue) models.KaiSchedulerQueueResources {
	res := models.KaiSchedulerQueueResources{
		Cpu: toQueueResource(r.Spec.Resources.Cpu),
		Gpu: toQueueResource(r.Spec.Resources.Gpu),
	}
	if r.Spec.Resources.Memory != nil {
		m := toQueueResource(*r.Spec.Resources.Memory)
		res.Memory = &m
	}
	return res
}

func loadRawKaiSchedulerQueues() ([]rawKaiSchedulerQueue, error) {
	data, err := os.ReadFile("./internal/clients/mock/kai_scheduler_queues.json")
	if err != nil {
		return nil, err
	}

	var list rawKaiSchedulerQueueList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	return list.Items, nil
}

// ListKaiSchedulerParentQueues reads and returns mock parent queues (items without a parentQueue)
func (c Client) ListKaiSchedulerParentQueues() ([]models.KaiSchedulerParentQueue, error) {
	items, err := loadRawKaiSchedulerQueues()
	if err != nil {
		return nil, err
	}

	var result []models.KaiSchedulerParentQueue
	for _, item := range items {
		if item.Spec.ParentQueue != "" {
			continue
		}
		result = append(result, models.KaiSchedulerParentQueue{
			ChildQueues: item.Status.ChildQueues,
			Name:        item.Metadata.Name,
			Resources:   toQueueResources(item),
		})
	}

	return result, nil
}

// GetKaiSchedulerChildQueues reads and returns mock child queues for a given parent queue name
func (c Client) GetKaiSchedulerChildQueues(parent string) ([]models.KaiSchedulerChildQueue, error) {
	items, err := loadRawKaiSchedulerQueues()
	if err != nil {
		return nil, err
	}

	// Verify the parent exists
	parentFound := false
	for _, item := range items {
		if item.Metadata.Name == parent && item.Spec.ParentQueue == "" {
			parentFound = true
			break
		}
	}
	if !parentFound {
		return nil, fiber.ErrNotFound
	}

	var result []models.KaiSchedulerChildQueue
	for _, item := range items {
		if item.Spec.ParentQueue != parent {
			continue
		}
		result = append(result, models.KaiSchedulerChildQueue{
			Name:      item.Metadata.Name,
			Parent:    item.Spec.ParentQueue,
			Resources: toQueueResources(item),
		})
	}

	return result, nil
}
