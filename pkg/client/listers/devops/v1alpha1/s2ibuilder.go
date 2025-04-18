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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/kubesphere/ks-devops/pkg/api/devops/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// S2iBuilderLister helps list S2iBuilders.
type S2iBuilderLister interface {
	// List lists all S2iBuilders in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.S2iBuilder, err error)
	// S2iBuilders returns an object that can list and get S2iBuilders.
	S2iBuilders(namespace string) S2iBuilderNamespaceLister
	S2iBuilderListerExpansion
}

// s2iBuilderLister implements the S2iBuilderLister interface.
type s2iBuilderLister struct {
	indexer cache.Indexer
}

// NewS2iBuilderLister returns a new S2iBuilderLister.
func NewS2iBuilderLister(indexer cache.Indexer) S2iBuilderLister {
	return &s2iBuilderLister{indexer: indexer}
}

// List lists all S2iBuilders in the indexer.
func (s *s2iBuilderLister) List(selector labels.Selector) (ret []*v1alpha1.S2iBuilder, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.S2iBuilder))
	})
	return ret, err
}

// S2iBuilders returns an object that can list and get S2iBuilders.
func (s *s2iBuilderLister) S2iBuilders(namespace string) S2iBuilderNamespaceLister {
	return s2iBuilderNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// S2iBuilderNamespaceLister helps list and get S2iBuilders.
type S2iBuilderNamespaceLister interface {
	// List lists all S2iBuilders in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.S2iBuilder, err error)
	// Get retrieves the S2iBuilder from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.S2iBuilder, error)
	S2iBuilderNamespaceListerExpansion
}

// s2iBuilderNamespaceLister implements the S2iBuilderNamespaceLister
// interface.
type s2iBuilderNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all S2iBuilders in the indexer for a given namespace.
func (s s2iBuilderNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.S2iBuilder, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.S2iBuilder))
	})
	return ret, err
}

// Get retrieves the S2iBuilder from the indexer for a given namespace and name.
func (s s2iBuilderNamespaceLister) Get(name string) (*v1alpha1.S2iBuilder, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("s2ibuilder"), name)
	}
	return obj.(*v1alpha1.S2iBuilder), nil
}
