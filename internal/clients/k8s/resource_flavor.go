package k8s

import (
	"cmyk/internal/models"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kueuev1beta2 "sigs.k8s.io/kueue/apis/kueue/v1beta2"
)

func (c Client) ListResourceFlavors() ([]models.ResourceFlavor, error) {
	list, err := c.KueueClientset.KueueV1beta2().ResourceFlavors().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed listing resource flavors: %w", err)
	}

	var result []models.ResourceFlavor
	for _, rf := range list.Items {
		result = append(result, toResourceFlavorModel(&rf))
	}
	return result, nil
}

func (c Client) GetResourceFlavor(name string) (*models.ResourceFlavor, error) {
	rf, err := c.KueueClientset.KueueV1beta2().ResourceFlavors().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed getting resource flavor: %w", err)
	}

	result := toResourceFlavorModel(rf)
	return &result, nil
}

func (c Client) CreateResourceFlavor(rf models.ResourceFlavor) (*models.ResourceFlavor, error) {
	obj := &kueuev1beta2.ResourceFlavor{
		ObjectMeta: metav1.ObjectMeta{
			Name: rf.Name,
		},
		Spec: kueuev1beta2.ResourceFlavorSpec{
			NodeLabels: rf.NodeLabels,
		},
	}

	for _, t := range rf.NodeTaints {
		obj.Spec.NodeTaints = append(obj.Spec.NodeTaints, corev1.Taint{
			Key:    t.Key,
			Value:  t.Value,
			Effect: corev1.TaintEffect(t.Effect),
		})
	}

	for _, t := range rf.Tolerations {
		tol := corev1.Toleration{
			Effect:   corev1.TaintEffect(t.Effect),
			Key:      t.Key,
			Operator: corev1.TolerationOperator(t.Operator),
			Value:    t.Value,
		}
		if t.TolerationSeconds != nil {
			tol.TolerationSeconds = t.TolerationSeconds
		}
		obj.Spec.Tolerations = append(obj.Spec.Tolerations, tol)
	}

	if rf.TopologyName != "" {
		ref := kueuev1beta2.TopologyReference(rf.TopologyName)
		obj.Spec.TopologyName = &ref
	}

	created, err := c.KueueClientset.KueueV1beta2().ResourceFlavors().Create(context.TODO(), obj, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed creating resource flavor: %w", err)
	}

	result := toResourceFlavorModel(created)
	return &result, nil
}

func (c Client) DeleteResourceFlavor(name string) error {
	err := c.KueueClientset.KueueV1beta2().ResourceFlavors().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed deleting resource flavor: %w", err)
	}
	return nil
}

func toResourceFlavorModel(rf *kueuev1beta2.ResourceFlavor) models.ResourceFlavor {
	m := models.ResourceFlavor{
		Name:       rf.Name,
		NodeLabels: rf.Spec.NodeLabels,
	}

	for _, t := range rf.Spec.NodeTaints {
		m.NodeTaints = append(m.NodeTaints, models.NodeTaint{
			Key:    t.Key,
			Value:  t.Value,
			Effect: string(t.Effect),
		})
	}

	for _, t := range rf.Spec.Tolerations {
		tol := models.Toleration{
			Effect:   string(t.Effect),
			Key:      t.Key,
			Operator: string(t.Operator),
			Value:    t.Value,
		}
		if t.TolerationSeconds != nil {
			tol.TolerationSeconds = t.TolerationSeconds
		}
		m.Tolerations = append(m.Tolerations, tol)
	}

	if rf.Spec.TopologyName != nil {
		m.TopologyName = string(*rf.Spec.TopologyName)
	}

	return m
}
