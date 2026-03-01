package k8s

import (
	kaiClientset "github.com/NVIDIA/KAI-scheduler/pkg/apis/client/clientset/versioned"
	schedulingv2 "github.com/NVIDIA/KAI-scheduler/pkg/apis/client/clientset/versioned/typed/scheduling/v2"
	"k8s.io/client-go/rest"
)

func newKAIClients(cfg *rest.Config) (schedulingv2.SchedulingV2Interface, error) {
	cs, err := kaiClientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return cs.SchedulingV2(), nil
}
