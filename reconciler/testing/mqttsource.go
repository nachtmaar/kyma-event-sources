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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	sourcesv1alpha1 "github.com/antoineco/kyma-event-sources/apis/sources/v1alpha1"
)

// TODO: the sink should be configurable
const tSinkURI = "http://event-publish-service.kyma-system.svc.cluster.local/"

const uid = types.UID("00000000-0000-0000-0000-000000000000")

func NewMQTTSource(ns, name string, opts ...MQTTSourceOption) *sourcesv1alpha1.MQTTSource {
	mqttSrc := &sourcesv1alpha1.MQTTSource{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
			UID:       uid,
		},
	}

	for _, opt := range opts {
		opt(mqttSrc)
	}

	return mqttSrc
}

type MQTTSourceOption func(s *sourcesv1alpha1.MQTTSource)

func WithDeployed(s *sourcesv1alpha1.MQTTSource) {
	s.Status.PropagateServiceReady(NewService("", "", WithServiceReady))
}

func WithNotDeployed(s *sourcesv1alpha1.MQTTSource) {
	s.Status.PropagateServiceReady(NewService("", ""))
}

func WithSink(s *sourcesv1alpha1.MQTTSource) {
	s.Status.MarkSink(tSinkURI)
}
