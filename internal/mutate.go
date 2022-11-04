package internal

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type PodDefaulter struct{}

func (PodDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("pod defaulter did not understand object: %T", obj)
	}

	for i := range pod.Spec.Containers {
		if pod.Spec.Containers[i].Name == "compute" {
			pod.Spec.Containers[i].SecurityContext.Capabilities.Add = append(pod.Spec.Containers[i].SecurityContext.Capabilities.Add, "SYS_RESOURCE")
		}
	}

	return nil
}
