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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"knative.dev/pkg/apis"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

const (
	// MQTTConditionSinkProvided has status True when the MQTTSource has
	// been configured with a sink target.
	MQTTConditionSinkProvided apis.ConditionType = "SinkProvided"

	// MQTTConditionReady has status True when the MQTTSource is ready to
	// send events.
	MQTTConditionReady = apis.ConditionReady
)

var mqttCondSet = apis.NewLivingConditionSet(
	MQTTConditionSinkProvided,
	MQTTConditionReady,
)

// MQTTSourceGVK returns a GroupVersionKind for the MQTTSource type.
func MQTTSourceGVK() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("MQTTSource")
}

// ToOwner return a OwnerReference corresponding to the given MQTTSource.
func (s *MQTTSource) ToOwner() *metav1.OwnerReference {
	return metav1.NewControllerRef(s, MQTTSourceGVK())
}

// ToKey returns a key corresponding to the MQTTSource.
func (s *MQTTSource) ToKey() string {
	return s.Namespace + "/" + s.Name
}

// ToDesc returns a description of the MQTTSource.
func (s *MQTTSource) ToDesc() string {
	return s.ToKey() + ", " + s.GroupVersionKind().String()
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (s *MQTTSourceStatus) InitializeConditions() {
	mqttCondSet.Manage(s).InitializeConditions()
}

// PropagateSinkProvided evaluates the provided sink URI and sets the sink
// provided condition.
func (s *MQTTSourceStatus) PropagateSinkProvided(uri string) {
	s.SinkURI = uri
	if uri == "" {
		mqttCondSet.Manage(s).MarkFalse(MQTTConditionSinkProvided,
			"NotFound", "The events sink was not found")
		return
	}
	mqttCondSet.Manage(s).MarkTrue(MQTTConditionSinkProvided)
}

// PropagateReady uses the readiness of the provided Knative Service to determine if
// the Ready condition should be marked as true or false.
func (s *MQTTSourceStatus) PropagateReady(ksvc *servingv1.Service) {
	if ksvc.Status.IsReady() {
		mqttCondSet.Manage(s).MarkTrue(MQTTConditionReady)
		return
	}

	condReady := ksvc.Status.GetCondition(MQTTConditionReady)
	mqttCondSet.Manage(s).MarkFalse(MQTTConditionReady, "KnativeServiceNotReady",
		"The adapter Service is not yet ready: %s", condReady.Message)
}
