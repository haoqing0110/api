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
	"context"
	"fmt"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/logging/logkey"
)

const controllerAgentName = "cachier-controller"

const annotationKey = "cachier.mattmoor.io/decorate"

// Reconciler is the controller implementation for PodSpecable resources
type Reconciler struct {
	// For reading the state of the world.
	lister cache.GenericLister
	//	imageLister cachinglisters.ImageLister

	// Sugared logger is easier to use but is not as performant as the
	// raw logger. In performance critical paths, call logger.Desugar()
	// and use the returned raw logger instead. In addition to the
	// performance benefits, raw logger also preserves type-safety at
	// the expense of slightly greater verbosity.
	Logger        *zap.SugaredLogger
	DynamicClient dynamic.Interface
}

// Check that we implement the controller.Reconciler interface.
var _ controller.Reconciler = (*Reconciler)(nil)

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
		DynamicClient: dynamicClient,
	}
	impl := controller.NewContext(context.Background(), r, controller.ControllerOptions{WorkQueueName: gvr.String(), Logger: logger})

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
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		logger.Errorf("invalid resource key: %s", key)
		return nil
	}

	// Get the thing resource with this namespace/name
	//untyped, err := c.lister.ByNamespace(namespace).Get(name)
	untyped, err := c.lister.Get(name)
	if errors.IsNotFound(err) {
		logger.Errorf("thing %q in work queue no longer exists", key)
		return nil
	} else if err != nil {
		return err
	}
	before := untyped.(*UpgradableType)
	klog.Warningf("HHQQ %v", before)

	after := before.DeepCopy()
	//	for i, _ := range after.Spec.PlacementRefs {
	//		after.Spec.PlacementRefs[i].Name = after.Spec.PlacementRefs[i].Name + "H"
	//	}

	patch, err := duck.CreatePatch(before, after)
	if err != nil {
		klog.Warningf("%v", err)
	}

	bytes, err := patch.MarshalJSON()
	if err != nil {
		klog.Warningf("%v", err)
	}

	gvk := after.GetGroupVersionKind()
	gvr, _ := meta.UnsafeGuessKindToResource(gvk)
	_, err = c.DynamicClient.Resource(gvr).Patch(ctx, after.Name, types.JSONPatchType,
		bytes, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("failed to apply to target %v", err)
	}
	return nil
}
