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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"

	spacev1alpha1 "github.com/kubestellar/kubestellar/space-framework/pkg/apis/space/v1alpha1"
	versioned "github.com/kubestellar/kubestellar/space-framework/pkg/client/clientset/versioned"
	internalinterfaces "github.com/kubestellar/kubestellar/space-framework/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kubestellar/kubestellar/space-framework/pkg/client/listers/space/v1alpha1"
)

// SpaceProviderDescInformer provides access to a shared informer and lister for
// SpaceProviderDescs.
type SpaceProviderDescInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.SpaceProviderDescLister
}

type spaceProviderDescInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewSpaceProviderDescInformer constructs a new informer for SpaceProviderDesc type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSpaceProviderDescInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSpaceProviderDescInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredSpaceProviderDescInformer constructs a new informer for SpaceProviderDesc type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSpaceProviderDescInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SpaceV1alpha1().SpaceProviderDescs().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SpaceV1alpha1().SpaceProviderDescs().Watch(context.TODO(), options)
			},
		},
		&spacev1alpha1.SpaceProviderDesc{},
		resyncPeriod,
		indexers,
	)
}

func (f *spaceProviderDescInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSpaceProviderDescInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *spaceProviderDescInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&spacev1alpha1.SpaceProviderDesc{}, f.defaultInformer)
}

func (f *spaceProviderDescInformer) Lister() v1alpha1.SpaceProviderDescLister {
	return v1alpha1.NewSpaceProviderDescLister(f.Informer().GetIndexer())
}
