package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/clients/dynamicclient"
	"knative.dev/pkg/logging"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type UpgradableType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +required
	Spec UpgradableSpec `json:"spec,omitempty"`
}

// UpgradableSpec contains Spec of the UpgradableType object
type UpgradableSpec struct {
	// +optional
	PlacementRefs []PlacementRefs `json:"placementRefs,omitempty"`
}

type PlacementRefs struct {
	PlacementRef `json:",inline"`
	// The rollout strategy to apply addon configurations change.
	// The rollout strategy only watches the addon configurations defined in ClusterManagementAddOn.
	// +kubebuilder:default={type: All}
	// +optional
	RolloutStrategy RolloutStrategy `json:"rolloutStrategy,omitempty"`
}

type PlacementRef struct {
	// Namespace is the namespace of the placement
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=1
	Namespace string `json:"namespace"`
	// Name is the name of the placement
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=1
	Name string `json:"name"`
}

// Rollout strategy to be used by workload applier controller.
// +k8s:deepcopy-gen=false
type RolloutStrategy struct {
	// Rollout strategy Types are All, Progressive and ProgressivePerGroup
	// 1) All means apply the workload to all clusters in the decision groups at once.
	// 2) Progressive means apply the workload to the selected clusters progressively per cluster. The workload will not be applied to the next cluster unless one of the current applied clusters reach the successful state or timeout.
	// 3) ProgressivePerGroup means apply the workload to decisionGroup clusters progressively per group. The workload will not be applied to the next decisionGroup unless all clusters in the current group reach the successful state or timeout.
	// +kubebuilder:default:=All
	// +optional
	Type RolloutStrategyType `json:"type,omitempty"`

	// List of the decision groups names or indexes to apply the workload first and fail if workload did not reach successful ((not proceed to apply workload to other decision groups/clusters).
	// +optional
	MandatoryDecisionGroups []MandatoryDecisionGroup `json:"mandatoryDecisionGroups,omitempty"`

	// Timeout define how long workload applier controller will wait till workload reach successful state in the cluster. Only considered for Rollout Type Progressive and ProgressivePerGroup.
	// Timeout default value is None meaning the workload applier will not proceed apply workload to other clusters if did not reach the successful state.
	// Timeout must be defined in [0-9h]|[0-9m]|[0-9s] format examples; 2h , 90m , 360s
	// +kubebuilder:validation:Pattern="^(([0-9])+[h|m|s])|None$"
	// +kubebuilder:default:=None
	// +optional
	Timeout string `json:"timeout,omitempty"`

	// MaxConcurrency is the max number of clusters to deploy workload concurrently. The default value for MaxConcurrency is determined from the clustersPerDecisionGroup defined in the placement->DecisionStrategy.
	// When RolloutStrategy Type defined as ProgressivePerGroup the MaxConcurrency must not be bigger than the clustersPerDecisionGroup defined in the placement->DecisionStrategy otherwise it is the workload applier controller responsibility to re-group the placement clusters and keep track each group status.
	// +kubebuilder:validation:Pattern="^((100|[0-9]{1,2})%|[0-9]+)$"
	// +optional
	MaxConcurrency string `json:"maxConcurrency,omitempty"`
}

// RolloutStrategy Type
type RolloutStrategyType string

const (
	//All means apply the workload to all clusters in the decision groups at once.
	All RolloutStrategyType = "All"
	//Progressive means apply the workload to the selected clusters progressively per cluster.
	Progressive RolloutStrategyType = "Progressive"
	//ProgressivePerGroup means apply the workload to the selected clusters progressively per group.
	ProgressivePerGroup RolloutStrategyType = "ProgressivePerGroup"
)

// MandatoryDecisionGroup set the decision group name and group index.
type MandatoryDecisionGroup struct {
	// GroupName of the decision group should match the placementDecisions label value with label key cluster.open-cluster-management.io/decision-group-name
	// +optional
	GroupName string `json:"groupName,omitempty"`

	// GroupIndex of the decision group should match the placementDecisions label value with label key cluster.open-cluster-management.io/decision-group-index
	// +optional
	GroupIndex int32 `json:"groupIndex,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UpgradableTypeList is a collection of UpgradableType.
type UpgradableTypeList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of AddOnPlacementScore
	Items []UpgradableTypeList `json:"items"`
}

// GetListType implements apis.Listable
func (r *UpgradableType) GetListType() runtime.Object {
	return &UpgradableType{}
}

// GetFullType implements duck.Implementable
func (*UpgradableType) GetFullType() duck.Populatable {
	return &UpgradableType{}
}

// Populate implements duck.Populatable
func (t *UpgradableType) Populate() {}

func init() {
	injection.Default.RegisterDuck(WithDuck)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct{}

func WithDuck(ctx context.Context) context.Context {
	dc := dynamicclient.Get(ctx)
	dif := &duck.CachedInformerFactory{
		Delegate: &duck.TypedInformerFactory{
			Client:       dc,
			Type:         (&UpgradableType{}).GetFullType(),
			ResyncPeriod: controller.GetResyncPeriod(ctx),
			StopChannel:  ctx.Done(),
		},
	}
	return context.WithValue(ctx, Key{}, dif)
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context) duck.InformerFactory {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch knative.dev/pkg/apis/duck.InformerFactory from context.")
	}
	return untyped.(duck.InformerFactory)
}
