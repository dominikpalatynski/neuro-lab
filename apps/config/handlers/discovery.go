package handlers

import (
	"encoding/json"
	"net/http"
	"types"

	"gorm.io/gorm"
)

type DiscoveryHandler struct {
	db *gorm.DB
}

func NewDiscoveryHandler(db *gorm.DB) *DiscoveryHandler {
	return &DiscoveryHandler{db: db}
}

func (h *DiscoveryHandler) GetAPIResources(w http.ResponseWriter, r *http.Request) {
	resourceList := types.APIResourceList{
		Kind:         "APIResourceList",
		APIVersion:   "v1",
		GroupVersion: "v1",
		Resources: []types.APIResource{
			{
				Name:         "devices",
				SingularName: "device",
				Kind:         "Device",
				Verbs:        []string{"create", "get", "list", "update", "delete"},
				ShortNames:   []string{},
			},
			{
				Name:         "test-sessions",
				SingularName: "test-session",
				Kind:         "TestSession",
				Verbs:        []string{"create", "get", "list", "update", "delete"},
				ShortNames:   []string{},
			},
			{
				Name:         "conditions",
				SingularName: "condition",
				Kind:         "Condition",
				Verbs:        []string{"create", "get", "list", "update", "delete"},
				ShortNames:   []string{},
			},
			{
				Name:         "condition-values",
				SingularName: "condition-value",
				Kind:         "ConditionValue",
				Verbs:        []string{"create", "get", "list", "update", "delete"},
				ShortNames:   []string{},
			},
			{
				Name:         "scenarios",
				SingularName: "scenario",
				Kind:         "Scenario",
				Verbs:        []string{"create", "get", "list", "update", "delete"},
				ShortNames:   []string{},
			},
			{
				Name:         "scenario-conditions",
				SingularName: "scenario-condition",
				Kind:         "ScenarioCondition",
				Verbs:        []string{"create", "get", "update", "delete"},
				ShortNames:   []string{},
			},
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resourceList)
}
