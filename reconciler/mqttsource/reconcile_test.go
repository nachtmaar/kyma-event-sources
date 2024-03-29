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

package mqttsource

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	k8stesting "k8s.io/client-go/testing"

	"knative.dev/eventing/pkg/duck"
	"knative.dev/eventing/pkg/reconciler"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	rt "knative.dev/pkg/reconciler/testing"
	fakeservingclient "knative.dev/serving/pkg/client/injection/client/fake"

	fakesourcesclient "github.com/antoineco/kyma-event-sources/client/generated/injection/client/fake"
	. "github.com/antoineco/kyma-event-sources/reconciler/testing"
)

const (
	tNs   = "testns"
	tName = "test"
	tImg  = "sources.kyma-project.io/mqtt:latest"
)

func TestReconcile(t *testing.T) {
	testCases := rt.TableTest{
		/* Error handling */

		{
			Name:    "Source was deleted",
			Key:     tNs + "/" + tName,
			Objects: nil,
			WantErr: false,
		},
		{
			Name:    "Invalid object key",
			Key:     tNs + "/" + tName + "/invalid",
			WantErr: true,
		},

		/* Service synchronization */

		{
			Name: "Initial source creation",
			Key:  tNs + "/" + tName,
			Objects: []runtime.Object{
				NewMQTTSource(tNs, tName),
			},
			WantCreates: []runtime.Object{
				NewService(tNs, tName,
					WithServiceController(tName),
					WithServiceContainer(tImg),
				),
			},
			WantUpdates: nil,
			WantStatusUpdates: []k8stesting.UpdateActionImpl{{
				Object: NewMQTTSource(tNs, tName,
					WithNotDeployed,
					WithSink,
				),
			}},
		},
		{
			Name: "Adapter Service up-to-date",
			Key:  tNs + "/" + tName,
			Objects: []runtime.Object{
				NewMQTTSource(tNs, tName,
					WithDeployed,
					WithSink,
				),
				NewService(tNs, tName,
					WithServiceContainer(tImg),
					WithServiceReady,
				),
			},
			WantCreates:       nil,
			WantUpdates:       nil,
			WantStatusUpdates: nil,
		},
		{
			Name: "Adapter Service spec does not match expectation",
			Key:  tNs + "/" + tName,
			Objects: []runtime.Object{
				NewMQTTSource(tNs, tName,
					WithDeployed,
					WithSink,
				),
				NewService(tNs, tName,
					WithServiceContainer("outdated"),
					WithServiceReady,
				),
			},
			WantCreates: nil,
			WantUpdates: []k8stesting.UpdateActionImpl{{
				Object: NewService(tNs, tName,
					WithServiceController(tName),
					WithServiceContainer(tImg),
					WithServiceReady),
			}},
			WantStatusUpdates: nil,
		},

		/* Status updates */

		{
			Name: "Adapter Service deployment in progress",
			Key:  tNs + "/" + tName,
			Objects: []runtime.Object{
				NewMQTTSource(tNs, tName,
					WithNotDeployed,
					WithSink,
				),
				NewService(tNs, tName,
					WithServiceContainer(tImg),
				),
			},
			WantCreates:       nil,
			WantUpdates:       nil,
			WantStatusUpdates: nil,
		},
		{
			Name: "Adapter Service became ready",
			Key:  tNs + "/" + tName,
			Objects: []runtime.Object{
				NewMQTTSource(tNs, tName,
					WithNotDeployed,
					WithSink,
				),
				NewService(tNs, tName,
					WithServiceContainer(tImg),
					WithServiceReady,
				),
			},
			WantCreates: nil,
			WantUpdates: nil,
			WantStatusUpdates: []k8stesting.UpdateActionImpl{{
				Object: NewMQTTSource(tNs, tName,
					WithDeployed,
					WithSink,
				),
			}},
		},
		{
			Name: "Adapter Service became not ready",
			Key:  tNs + "/" + tName,
			Objects: []runtime.Object{
				NewMQTTSource(tNs, tName,
					WithDeployed,
					WithSink,
				),
				NewService(tNs, tName,
					WithServiceContainer(tImg),
				),
			},
			WantCreates: nil,
			WantUpdates: nil,
			WantStatusUpdates: []k8stesting.UpdateActionImpl{{
				Object: NewMQTTSource(tNs, tName,
					WithNotDeployed,
					WithSink,
				),
			}},
		},
	}

	var ctor Ctor = func(ctx context.Context, ls *Listers, cmw configmap.Watcher) controller.Reconciler {
		return &Reconciler{
			Base:             reconciler.NewBase(ctx, controllerAgentName, cmw),
			mqttsourceLister: ls.GetMQTTSourceLister(),
			ksvcLister:       ls.GetServiceLister(),
			sourcesClient:    fakesourcesclient.Get(ctx).SourcesV1alpha1(),
			servingClient:    fakeservingclient.Get(ctx).ServingV1(),
			adapterImage:     tImg,
			sinkReconciler:   duck.NewSinkReconciler(ctx, func(string) {}),
		}
	}

	testCases.Test(t, MakeFactory(ctor))
}
