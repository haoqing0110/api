// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
	v1alpha1 "open-cluster-management.io/api/client/cluster/clientset/versioned/typed/cluster/v1alpha1"
)

type FakeClusterV1alpha1 struct {
	*testing.Fake
}

func (c *FakeClusterV1alpha1) ClusterClaims() v1alpha1.ClusterClaimInterface {
	return &FakeClusterClaims{c}
}

func (c *FakeClusterV1alpha1) ManagedClusterScalars(namespace string) v1alpha1.ManagedClusterScalarInterface {
	return &FakeManagedClusterScalars{c, namespace}
}

func (c *FakeClusterV1alpha1) ManagedClusterSets() v1alpha1.ManagedClusterSetInterface {
	return &FakeManagedClusterSets{c}
}

func (c *FakeClusterV1alpha1) ManagedClusterSetBindings(namespace string) v1alpha1.ManagedClusterSetBindingInterface {
	return &FakeManagedClusterSetBindings{c, namespace}
}

func (c *FakeClusterV1alpha1) Placements(namespace string) v1alpha1.PlacementInterface {
	return &FakePlacements{c, namespace}
}

func (c *FakeClusterV1alpha1) PlacementDecisions(namespace string) v1alpha1.PlacementDecisionInterface {
	return &FakePlacementDecisions{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeClusterV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
