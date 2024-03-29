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

package objects

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/diff"

	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	sourcesv1alpha1 "github.com/antoineco/kyma-event-sources/apis/sources/v1alpha1"
)

func TestNewService(t *testing.T) {
	const (
		ns   = "testns"
		name = "test"
		img  = "registry/image:tag"
	)

	testMQTTSrc := &sourcesv1alpha1.MQTTSource{ObjectMeta: metav1.ObjectMeta{
		Namespace: ns,
		Name:      name,
		UID:       "00000000-0000-0000-0000-000000000000",
	}}

	testExistingKsvc := &servingv1.Service{ObjectMeta: metav1.ObjectMeta{
		Namespace:       ns,
		Name:            name,
		ResourceVersion: "1",
		Annotations: map[string]string{
			knativeServingAnnotations[0]: "some-user",
		},
	}}

	expectKsvc := &servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       ns,
			Name:            name,
			Annotations:     testExistingKsvc.Annotations,
			ResourceVersion: testExistingKsvc.ResourceVersion,
			OwnerReferences: []metav1.OwnerReference{
				*testMQTTSrc.ToOwner(),
			},
		},
		Spec: servingv1.ServiceSpec{
			ConfigurationSpec: servingv1.ConfigurationSpec{
				Template: servingv1.RevisionTemplateSpec{
					Spec: servingv1.RevisionSpec{
						PodSpec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Image: img,
							}},
						},
					},
				},
			},
		},
	}

	ksvc := NewService(ns, name,
		WithContainerImage(img),
		WithControllerRef(testMQTTSrc.ToOwner()),
		WithExisting(testExistingKsvc),
	)

	if d := diff.ObjectReflectDiff(expectKsvc, ksvc); d != "<no diffs>" {
		t.Errorf("Unexpected diff: (a:expect, b:got) %s", d)
	}
}
