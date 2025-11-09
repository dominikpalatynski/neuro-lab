package handlers

import (
	"database"
	"encoding/json"
	"net/http"

	"config/types"
	"config/utils"

	"gorm.io/gorm"
)

type ScenarioConditionHandler struct {
	db *gorm.DB
}

func NewScenarioConditionHandler(db *gorm.DB) *ScenarioConditionHandler {
	return &ScenarioConditionHandler{db: db}
}

func (h *ScenarioConditionHandler) CreateScenarioCondition(w http.ResponseWriter, r *http.Request) {

	var req types.CreateScenarioConditionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that the scenario exists
	scenario := database.Scenario{}
	if result := h.db.First(&scenario, req.ScenarioID); result.Error != nil {
		http.Error(w, "Scenario does not exist", http.StatusBadRequest)
		return
	}

	// Validate that the condition value exists
	conditionValue := database.ConditionValue{}
	if result := h.db.First(&conditionValue, req.ConditionValueID); result.Error != nil {
		http.Error(w, "ConditionValue does not exist", http.StatusBadRequest)
		return
	}

	scenarioCondition := database.ScenarioCondition{
		ScenarioID:       req.ScenarioID,
		ConditionValueID: req.ConditionValueID,
	}

	result := h.db.Create(&scenarioCondition)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Reload with Scenario and ConditionValue relationships
	if err := h.db.Preload("Scenario").Preload("ConditionValue").First(&scenarioCondition, scenarioCondition.ID).Error; err != nil {
		http.Error(w, "Error loading relationships: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(scenarioCondition); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ScenarioConditionHandler) UpdateScenarioCondition(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid scenario condition ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req types.UpdateScenarioConditionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	req.ID = id

	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that the scenario exists
	scenario := database.Scenario{}
	if result := h.db.First(&scenario, req.ScenarioID); result.Error != nil {
		http.Error(w, "Scenario does not exist", http.StatusBadRequest)
		return
	}

	// Validate that the condition value exists
	conditionValue := database.ConditionValue{}
	if result := h.db.First(&conditionValue, req.ConditionValueID); result.Error != nil {
		http.Error(w, "ConditionValue does not exist", http.StatusBadRequest)
		return
	}

	scenarioCondition := database.ScenarioCondition{
		ScenarioID:       req.ScenarioID,
		ConditionValueID: req.ConditionValueID,
	}
	scenarioCondition.ID = id

	result := h.db.Save(&scenarioCondition)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Reload with Scenario and ConditionValue relationships
	if err := h.db.Preload("Scenario").Preload("ConditionValue").First(&scenarioCondition, scenarioCondition.ID).Error; err != nil {
		http.Error(w, "Error loading relationships: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(scenarioCondition); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ScenarioConditionHandler) DeleteScenarioCondition(w http.ResponseWriter, r *http.Request) {

	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid scenario condition ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	scenarioCondition := database.ScenarioCondition{}
	result := h.db.Delete(&scenarioCondition, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ScenarioConditionHandler) GetScenarioCondition(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid scenario condition ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	scenarioCondition := database.ScenarioCondition{}
	result := h.db.Preload("Scenario").Preload("ConditionValue").First(&scenarioCondition, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(scenarioCondition); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
