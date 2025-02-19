/*
Copyright The KubeStellar Authors.

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

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"

	v1alpha1 "github.com/kubestellar/kubestellar/pkg/apis/edge/v1alpha1"
)

// FakeCustomizers implements CustomizerInterface
type FakeCustomizers struct {
	Fake *FakeEdgeV1alpha1
	ns   string
}

var customizersResource = schema.GroupVersionResource{Group: "edge.kubestellar.io", Version: "v1alpha1", Resource: "customizers"}

var customizersKind = schema.GroupVersionKind{Group: "edge.kubestellar.io", Version: "v1alpha1", Kind: "Customizer"}

// Get takes name of the customizer, and returns the corresponding customizer object, and an error if there is any.
func (c *FakeCustomizers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Customizer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(customizersResource, c.ns, name), &v1alpha1.Customizer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Customizer), err
}

// List takes label and field selectors, and returns the list of Customizers that match those selectors.
func (c *FakeCustomizers) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.CustomizerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(customizersResource, customizersKind, c.ns, opts), &v1alpha1.CustomizerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.CustomizerList{ListMeta: obj.(*v1alpha1.CustomizerList).ListMeta}
	for _, item := range obj.(*v1alpha1.CustomizerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested customizers.
func (c *FakeCustomizers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(customizersResource, c.ns, opts))

}

// Create takes the representation of a customizer and creates it.  Returns the server's representation of the customizer, and an error, if there is any.
func (c *FakeCustomizers) Create(ctx context.Context, customizer *v1alpha1.Customizer, opts v1.CreateOptions) (result *v1alpha1.Customizer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(customizersResource, c.ns, customizer), &v1alpha1.Customizer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Customizer), err
}

// Update takes the representation of a customizer and updates it. Returns the server's representation of the customizer, and an error, if there is any.
func (c *FakeCustomizers) Update(ctx context.Context, customizer *v1alpha1.Customizer, opts v1.UpdateOptions) (result *v1alpha1.Customizer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(customizersResource, c.ns, customizer), &v1alpha1.Customizer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Customizer), err
}

// Delete takes name of the customizer and deletes it. Returns an error if one occurs.
func (c *FakeCustomizers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(customizersResource, c.ns, name, opts), &v1alpha1.Customizer{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCustomizers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(customizersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.CustomizerList{})
	return err
}

// Patch applies the patch and returns the patched customizer.
func (c *FakeCustomizers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Customizer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(customizersResource, c.ns, name, pt, data, subresources...), &v1alpha1.Customizer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Customizer), err
}
