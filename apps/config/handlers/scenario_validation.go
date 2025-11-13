package handlers

import (
	"database"
	"encoding/json"
	"net/http"
	"types"

	apierrors "github.com/neuro-lab/errors"

	"gorm.io/gorm"
)

type ScenarioValidationHandler struct {
	db *gorm.DB
}

func NewScenarioValidationHandler(db *gorm.DB) *ScenarioValidationHandler {
	return &ScenarioValidationHandler{db: db}
}

func (h *ScenarioValidationHandler) ValidateScenario(w http.ResponseWriter, r *http.Request) {
	var req types.ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, apierrors.NewBadRequestError("Invalid request body: "+err.Error(), r.URL.Path))
		return
	}

	// Validate request using validator
	if err := validate.Struct(req); err != nil {
		apierrors.WriteError(w, apierrors.NewValidationError(err, r.URL.Path))
		return
	}

	// Check if scenario exists in database
	scenario := database.Scenario{}
	result := h.db.First(&scenario, req.ScenarioID)
	if result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	// Return success response with scenario data
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(scenario)
}
