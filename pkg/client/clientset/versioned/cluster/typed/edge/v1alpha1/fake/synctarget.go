//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by kcp code-generator. DO NOT EDIT.

package fake

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/testing"

	kcptesting "github.com/kcp-dev/client-go/third_party/k8s.io/client-go/testing"
	"github.com/kcp-dev/logicalcluster/v3"

	edgev1alpha1 "github.com/kubestellar/kubestellar/pkg/apis/edge/v1alpha1"
	edgev1alpha1client "github.com/kubestellar/kubestellar/pkg/client/clientset/versioned/typed/edge/v1alpha1"
)

var syncTargetsResource = schema.GroupVersionResource{Group: "edge.kubestellar.io", Version: "v1alpha1", Resource: "synctargets"}
var syncTargetsKind = schema.GroupVersionKind{Group: "edge.kubestellar.io", Version: "v1alpha1", Kind: "SyncTarget"}

type syncTargetsClusterClient struct {
	*kcptesting.Fake
}

// Cluster scopes the client down to a particular cluster.
func (c *syncTargetsClusterClient) Cluster(clusterPath logicalcluster.Path) edgev1alpha1client.SyncTargetInterface {
	if clusterPath == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return &syncTargetsClient{Fake: c.Fake, ClusterPath: clusterPath}
}

// List takes label and field selectors, and returns the list of SyncTargets that match those selectors across all clusters.
func (c *syncTargetsClusterClient) List(ctx context.Context, opts metav1.ListOptions) (*edgev1alpha1.SyncTargetList, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootListAction(syncTargetsResource, syncTargetsKind, logicalcluster.Wildcard, opts), &edgev1alpha1.SyncTargetList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &edgev1alpha1.SyncTargetList{ListMeta: obj.(*edgev1alpha1.SyncTargetList).ListMeta}
	for _, item := range obj.(*edgev1alpha1.SyncTargetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested SyncTargets across all clusters.
func (c *syncTargetsClusterClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(kcptesting.NewRootWatchAction(syncTargetsResource, logicalcluster.Wildcard, opts))
}

type syncTargetsClient struct {
	*kcptesting.Fake
	ClusterPath logicalcluster.Path
}

func (c *syncTargetsClient) Create(ctx context.Context, syncTarget *edgev1alpha1.SyncTarget, opts metav1.CreateOptions) (*edgev1alpha1.SyncTarget, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootCreateAction(syncTargetsResource, c.ClusterPath, syncTarget), &edgev1alpha1.SyncTarget{})
	if obj == nil {
		return nil, err
	}
	return obj.(*edgev1alpha1.SyncTarget), err
}

func (c *syncTargetsClient) Update(ctx context.Context, syncTarget *edgev1alpha1.SyncTarget, opts metav1.UpdateOptions) (*edgev1alpha1.SyncTarget, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootUpdateAction(syncTargetsResource, c.ClusterPath, syncTarget), &edgev1alpha1.SyncTarget{})
	if obj == nil {
		return nil, err
	}
	return obj.(*edgev1alpha1.SyncTarget), err
}

func (c *syncTargetsClient) UpdateStatus(ctx context.Context, syncTarget *edgev1alpha1.SyncTarget, opts metav1.UpdateOptions) (*edgev1alpha1.SyncTarget, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootUpdateSubresourceAction(syncTargetsResource, c.ClusterPath, "status", syncTarget), &edgev1alpha1.SyncTarget{})
	if obj == nil {
		return nil, err
	}
	return obj.(*edgev1alpha1.SyncTarget), err
}

func (c *syncTargetsClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.Invokes(kcptesting.NewRootDeleteActionWithOptions(syncTargetsResource, c.ClusterPath, name, opts), &edgev1alpha1.SyncTarget{})
	return err
}

func (c *syncTargetsClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := kcptesting.NewRootDeleteCollectionAction(syncTargetsResource, c.ClusterPath, listOpts)

	_, err := c.Fake.Invokes(action, &edgev1alpha1.SyncTargetList{})
	return err
}

func (c *syncTargetsClient) Get(ctx context.Context, name string, options metav1.GetOptions) (*edgev1alpha1.SyncTarget, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootGetAction(syncTargetsResource, c.ClusterPath, name), &edgev1alpha1.SyncTarget{})
	if obj == nil {
		return nil, err
	}
	return obj.(*edgev1alpha1.SyncTarget), err
}

// List takes label and field selectors, and returns the list of SyncTargets that match those selectors.
func (c *syncTargetsClient) List(ctx context.Context, opts metav1.ListOptions) (*edgev1alpha1.SyncTargetList, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootListAction(syncTargetsResource, syncTargetsKind, c.ClusterPath, opts), &edgev1alpha1.SyncTargetList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &edgev1alpha1.SyncTargetList{ListMeta: obj.(*edgev1alpha1.SyncTargetList).ListMeta}
	for _, item := range obj.(*edgev1alpha1.SyncTargetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

func (c *syncTargetsClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(kcptesting.NewRootWatchAction(syncTargetsResource, c.ClusterPath, opts))
}

func (c *syncTargetsClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*edgev1alpha1.SyncTarget, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootPatchSubresourceAction(syncTargetsResource, c.ClusterPath, name, pt, data, subresources...), &edgev1alpha1.SyncTarget{})
	if obj == nil {
		return nil, err
	}
	return obj.(*edgev1alpha1.SyncTarget), err
}
