/*
Copyright 2017 The Kubernetes Authors.

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

package v1alpha1

import (

	// caching "github.com/knative/caching/pkg/apis/caching/v1alpha1"
	//cachingclientset "github.com/knative/caching/pkg/client/clientset/versioned"
	//cachinginformers "github.com/knative/caching/pkg/client/informers/externalversions/caching/v1alpha1"
	//cachinglisters "github.com/knative/caching/pkg/client/listers/caching/v1alpha1"

	"context"
	"fmt"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/logging/logkey"
)

const controllerAgentName = "rollout-controller"

//const annotationKey = "cachier.mattmoor.io/decorate"

// Reconciler is the controller implementation for PodSpecable resources
type Reconciler struct {
	// For reading the state of the world.
	lister cache.GenericLister

	Logger *zap.SugaredLogger
}

// Check that we implement the controller.Reconciler interface.
// var _ controller.Reconciler = (*Reconciler)(nil)

// NewController returns a new PodSpecable controller
func NewController(
	logger *zap.SugaredLogger,
	dynamicClient dynamic.Interface,
	psif duck.InformerFactory,
	gvk schema.GroupVersionKind,
) *controller.Impl {

	// GVK => GVR
	gvr, _ := meta.UnsafeGuessKindToResource(gvk)

	// Get an informer / lister pair for this resource group.
	informer, lister, err := psif.Get(context.Background(), gvr)
	if err != nil {
		logger.Fatalf("Error building informer for %v: %v", gvr, err)
	}

	r := &Reconciler{
		lister: lister,
		Logger: logger.Named(controllerAgentName).
			With(zap.String(logkey.ControllerType, controllerAgentName)),
	}
	impl := controller.NewContext(context.Background(), r, controller.ControllerOptions{})

	r.Logger.Info("Setting up event handlers")

	// As resources in the tracked resource group change, have our informer
	// queue those resources for reconciliation.
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    impl.Enqueue,
		UpdateFunc: controller.PassNew(impl.Enqueue),
	})

	return impl
}

// Reconcile implements controller.Reconciler
func (c *Reconciler) Reconcile(ctx context.Context, key string) error {
	logger := logging.FromContext(ctx)
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		logger.Errorf("invalid resource key: %s", key)
		return nil
	}

	// Get the thing resource with this namespace/name
	untyped, err := c.lister.ByNamespace(namespace).Get(name)
	if errors.IsNotFound(err) {
		logger.Errorf("thing %q in work queue no longer exists", key)
		return nil
	} else if err != nil {
		return err
	}
	thing := untyped.(*UpgradableType)
	fmt.Printf("%v", thing)

	return nil
}
