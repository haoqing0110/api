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

package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
	duckv1alpha1 "open-cluster-management.io/api/duck/v1alpha1"
)

const (
	threadsPerController = 2
)

func main() {
	var masterURL string
	flag.StringVar(&masterURL, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")

	var kubeconfig string
	flag.StringVar(&kubeconfig, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	kubeconfig = "/root/.kube/config"

	var resources gvkListFlag
	flag.Var(&resources, "resource", "The list of resources to operate over, in the form: Kind.version.group (e.g. Deployment.v1.app)")
	resources = []schema.GroupVersionKind{
		{
			Group:   "cluster.open-cluster-management.io",
			Version: "v1beta1",
			Kind:    "Placement",
		},
	}

	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	logger := logging.FromContext(context.TODO()).Named("controller")

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		logger.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		logger.Fatalf("Error building build clientset: %v", err)
	}

	resyncPeriod := 10 * time.Hour

	tif := &duck.TypedInformerFactory{
		Client:       dynamicClient,
		Type:         &duckv1alpha1.UpgradableType{},
		ResyncPeriod: resyncPeriod,
		StopChannel:  stopCh,
	}

	controllers := make([]*controller.Impl, 0, len(resources))

	for _, gvk := range resources {
		controllers = append(controllers, duckv1alpha1.NewController(
			logger, dynamicClient, tif, gvk))
	}

	// Start all of the controllers.
	for _, ctrlr := range controllers {
		go func(ctrlr *controller.Impl) {
			// We don't expect this to return until stop is called,
			// but if it does, propagate it back.
			if err := ctrlr.Run(context.Background()); err != nil {
				logger.Fatalf("Error running controller: %s", err.Error())
			}
		}(ctrlr)
	}

	<-stopCh
}

// Custom flag type for reading GroupVersionKind.
type gvkListFlag []schema.GroupVersionKind

func (i *gvkListFlag) String() string {
	strs := []string{}
	for _, x := range []schema.GroupVersionKind(*i) {
		strs = append(strs, x.String())
	}
	return strings.Join(strs, ",")
}

func (i *gvkListFlag) Set(value string) error {
	gvk, _ := schema.ParseKindArg(value)
	if gvk == nil {
		return fmt.Errorf("not a valid GroupVersionKind: %q", value)
	}
	*i = append(*i, *gvk)
	return nil
}
