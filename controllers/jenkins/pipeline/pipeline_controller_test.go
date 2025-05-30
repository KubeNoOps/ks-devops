/*
Copyright 2020 KubeSphere Authors

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

package pipeline

import (
	"reflect"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"

	"github.com/golang/mock/gomock"
	fakeDevOps "github.com/kubesphere/ks-devops/pkg/client/devops/fake"
	"github.com/kubesphere/ks-devops/pkg/client/devops/jclient"

	"github.com/kubesphere/ks-devops/pkg/constants"

	modelsdevops "github.com/kubesphere/ks-devops/pkg/models/devops"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/diff"
	kubeinformers "k8s.io/client-go/informers"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"

	devops "github.com/kubesphere/ks-devops/pkg/api/devops/v1alpha3"

	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
	"github.com/kubesphere/ks-devops/pkg/client/clientset/versioned/fake"
	informers "github.com/kubesphere/ks-devops/pkg/client/informers/externalversions"
)

var (
	alwaysReady        = func() bool { return true }
	noResyncPeriodFunc = func() time.Duration { return 0 }
)

type fixture struct {
	t *testing.T

	client          *fake.Clientset
	kubeclient      *k8sfake.Clientset
	namespaceLister []*v1.Namespace
	pipelineLister  []*devops.Pipeline
	actions         []core.Action
	kubeactions     []core.Action

	kubeobjects []runtime.Object
	// Objects from here preloaded into NewSimpleFake.
	objects []runtime.Object
	// Objects from here preloaded into devops
	initDevOpsProject string
	initPipeline      []*devops.Pipeline
	expectPipeline    []*devops.Pipeline
	jenkinsClient     *jclient.JenkinsClient
	roundTripper      *mhttp.MockRoundTripper
}

func newFixture(t *testing.T) *fixture {
	f := &fixture{}
	f.t = t
	f.objects = []runtime.Object{}

	ctrl := gomock.NewController(t)
	roundTripper := mhttp.NewMockRoundTripper(ctrl)
	f.roundTripper = roundTripper
	f.jenkinsClient = fakeDevOps.NewFakeJenkinsClient(roundTripper)
	return f
}

func newNamespace(name string, projectName string) *v1.Namespace {
	ns := &v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: map[string]string{constants.DevOpsProjectLabelKey: projectName},
		},
	}
	TRUE := true
	ns.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion:         devops.GroupVersion.String(),
			Kind:               devops.ResourceKindDevOpsProject,
			Name:               projectName,
			BlockOwnerDeletion: &TRUE,
			Controller:         &TRUE,
		},
	}

	return ns
}

func newPipeline(namespace, name string, spec devops.PipelineSpec, withFinalizers bool, syncOk bool) *devops.Pipeline {
	pipeline := &devops.Pipeline{
		TypeMeta: metav1.TypeMeta{
			Kind:       devops.ResourceKindPipeline,
			APIVersion: devops.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   namespace,
			Name:        name,
			Annotations: map[string]string{},
		},
		Spec:   spec,
		Status: devops.PipelineStatus{},
	}
	if withFinalizers {
		pipeline.Finalizers = append(pipeline.Finalizers, devops.PipelineFinalizerName)
	}
	if syncOk {
		pipeline.Annotations[devops.PipelineSyncStatusAnnoKey] = modelsdevops.StatusSuccessful
	}
	return pipeline
}

func newDeletingPipeline(namespace, name string) *devops.Pipeline {
	now := metav1.Now()
	pipeline := &devops.Pipeline{
		TypeMeta: metav1.TypeMeta{
			Kind:       devops.ResourceKindPipeline,
			APIVersion: devops.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:         namespace,
			Name:              name,
			DeletionTimestamp: &now,
		},
	}
	pipeline.Finalizers = append(pipeline.Finalizers, devops.PipelineFinalizerName)

	return pipeline
}

func (f *fixture) newController() (*Controller, informers.SharedInformerFactory, kubeinformers.SharedInformerFactory, *fakeDevOps.Devops) {
	f.client = fake.NewSimpleClientset(f.objects...)
	f.kubeclient = k8sfake.NewSimpleClientset(f.kubeobjects...)

	i := informers.NewSharedInformerFactory(f.client, noResyncPeriodFunc())
	k8sI := kubeinformers.NewSharedInformerFactory(f.kubeclient, noResyncPeriodFunc())
	dI := fakeDevOps.NewWithPipelines(f.initDevOpsProject, f.initPipeline...)
	c := NewController(f.kubeclient, f.client, dI, k8sI.Core().V1().Namespaces(),
		i.Devops().V1alpha3().Pipelines())

	c.pipelineSynced = alwaysReady
	c.eventRecorder = &record.FakeRecorder{}

	for _, f := range f.pipelineLister {
		_ = i.Devops().V1alpha3().Pipelines().Informer().GetIndexer().Add(f)
	}

	for _, d := range f.namespaceLister {
		_ = k8sI.Core().V1().Namespaces().Informer().GetIndexer().Add(d)
	}

	return c, i, k8sI, dI
}

func (f *fixture) run(fooName string) {
	f.runController(fooName, true, false)
}

func (f *fixture) runController(projectName string, startInformers bool, expectError bool) {
	c, i, k8sI, dI := f.newController()
	if startInformers {
		stopCh := make(chan struct{})
		defer close(stopCh)
		i.Start(stopCh)
		k8sI.Start(stopCh)
	}

	err := c.syncHandler(projectName)
	if !expectError && err != nil {
		f.t.Errorf("error syncing foo: %v", err)
	} else if expectError && err == nil {
		f.t.Error("expected error syncing foo, got nil")
	}

	actions := filterInformerActions(f.client.Actions())
	k8sActions := filterInformerActions(f.kubeclient.Actions())
	for i, action := range k8sActions {
		if len(f.kubeactions) < i+1 {
			f.t.Errorf("%d unexpected actions: %+v", len(k8sActions)-len(f.kubeactions), k8sActions[i:])
			break
		}

		expectedAction := f.kubeactions[i]
		checkAction(expectedAction, action, f.t)
	}

	if len(f.kubeactions) > len(k8sActions) {
		f.t.Errorf("%d additional expected actions:%+v", len(f.kubeactions)-len(k8sActions), f.kubeactions[len(k8sActions):])
	}

	if len(f.actions) > len(actions) {
		f.t.Errorf("%d additional expected actions:%+v", len(f.actions)-len(actions), f.actions[len(actions):])
	}
	if len(dI.Pipelines[f.initDevOpsProject]) != len(f.expectPipeline) {
		f.t.Errorf(" unexpected objects: %v", dI.Projects)
	}
	for _, pipeline := range f.expectPipeline {
		actualPipeline := dI.Pipelines[f.initDevOpsProject][pipeline.Name]
		if !reflect.DeepEqual(actualPipeline, pipeline) {
			f.t.Errorf(" pipeline %+v not match %+v", pipeline, actualPipeline)
		}
	}
}

// checkAction verifies that expected and actual actions are equal and both have
// same attached resources
func checkAction(expected, actual core.Action, t *testing.T) {
	if !(expected.Matches(actual.GetVerb(), actual.GetResource().Resource) && actual.GetSubresource() == expected.GetSubresource()) {
		t.Errorf("Expected\n\t%#v\ngot\n\t%#v", expected, actual)
		return
	}

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Action has wrong type. Expected: %t. Got: %t", expected, actual)
		return
	}

	switch a := actual.(type) {
	case core.CreateActionImpl:
		e, _ := expected.(core.CreateActionImpl)
		expObject := e.GetObject()
		object := a.GetObject()

		if !reflect.DeepEqual(expObject, object) {
			t.Errorf("Action %s %s has wrong object\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintSideBySide(expObject, object))
		}
	case core.UpdateActionImpl:
		e, _ := expected.(core.UpdateActionImpl)
		expObject := e.GetObject()
		object := a.GetObject()

		if !reflect.DeepEqual(expObject, object) {
			t.Errorf("Action %s %s has wrong object\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintSideBySide(expObject, object))
		}
	case core.PatchActionImpl:
		e, _ := expected.(core.PatchActionImpl)
		expPatch := e.GetPatch()
		patch := a.GetPatch()

		if !reflect.DeepEqual(expPatch, patch) {
			t.Errorf("Action %s %s has wrong patch\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintSideBySide(expPatch, patch))
		}
	default:
		t.Errorf("Uncaptured Action %s %s, you should explicitly add a case to capture it",
			actual.GetVerb(), actual.GetResource().Resource)
	}
}

// filterInformerActions filters list and watch actions for testing resources.
// Since list and watch don't change resource state we can filter it to lower
// nose level in our tests.
func filterInformerActions(actions []core.Action) []core.Action {
	ret := []core.Action{}
	for _, action := range actions {
		if len(action.GetNamespace()) == 0 &&
			(action.Matches("list", devops.ResourcePluralPipeline) ||
				action.Matches("watch", devops.ResourcePluralPipeline) ||
				action.Matches("list", "namespaces") ||
				action.Matches("watch", "namespaces")) {
			continue
		}
		ret = append(ret, action)
	}

	return ret
}

func (f *fixture) expectUpdatePipelineAction(p *devops.Pipeline) {
	action := core.NewUpdateAction(schema.GroupVersionResource{
		Version:  devops.GroupVersion.Version,
		Resource: devops.ResourcePluralPipeline,
		Group:    devops.GroupVersion.Group,
	}, p.Namespace, p)
	f.actions = append(f.actions, action)
}

func getKey(p *devops.Pipeline, t *testing.T) string {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(p)
	if err != nil {
		t.Errorf("Unexpected error getting key for pipeline %v: %v", p.Name, err)
		return ""
	}
	return key
}

func TestDoNothing(t *testing.T) {
	f := newFixture(t)
	nsName := "test-123"
	pipelineName := "test"
	projectName := "test_project"

	ns := newNamespace(nsName, projectName)
	pipeline := newPipeline(nsName, pipelineName, devops.PipelineSpec{}, true, true)

	f.pipelineLister = append(f.pipelineLister, pipeline)
	f.namespaceLister = append(f.namespaceLister, ns)
	f.objects = append(f.objects, pipeline)
	f.initDevOpsProject = nsName
	f.initPipeline = []*devops.Pipeline{pipeline}
	f.expectPipeline = []*devops.Pipeline{pipeline}

	f.run(getKey(pipeline, t))
}

func TestAddPipelineFinalizers(t *testing.T) {
	f := newFixture(t)
	nsName := "test-123"
	pipelineName := "test"
	projectName := "test_project"

	ns := newNamespace(nsName, projectName)
	pipeline := newPipeline(nsName, pipelineName, devops.PipelineSpec{}, false, false)

	expectPipeline := newPipeline(nsName, pipelineName, devops.PipelineSpec{}, true, true)

	f.pipelineLister = append(f.pipelineLister, pipeline)
	f.namespaceLister = append(f.namespaceLister, ns)
	f.objects = append(f.objects, pipeline)
	f.initDevOpsProject = nsName
	f.initPipeline = []*devops.Pipeline{pipeline}
	f.expectPipeline = []*devops.Pipeline{pipeline}
	f.expectUpdatePipelineAction(expectPipeline)
	f.run(getKey(pipeline, t))
}

func TestCreatePipeline(t *testing.T) {
	f := newFixture(t)
	nsName := "test-123"
	pipelineName := "test"
	projectName := "test_project"
	spec := devops.PipelineSpec{
		Type: devops.NoScmPipelineType,
		Pipeline: &devops.NoScmPipeline{
			Name: pipelineName,
		},
	}
	pipeline := newPipeline(nsName, pipelineName, spec, false, false)
	ns := newNamespace(nsName, projectName)

	f.pipelineLister = append(f.pipelineLister, pipeline)
	f.namespaceLister = append(f.namespaceLister, ns)
	f.objects = append(f.objects, pipeline)
	f.initDevOpsProject = nsName

	expectPipeline := pipeline.DeepCopy()
	expectPipeline.Finalizers = []string{devops.PipelineFinalizerName}
	expectPipeline.Annotations = map[string]string{
		devops.PipelineSyncStatusAnnoKey: constants.StatusSuccessful,
	}
	f.expectPipeline = []*devops.Pipeline{expectPipeline}

	f.run(getKey(pipeline, t))
}

func TestDeletePipeline(t *testing.T) {
	f := newFixture(t)
	nsName := "test-123"
	pipelineName := "test"
	projectName := "test_project"

	ns := newNamespace(nsName, projectName)
	pipeline := newDeletingPipeline(nsName, pipelineName)

	expectPipeline := pipeline.DeepCopy()
	expectPipeline.Finalizers = []string{}
	f.pipelineLister = append(f.pipelineLister, pipeline)
	f.namespaceLister = append(f.namespaceLister, ns)
	f.objects = append(f.objects, pipeline)
	f.initDevOpsProject = nsName
	f.initPipeline = []*devops.Pipeline{pipeline}
	f.expectPipeline = []*devops.Pipeline{}
	f.expectUpdatePipelineAction(expectPipeline)
	f.run(getKey(pipeline, t))
}

func TestDeleteNotExistPipeline(t *testing.T) {
	f := newFixture(t)
	nsName := "test-123"
	pipelineName := "test"
	projectName := "test_project"

	ns := newNamespace(nsName, projectName)
	pipeline := newDeletingPipeline(nsName, pipelineName)

	expectPipeline := pipeline.DeepCopy()
	expectPipeline.Finalizers = []string{}
	f.pipelineLister = append(f.pipelineLister, pipeline)
	f.namespaceLister = append(f.namespaceLister, ns)
	f.objects = append(f.objects, pipeline)
	f.initDevOpsProject = nsName
	f.initPipeline = []*devops.Pipeline{}
	f.expectPipeline = []*devops.Pipeline{}
	f.expectUpdatePipelineAction(expectPipeline)
	f.run(getKey(pipeline, t))
}

func TestUpdatePipelineConfig(t *testing.T) {
	f := newFixture(t)
	nsName := "test-123"
	pipelineName := "test"
	projectName := "test_project"

	ns := newNamespace(nsName, projectName)
	initPipeline := newPipeline(nsName, pipelineName, devops.PipelineSpec{}, true, false)
	modifiedPipeline := newPipeline(nsName, pipelineName, devops.PipelineSpec{Type: "aa"}, true, false)
	expectPipeline := newPipeline(nsName, pipelineName, devops.PipelineSpec{Type: "aa"}, true, true)
	f.pipelineLister = append(f.pipelineLister, modifiedPipeline)
	f.namespaceLister = append(f.namespaceLister, ns)
	f.objects = append(f.objects, modifiedPipeline)
	f.initDevOpsProject = nsName
	f.initPipeline = []*devops.Pipeline{initPipeline}
	f.expectPipeline = []*devops.Pipeline{expectPipeline}
	f.run(getKey(modifiedPipeline, t))
}
