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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/antoineco/mqtt-event-source/apis/sources/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMQTTSources implements MQTTSourceInterface
type FakeMQTTSources struct {
	Fake *FakeSourcesV1alpha1
	ns   string
}

var mqttsourcesResource = schema.GroupVersionResource{Group: "sources.kyma-project.io", Version: "v1alpha1", Resource: "mqttsources"}

var mqttsourcesKind = schema.GroupVersionKind{Group: "sources.kyma-project.io", Version: "v1alpha1", Kind: "MQTTSource"}

// Get takes name of the mQTTSource, and returns the corresponding mQTTSource object, and an error if there is any.
func (c *FakeMQTTSources) Get(name string, options v1.GetOptions) (result *v1alpha1.MQTTSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(mqttsourcesResource, c.ns, name), &v1alpha1.MQTTSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MQTTSource), err
}

// List takes label and field selectors, and returns the list of MQTTSources that match those selectors.
func (c *FakeMQTTSources) List(opts v1.ListOptions) (result *v1alpha1.MQTTSourceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(mqttsourcesResource, mqttsourcesKind, c.ns, opts), &v1alpha1.MQTTSourceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.MQTTSourceList{ListMeta: obj.(*v1alpha1.MQTTSourceList).ListMeta}
	for _, item := range obj.(*v1alpha1.MQTTSourceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested mQTTSources.
func (c *FakeMQTTSources) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(mqttsourcesResource, c.ns, opts))

}

// Create takes the representation of a mQTTSource and creates it.  Returns the server's representation of the mQTTSource, and an error, if there is any.
func (c *FakeMQTTSources) Create(mQTTSource *v1alpha1.MQTTSource) (result *v1alpha1.MQTTSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(mqttsourcesResource, c.ns, mQTTSource), &v1alpha1.MQTTSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MQTTSource), err
}

// Update takes the representation of a mQTTSource and updates it. Returns the server's representation of the mQTTSource, and an error, if there is any.
func (c *FakeMQTTSources) Update(mQTTSource *v1alpha1.MQTTSource) (result *v1alpha1.MQTTSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(mqttsourcesResource, c.ns, mQTTSource), &v1alpha1.MQTTSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MQTTSource), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeMQTTSources) UpdateStatus(mQTTSource *v1alpha1.MQTTSource) (*v1alpha1.MQTTSource, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(mqttsourcesResource, "status", c.ns, mQTTSource), &v1alpha1.MQTTSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MQTTSource), err
}

// Delete takes name of the mQTTSource and deletes it. Returns an error if one occurs.
func (c *FakeMQTTSources) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(mqttsourcesResource, c.ns, name), &v1alpha1.MQTTSource{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMQTTSources) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(mqttsourcesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.MQTTSourceList{})
	return err
}

// Patch applies the patch and returns the patched mQTTSource.
func (c *FakeMQTTSources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MQTTSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(mqttsourcesResource, c.ns, name, data, subresources...), &v1alpha1.MQTTSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MQTTSource), err
}