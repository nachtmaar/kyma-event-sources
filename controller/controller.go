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

// Package controller implements a controller for the MQTTSource custom resource.
package controller

import (
	"context"

	"knative.dev/eventing/pkg/reconciler"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"

	mqttsourceinformers "github.com/antoineco/mqtt-event-source/client/generated/injection/informers/sources/v1alpha1/mqttsource"
)

const (
	// reconcilerName is the name of the reconciler
	reconcilerName = "MQTTSources"

	// controllerAgentName is the string used by this controller to identify
	// itself when creating events.
	controllerAgentName = "mqttsource-controller"
)

// New returns a new controller that reconciles MQTTSource objects.
func New(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	mqttSourceInformer := mqttsourceinformers.Get(ctx)

	r := &Reconciler{
		Base:             reconciler.NewBase(ctx, controllerAgentName, cmw),
		mqttsourceLister: mqttSourceInformer.Lister(),
	}
	impl := controller.NewImpl(r, r.Logger, reconcilerName)

	// set event handler
	mqttSourceInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	return impl
}