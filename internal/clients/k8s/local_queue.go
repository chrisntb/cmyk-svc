package k8s

import (
	"cmyk/internal/models"
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kueuev1beta2 "sigs.k8s.io/kueue/apis/kueue/v1beta2"
)

func (c Client) ListLocalQueues() ([]models.LocalQueue, error) {
	list, err := c.KueueClientset.KueueV1beta2().LocalQueues(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed listing local queues: %w", err)
	}

	var result []models.LocalQueue
	for _, lq := range list.Items {
		result = append(result, toLocalQueueModel(&lq))
	}
	return result, nil
}

func (c Client) GetLocalQueue(namespace, name string) (*models.LocalQueue, error) {
	lq, err := c.KueueClientset.KueueV1beta2().LocalQueues(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed getting local queue: %w", err)
	}

	result := toLocalQueueModel(lq)
	return &result, nil
}

func (c Client) CreateLocalQueue(namespace string, lq models.LocalQueue) (*models.LocalQueue, error) {
	obj := &kueuev1beta2.LocalQueue{
		ObjectMeta: metav1.ObjectMeta{
			Name:      lq.Name,
			Namespace: namespace,
		},
		Spec: kueuev1beta2.LocalQueueSpec{
			ClusterQueue: kueuev1beta2.ClusterQueueReference(lq.ClusterQueue),
		},
	}

	if lq.StopPolicy != "" {
		sp := kueuev1beta2.StopPolicy(lq.StopPolicy)
		obj.Spec.StopPolicy = &sp
	}

	created, err := c.KueueClientset.KueueV1beta2().LocalQueues(namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed creating local queue: %w", err)
	}

	result := toLocalQueueModel(created)
	return &result, nil
}

func (c Client) DeleteLocalQueue(namespace, name string) error {
	err := c.KueueClientset.KueueV1beta2().LocalQueues(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed deleting local queue: %w", err)
	}
	return nil
}

func toLocalQueueModel(lq *kueuev1beta2.LocalQueue) models.LocalQueue {
	m := models.LocalQueue{
		AdmittedWorkloads:  lq.Status.AdmittedWorkloads,
		ClusterQueue:       string(lq.Spec.ClusterQueue),
		Name:               lq.Name,
		Namespace:          lq.Namespace,
		PendingWorkloads:   lq.Status.PendingWorkloads,
		ReservingWorkloads: lq.Status.ReservingWorkloads,
	}

	if lq.Spec.StopPolicy != nil {
		m.StopPolicy = string(*lq.Spec.StopPolicy)
	}

	for _, c := range lq.Status.Conditions {
		m.Conditions = append(m.Conditions, models.Condition{
			LastTransitionTime: c.LastTransitionTime.Format("2006-01-02T15:04:05Z"),
			Message:            c.Message,
			ObservedGeneration: c.ObservedGeneration,
			Reason:             c.Reason,
			Status:             string(c.Status),
			Type:               c.Type,
		})
	}

	for _, f := range lq.Status.FlavorsReservation {
		fu := models.FlavorUsage{
			Name: string(f.Name),
		}
		for _, r := range f.Resources {
			fu.Resources = append(fu.Resources, models.ResourceUsage{
				Name:  string(r.Name),
				Total: r.Total.String(),
			})
		}
		m.FlavorsReservation = append(m.FlavorsReservation, fu)
	}

	for _, f := range lq.Status.FlavorsUsage {
		fu := models.FlavorUsage{
			Name: string(f.Name),
		}
		for _, r := range f.Resources {
			fu.Resources = append(fu.Resources, models.ResourceUsage{
				Name:  string(r.Name),
				Total: r.Total.String(),
			})
		}
		m.FlavorsUsage = append(m.FlavorsUsage, fu)
	}

	return m
}
