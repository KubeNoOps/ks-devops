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

// xCode generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	devopsv1alpha1 "github.com/kubesphere/ks-devops/pkg/api/devops/v1alpha1"
	versioned "github.com/kubesphere/ks-devops/pkg/client/clientset/versioned"
	internalinterfaces "github.com/kubesphere/ks-devops/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kubesphere/ks-devops/pkg/client/listers/devops/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// S2iRunInformer provides access to a shared informer and lister for
// S2iRuns.
type S2iRunInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.S2iRunLister
}

type s2iRunInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewS2iRunInformer constructs a new informer for S2iRun type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewS2iRunInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredS2iRunInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredS2iRunInformer constructs a new informer for S2iRun type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredS2iRunInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.DevopsV1alpha1().S2iRuns(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.DevopsV1alpha1().S2iRuns(namespace).Watch(context.TODO(), options)
			},
		},
		&devopsv1alpha1.S2iRun{},
		resyncPeriod,
		indexers,
	)
}

func (f *s2iRunInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredS2iRunInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *s2iRunInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&devopsv1alpha1.S2iRun{}, f.defaultInformer)
}

func (f *s2iRunInformer) Lister() v1alpha1.S2iRunLister {
	return v1alpha1.NewS2iRunLister(f.Informer().GetIndexer())
}
