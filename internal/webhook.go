package internal

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var Webhook = admission.WithCustomDefaulter(new(corev1.Pod), PodDefaulter{})
