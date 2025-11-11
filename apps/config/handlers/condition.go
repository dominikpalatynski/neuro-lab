package handlers

import (
	"database"
	"encoding/json"
	"net/http"

	"config/errors"
	"config/types"
	"config/utils"

	"gorm.io/gorm"
)

type ConditionHandler struct {
	db *gorm.DB
}

func NewConditionHandler(db *gorm.DB) *ConditionHandler {
	return &ConditionHandler{db: db}
}

func (h *ConditionHandler) CreateCondition(w http.ResponseWriter, r *http.Request) {
	var req types.CreateConditionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.NewBadRequestError("Invalid request body: "+err.Error(), r.URL.Path))
		return
	}

	if err := validate.Struct(req); err != nil {
		errors.WriteError(w, errors.NewValidationError(err, r.URL.Path))
		return
	}

	condition := database.Condition{
		Name: req.Name,
	}

	result := h.db.Create(&condition)
	if result.Error != nil {
		errors.WriteError(w, errors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(condition)
}

func (h *ConditionHandler) UpdateCondition(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		errors.WriteError(w, errors.NewBadRequestError("Invalid condition ID: "+err.Error(), r.URL.Path))
		return
	}

	var req types.UpdateConditionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.NewBadRequestError("Invalid request body: "+err.Error(), r.URL.Path))
		return
	}
	req.ID = id

	if err := validate.Struct(req); err != nil {
		errors.WriteError(w, errors.NewValidationError(err, r.URL.Path))
		return
	}

	condition := database.Condition{
		Name: req.Name,
	}
	condition.ID = id

	result := h.db.Save(&condition)
	if result.Error != nil {
		errors.WriteError(w, errors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(condition)
}

func (h *ConditionHandler) DeleteCondition(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		errors.WriteError(w, errors.NewBadRequestError("Invalid condition ID: "+err.Error(), r.URL.Path))
		return
	}

	// First check if the condition exists
	condition := database.Condition{}
	result := h.db.First(&condition, id)
	if result.Error != nil {
		errors.WriteError(w, errors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	// Delete the condition
	result = h.db.Delete(&condition, id)
	if result.Error != nil {
		errors.WriteError(w, errors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ConditionHandler) GetCondition(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		errors.WriteError(w, errors.NewBadRequestError("Invalid condition ID: "+err.Error(), r.URL.Path))
		return
	}

	condition := database.Condition{}
	result := h.db.First(&condition, id)
	if result.Error != nil {
		errors.WriteError(w, errors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(condition)
}

func (h *ConditionHandler) GetConditions(w http.ResponseWriter, r *http.Request) {
	conditions := []database.Condition{}
	result := h.db.Find(&conditions)
	if result.Error != nil {
		errors.WriteError(w, errors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conditions)
}
