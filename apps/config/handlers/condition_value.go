package handlers

import (
	"database"
	"encoding/json"
	"net/http"

	"config/types"
	"config/utils"

	"gorm.io/gorm"
)

type ConditionValueHandler struct {
	db *gorm.DB
}

func NewConditionValueHandler(db *gorm.DB) *ConditionValueHandler {
	return &ConditionValueHandler{db: db}
}

func (h *ConditionValueHandler) CreateConditionValue(w http.ResponseWriter, r *http.Request) {

	var req types.CreateConditionValueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that the condition exists
	condition := database.Condition{}
	if result := h.db.First(&condition, req.ConditionID); result.Error != nil {
		http.Error(w, "Condition does not exist", http.StatusBadRequest)
		return
	}

	conditionValue := database.ConditionValue{
		Value:       req.Value,
		ConditionID: req.ConditionID,
	}

	result := h.db.Create(&conditionValue)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Reload with Condition relationship
	if err := h.db.Preload("Condition").First(&conditionValue, conditionValue.ID).Error; err != nil {
		http.Error(w, "Error loading condition: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(conditionValue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ConditionValueHandler) UpdateConditionValue(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid condition value ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req types.UpdateConditionValueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	req.ID = id

	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that the condition exists
	condition := database.Condition{}
	if result := h.db.First(&condition, req.ConditionID); result.Error != nil {
		http.Error(w, "Condition does not exist", http.StatusBadRequest)
		return
	}

	conditionValue := database.ConditionValue{
		Value:       req.Value,
		ConditionID: req.ConditionID,
	}
	conditionValue.ID = id

	result := h.db.Save(&conditionValue)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Reload with Condition relationship
	if err := h.db.Preload("Condition").First(&conditionValue, conditionValue.ID).Error; err != nil {
		http.Error(w, "Error loading condition: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(conditionValue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ConditionValueHandler) DeleteConditionValue(w http.ResponseWriter, r *http.Request) {

	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid condition value ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	conditionValue := database.ConditionValue{}
	result := h.db.Delete(&conditionValue, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ConditionValueHandler) GetConditionValue(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid condition value ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	conditionValue := database.ConditionValue{}
	result := h.db.Preload("Condition").First(&conditionValue, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(conditionValue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
