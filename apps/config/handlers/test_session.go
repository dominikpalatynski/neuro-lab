package handlers

import (
	"database"
	"encoding/json"
	"net/http"

	"config/types"
	"config/utils"

	"gorm.io/gorm"
)

type TestSessionHandler struct {
	db *gorm.DB
}

func NewTestSessionHandler(db *gorm.DB) *TestSessionHandler {
	return &TestSessionHandler{db: db}
}

func (h *TestSessionHandler) CreateTestSession(w http.ResponseWriter, r *http.Request) {

	var req types.CreateTestSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	device := database.Device{}
	if result := h.db.First(&device, req.DeviceID); result.Error != nil {
		http.Error(w, "Device with given ID does not exist", http.StatusBadRequest)
		return
	}

	testSession := database.TestSession{
		Name:     req.Name,
		DeviceID: req.DeviceID,
	}

	result := h.db.Create(&testSession)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(testSession); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TestSessionHandler) UpdateTestSession(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid test session ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req types.UpdateTestSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	req.ID = id

	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that the device exists
	device := database.Device{}
	if result := h.db.First(&device, req.DeviceID); result.Error != nil {
		http.Error(w, "Device with ID "+string(rune(req.DeviceID))+" does not exist", http.StatusBadRequest)
		return
	}

	testSession := database.TestSession{
		Name:     req.Name,
		DeviceID: req.DeviceID,
	}
	testSession.ID = id

	result := h.db.Save(&testSession)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(testSession); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TestSessionHandler) DeleteTestSession(w http.ResponseWriter, r *http.Request) {

	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid test session ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	testSession := database.TestSession{}
	result := h.db.Delete(&testSession, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TestSessionHandler) GetTestSession(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid test session ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	testSession := database.TestSession{}
	result := h.db.First(&testSession, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(testSession); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
