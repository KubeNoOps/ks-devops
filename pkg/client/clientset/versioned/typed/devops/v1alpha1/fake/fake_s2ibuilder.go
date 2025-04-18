/*
Copyright 2020 The KubeSphere Authors.

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
	"context"

	v1alpha1 "github.com/kubesphere/ks-devops/pkg/api/devops/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeS2iBuilders implements S2iBuilderInterface
type FakeS2iBuilders struct {
	Fake *FakeDevopsV1alpha1
	ns   string
}

var s2ibuildersResource = schema.GroupVersionResource{Group: "devops.kubesphere.io", Version: "v1alpha1", Resource: "s2ibuilders"}

var s2ibuildersKind = schema.GroupVersionKind{Group: "devops.kubesphere.io", Version: "v1alpha1", Kind: "S2iBuilder"}

// Get takes name of the s2iBuilder, and returns the corresponding s2iBuilder object, and an error if there is any.
func (c *FakeS2iBuilders) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.S2iBuilder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(s2ibuildersResource, c.ns, name), &v1alpha1.S2iBuilder{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.S2iBuilder), err
}

// List takes label and field selectors, and returns the list of S2iBuilders that match those selectors.
func (c *FakeS2iBuilders) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.S2iBuilderList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(s2ibuildersResource, s2ibuildersKind, c.ns, opts), &v1alpha1.S2iBuilderList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.S2iBuilderList{ListMeta: obj.(*v1alpha1.S2iBuilderList).ListMeta}
	for _, item := range obj.(*v1alpha1.S2iBuilderList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested s2iBuilders.
func (c *FakeS2iBuilders) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(s2ibuildersResource, c.ns, opts))

}

// Create takes the representation of a s2iBuilder and creates it.  Returns the server's representation of the s2iBuilder, and an error, if there is any.
func (c *FakeS2iBuilders) Create(ctx context.Context, s2iBuilder *v1alpha1.S2iBuilder, opts v1.CreateOptions) (result *v1alpha1.S2iBuilder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(s2ibuildersResource, c.ns, s2iBuilder), &v1alpha1.S2iBuilder{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.S2iBuilder), err
}

// Update takes the representation of a s2iBuilder and updates it. Returns the server's representation of the s2iBuilder, and an error, if there is any.
func (c *FakeS2iBuilders) Update(ctx context.Context, s2iBuilder *v1alpha1.S2iBuilder, opts v1.UpdateOptions) (result *v1alpha1.S2iBuilder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(s2ibuildersResource, c.ns, s2iBuilder), &v1alpha1.S2iBuilder{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.S2iBuilder), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeS2iBuilders) UpdateStatus(ctx context.Context, s2iBuilder *v1alpha1.S2iBuilder, opts v1.UpdateOptions) (*v1alpha1.S2iBuilder, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(s2ibuildersResource, "status", c.ns, s2iBuilder), &v1alpha1.S2iBuilder{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.S2iBuilder), err
}

// Delete takes name of the s2iBuilder and deletes it. Returns an error if one occurs.
func (c *FakeS2iBuilders) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(s2ibuildersResource, c.ns, name), &v1alpha1.S2iBuilder{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeS2iBuilders) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(s2ibuildersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.S2iBuilderList{})
	return err
}

// Patch applies the patch and returns the patched s2iBuilder.
func (c *FakeS2iBuilders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.S2iBuilder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(s2ibuildersResource, c.ns, name, pt, data, subresources...), &v1alpha1.S2iBuilder{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.S2iBuilder), err
}
