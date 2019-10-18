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

package controller

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/cache"
	"knative.dev/eventing/pkg/reconciler"
	"knative.dev/pkg/controller"

	"github.com/antoineco/mqtt-event-source/apis/sources/v1alpha1"
	sourceslisters "github.com/antoineco/mqtt-event-source/client/generated/lister/sources/v1alpha1"
	"github.com/antoineco/mqtt-event-source/controller/errors"
)

// Reconciler reconciles MQTTSource resources.
type Reconciler struct {
	// wrapper for core controller components (clients, logger, ...)
	*reconciler.Base

	// listers index properties about resources
	mqttsourceLister sourceslisters.MQTTSourceLister
}

// Reconcile compares the actual state of a MQTTSource object referenced by key with its desired state, and attempts to
// converge the two.
func (r *Reconciler) Reconcile(ctx context.Context, key string) error {
	_, err := mqttsourceByKey(key, r.mqttsourceLister)
	if err != nil {
		return errors.Handle(err, ctx, "Failed to get object from local store")
	}

	return nil
}

// mqttsourceByKey retrieves a MQTTSource object from a lister by key (ns/name).
func mqttsourceByKey(key string, l sourceslisters.MQTTSourceLister) (*v1alpha1.MQTTSource, error) {
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		controller.NewPermanentError(pkgerrors.Wrap(err, "invalid object key"))
	}

	mqttSrc, err := l.MQTTSources(ns).Get(name)
	switch {
	case apierrors.IsNotFound(err):
		return nil, errors.NewSkippable(pkgerrors.Wrap(err, "object no longer exists"))
	case err != nil:
		return nil, err
	}

	return mqttSrc, nil
}
