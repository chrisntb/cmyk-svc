package k8s

import (
	"cmyk/internal/models"
	"context"
	"fmt"

	kaiSchedulingV2 "github.com/NVIDIA/KAI-scheduler/pkg/apis/scheduling/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) ListKaiSchedulerParentQueues() ([]models.KaiSchedulerParentQueue, error) {
	list, err := c.KAISchedulerClient.Queues("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed listing kai scheduler queues: %w", err)
	}

	var result []models.KaiSchedulerParentQueue
	for _, q := range list.Items {
		if q.Spec.ParentQueue != "" {
			continue
		}
		result = append(result, toKaiSchedulerParentQueue(q))
	}

	return result, nil
}

func (c Client) GetKaiSchedulerChildQueues(parent string) ([]models.KaiSchedulerChildQueue, error) {
	list, err := c.KAISchedulerClient.Queues("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed listing kai scheduler queues: %w", err)
	}

	parentFound := false
	var result []models.KaiSchedulerChildQueue
	for _, q := range list.Items {
		if q.Spec.ParentQueue == "" && q.Name == parent {
			parentFound = true
			continue
		}
		if q.Spec.ParentQueue == parent {
			result = append(result, toKaiSchedulerChildQueue(q))
		}
	}

	if !parentFound {
		return nil, apierrors.NewNotFound(kaiSchedulingV2.Resource("queues"), parent)
	}

	return result, nil
}

func toKaiSchedulerParentQueue(q kaiSchedulingV2.Queue) models.KaiSchedulerParentQueue {
	return models.KaiSchedulerParentQueue{
		ChildQueues: q.Status.ChildQueues,
		Name:        q.Name,
		Resources:   toKaiSchedulerQueueResources(q.Spec.Resources),
	}
}

func toKaiSchedulerChildQueue(q kaiSchedulingV2.Queue) models.KaiSchedulerChildQueue {
	return models.KaiSchedulerChildQueue{
		Name:      q.Name,
		Parent:    q.Spec.ParentQueue,
		Resources: toKaiSchedulerQueueResources(q.Spec.Resources),
	}
}

func toKaiSchedulerQueueResources(r *kaiSchedulingV2.QueueResources) models.KaiSchedulerQueueResources {
	if r == nil {
		return models.KaiSchedulerQueueResources{}
	}
	mem := toKaiSchedulerQueueResource(r.Memory)
	return models.KaiSchedulerQueueResources{
		Cpu:    toKaiSchedulerQueueResource(r.CPU),
		Gpu:    toKaiSchedulerQueueResource(r.GPU),
		Memory: &mem,
	}
}

func toKaiSchedulerQueueResource(r kaiSchedulingV2.QueueResource) models.KaiSchedulerQueueResource {
	return models.KaiSchedulerQueueResource{
		Limit:           r.Limit,
		OverQuotaWeight: r.OverQuotaWeight,
		Quota:           r.Quota,
	}
}
