package v1alpha1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE
var map_Key = map[string]string{
	"": "Key is used for associating the Informer inside the context.Context.",
}

func (Key) SwaggerDoc() map[string]string {
	return map_Key
}

var map_MandatoryDecisionGroup = map[string]string{
	"":           "MandatoryDecisionGroup set the decision group name and group index.",
	"groupName":  "GroupName of the decision group should match the placementDecisions label value with label key cluster.open-cluster-management.io/decision-group-name",
	"groupIndex": "GroupIndex of the decision group should match the placementDecisions label value with label key cluster.open-cluster-management.io/decision-group-index",
}

func (MandatoryDecisionGroup) SwaggerDoc() map[string]string {
	return map_MandatoryDecisionGroup
}

var map_PlacementRef = map[string]string{
	"namespace": "Namespace is the namespace of the placement",
	"name":      "Name is the name of the placement",
}

func (PlacementRef) SwaggerDoc() map[string]string {
	return map_PlacementRef
}

var map_PlacementRefs = map[string]string{
	"rolloutStrategy": "The rollout strategy to apply addon configurations change. The rollout strategy only watches the addon configurations defined in ClusterManagementAddOn.",
}

func (PlacementRefs) SwaggerDoc() map[string]string {
	return map_PlacementRefs
}

var map_RolloutStrategy = map[string]string{
	"":                        "Rollout strategy to be used by workload applier controller.",
	"type":                    "Rollout strategy Types are All, Progressive and ProgressivePerGroup 1) All means apply the workload to all clusters in the decision groups at once. 2) Progressive means apply the workload to the selected clusters progressively per cluster. The workload will not be applied to the next cluster unless one of the current applied clusters reach the successful state or timeout. 3) ProgressivePerGroup means apply the workload to decisionGroup clusters progressively per group. The workload will not be applied to the next decisionGroup unless all clusters in the current group reach the successful state or timeout.",
	"mandatoryDecisionGroups": "List of the decision groups names or indexes to apply the workload first and fail if workload did not reach successful ((not proceed to apply workload to other decision groups/clusters).",
	"timeout":                 "Timeout define how long workload applier controller will wait till workload reach successful state in the cluster. Only considered for Rollout Type Progressive and ProgressivePerGroup. Timeout default value is None meaning the workload applier will not proceed apply workload to other clusters if did not reach the successful state. Timeout must be defined in [0-9h]|[0-9m]|[0-9s] format examples; 2h , 90m , 360s",
	"maxConcurrency":          "MaxConcurrency is the max number of clusters to deploy workload concurrently. The default value for MaxConcurrency is determined from the clustersPerDecisionGroup defined in the placement->DecisionStrategy. When RolloutStrategy Type defined as ProgressivePerGroup the MaxConcurrency must not be bigger than the clustersPerDecisionGroup defined in the placement->DecisionStrategy otherwise it is the workload applier controller responsibility to re-group the placement clusters and keep track each group status.",
}

func (RolloutStrategy) SwaggerDoc() map[string]string {
	return map_RolloutStrategy
}

var map_UpgradableSpec = map[string]string{
	"": "UpgradableSpec contains Spec of the UpgradableType object",
}

func (UpgradableSpec) SwaggerDoc() map[string]string {
	return map_UpgradableSpec
}

var map_UpgradableTypeList = map[string]string{
	"":         "UpgradableTypeList is a collection of UpgradableType.",
	"metadata": "Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
	"items":    "Items is a list of AddOnPlacementScore",
}

func (UpgradableTypeList) SwaggerDoc() map[string]string {
	return map_UpgradableTypeList
}

// AUTO-GENERATED FUNCTIONS END HERE
