/*
Copyright 2019 The Kyma Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testing

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/pkg/apis"
	"knative.dev/pkg/ptr"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	sourcesv1alpha1 "github.com/antoineco/kyma-event-sources/apis/sources/v1alpha1"
)

// NewService creates a Service object.
func NewService(ns, name string, opts ...ServiceOption) *servingv1.Service {
	s := &servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// ServiceOption is a functional option for Service objects.
type ServiceOption func(*servingv1.Service)

// WithServiceReady marks the Service as Ready.
func WithServiceReady(s *servingv1.Service) {
	s.Status.SetConditions(apis.Conditions{{
		Type:   apis.ConditionReady,
		Status: corev1.ConditionTrue,
	}})
}

func WithServiceController(srcName string) ServiceOption {
	return func(s *servingv1.Service) {
		gvk := sourcesv1alpha1.MQTTSourceGVK()

		s.OwnerReferences = []metav1.OwnerReference{{
			APIVersion:         gvk.GroupVersion().String(),
			Kind:               gvk.Kind,
			Name:               srcName,
			UID:                uid,
			Controller:         ptr.Bool(true),
			BlockOwnerDeletion: ptr.Bool(true),
		}}
	}
}

func WithServiceContainer(img string) ServiceOption {
	return func(s *servingv1.Service) {
		s.Spec.ConfigurationSpec.Template.Spec.PodSpec.Containers = []corev1.Container{{
			Image: img,
		}}
	}
}
