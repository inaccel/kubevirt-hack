package internal

import (
	"context"
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
)

type PodDefaulter struct{}

func NewPodDefaulter() *PodDefaulter {
	return &PodDefaulter{}
}

func (PodDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("pod defaulter did not understand object: %T", obj)
	}

	for i := range pod.Spec.Containers {
		if pod.Spec.Containers[i].Name == "compute" {
			for j := range pod.Spec.Containers[i].Command {
				if pod.Spec.Containers[i].Command[j] == "--hook-sidecars" {
					hookSidecars, err := strconv.Atoi(pod.Spec.Containers[i].Command[j+1])
					if err != nil {
						return err
					}
					hookSidecars++
					pod.Spec.Containers[i].Command[j+1] = strconv.Itoa(hookSidecars)
					break
				}
			}
			pod.Spec.Containers[i].SecurityContext.Capabilities.Add = append(pod.Spec.Containers[i].SecurityContext.Capabilities.Add, "SYS_RESOURCE")
			pod.Spec.Containers[i].VolumeMounts = append(pod.Spec.Containers[i].VolumeMounts, corev1.VolumeMount{
				Name:      "inaccel",
				MountPath: "/var/run/kubevirt-hooks/inaccel.sock",
			})
			break
		}
	}
	pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
		Name: "inaccel",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/run/kubevirt-hooks/inaccel.sock",
				Type: ptr.To(corev1.HostPathSocket),
			},
		},
	})

	return nil
}
