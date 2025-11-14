package handlers

import (
	"database"
	"encoding/json"
	"net/http"
	"strconv"
	"types"

	"config/utils"

	"github.com/go-chi/chi/v5"
	apierrors "github.com/neuro-lab/errors"
	"gorm.io/gorm"
)

type ScenarioHandler struct {
	db *gorm.DB
}

func NewScenarioHandler(db *gorm.DB) *ScenarioHandler {
	return &ScenarioHandler{db: db}
}

func (h *ScenarioHandler) CreateScenario(w http.ResponseWriter, r *http.Request) {

	var req types.CreateScenarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that the test session exists
	testSession := database.TestSession{}
	if result := h.db.First(&testSession, req.TestSessionID); result.Error != nil {
		http.Error(w, "TestSession does not exist", http.StatusBadRequest)
		return
	}

	scenario := database.Scenario{
		Name:          req.Name,
		TestSessionID: req.TestSessionID,
	}

	result := h.db.Create(&scenario)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Reload with TestSession relationship
	if err := h.db.Preload("TestSession").First(&scenario, scenario.ID).Error; err != nil {
		http.Error(w, "Error loading test session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(scenario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create Scenario with Condition Values
func (h *ScenarioHandler) CreateScenarioWithConditionValues(w http.ResponseWriter, r *http.Request) {
	var req types.CreateScenarioWithConditionValuesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that the test session exists
	testSession := database.TestSession{}
	if result := h.db.First(&testSession, req.TestSessionID); result.Error != nil {
		http.Error(w, "TestSession does not exist", http.StatusBadRequest)
		return
	}

	// Validate that the condition values exist
	conditionValues := []database.ConditionValue{}
	if result := h.db.First(&conditionValues, req.ConditionValueIDs); result.Error != nil {
		http.Error(w, "ConditionValues do not exist", http.StatusBadRequest)
		return
	}

	scenario := database.Scenario{
		Name:          req.Name,
		TestSessionID: req.TestSessionID,
	}

	transactionResult := h.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&scenario)
		if result.Error != nil {
			return result.Error
		}

		scenarioConditions := []database.ScenarioCondition{}
		for _, conditionValueID := range req.ConditionValueIDs {
			scenarioConditions = append(scenarioConditions, database.ScenarioCondition{
				ScenarioID:       scenario.ID,
				ConditionValueID: conditionValueID,
			})
		}

		result = tx.Create(&scenarioConditions)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if transactionResult != nil {
		http.Error(w, transactionResult.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(scenario); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ScenarioHandler) UpdateScenario(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid scenario ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req types.UpdateScenarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	req.ID = id

	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that the test session exists
	testSession := database.TestSession{}
	if result := h.db.First(&testSession, req.TestSessionID); result.Error != nil {
		http.Error(w, "TestSession does not exist", http.StatusBadRequest)
		return
	}

	scenario := database.Scenario{
		Name:          req.Name,
		TestSessionID: req.TestSessionID,
	}
	scenario.ID = id

	result := h.db.Save(&scenario)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Reload with TestSession relationship
	if err := h.db.Preload("TestSession").First(&scenario, scenario.ID).Error; err != nil {
		http.Error(w, "Error loading test session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(scenario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ScenarioHandler) DeleteScenario(w http.ResponseWriter, r *http.Request) {

	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid scenario ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	scenario := database.Scenario{}
	result := h.db.Delete(&scenario, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ScenarioHandler) GetScenario(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid scenario ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	scenario := database.Scenario{}
	result := h.db.Preload("TestSession").Preload("ScenarioConditions").Preload("ScenarioConditions.ConditionValue").First(&scenario, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(scenario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ScenarioHandler) GetScenarios(w http.ResponseWriter, r *http.Request) {
	scenarios := []database.Scenario{}
	result := h.db.Preload("TestSession").Preload("ScenarioConditions").Preload("ScenarioConditions.ConditionValue").Find(&scenarios)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(scenarios); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ScenarioHandler) GetScenariosByTestSession(w http.ResponseWriter, r *http.Request) {
	testSessionIDStr := chi.URLParam(r, "testSessionID")
	testSessionID64, err := strconv.ParseUint(testSessionIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid test session ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	testSessionID := uint(testSessionID64)

	scenarios := []database.Scenario{}
	result := h.db.Preload("TestSession").Preload("ScenarioConditions").Preload("ScenarioConditions.ConditionValue").Where("test_session_id = ?", testSessionID).Find(&scenarios)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(scenarios); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ScenarioHandler) ActivateScenario(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		apierrors.WriteError(w, apierrors.NewBadRequestError("Invalid scenario ID: "+err.Error(), r.URL.Path))
		return
	}

	scenario := database.Scenario{}
	if result := h.db.First(&scenario, id); result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	if scenario.Status == database.StatusActive {
		apierrors.WriteError(w, apierrors.NewBadRequestError("Scenario is already active", r.URL.Path))
		return
	}

	scenario.Status = database.StatusActive
	result := h.db.Save(&scenario)
	if result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ScenarioHandler) DeactivateScenario(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		apierrors.WriteError(w, apierrors.NewBadRequestError("Invalid scenario ID: "+err.Error(), r.URL.Path))
		return
	}

	scenario := database.Scenario{}
	if result := h.db.First(&scenario, id); result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	if scenario.Status == database.StatusInactive {
		apierrors.WriteError(w, apierrors.NewBadRequestError("Scenario is already deactivated", r.URL.Path))
		return
	}

	scenario.Status = database.StatusInactive
	result := h.db.Save(&scenario)
	if result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ScenarioHandler) CompleteScenario(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		apierrors.WriteError(w, apierrors.NewBadRequestError("Invalid scenario ID: "+err.Error(), r.URL.Path))
		return
	}

	scenario := database.Scenario{}
	if result := h.db.First(&scenario, id); result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	if scenario.Status == database.StatusCompleted {
		apierrors.WriteError(w, apierrors.NewBadRequestError("Scenario is already completed", r.URL.Path))
		return
	}

	scenario.Status = database.StatusCompleted
	result := h.db.Save(&scenario)
	if result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.WriteHeader(http.StatusOK)
}
