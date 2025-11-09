package handlers

import (
	"database"
	"encoding/json"
	"net/http"

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
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	condition := database.Condition{
		Name: req.Name,
	}

	result := h.db.Create(&condition)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(condition); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ConditionHandler) UpdateCondition(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid condition ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req types.UpdateConditionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	req.ID = id

	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	condition := database.Condition{
		Name: req.Name,
	}
	condition.ID = id

	result := h.db.Save(&condition)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(condition); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ConditionHandler) DeleteCondition(w http.ResponseWriter, r *http.Request) {

	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid condition ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	condition := database.Condition{}
	result := h.db.Delete(&condition, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ConditionHandler) GetCondition(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid condition ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	condition := database.Condition{}
	result := h.db.First(&condition, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(condition); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ConditionHandler) GetConditions(w http.ResponseWriter, r *http.Request) {
	conditions := []database.Condition{}
	result := h.db.Find(&conditions)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(conditions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
