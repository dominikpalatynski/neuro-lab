package handlers

import (
	"database"
	"encoding/json"
	"net/http"

	"config/types"

	"config/utils"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

type DeviceHandler struct {
	db *gorm.DB
}

func NewDeviceHandler(db *gorm.DB) *DeviceHandler {
	return &DeviceHandler{db: db}
}

func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {

	var req types.CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, "Validation failed: "+validationErrors.Error(), http.StatusBadRequest)
		return
	}

	device := database.Device{
		Name: req.Name,
	}

	result := h.db.Create(&device)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(device); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid device ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req types.UpdateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	req.ID = id

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, "Validation failed: "+validationErrors.Error(), http.StatusBadRequest)
		return
	}

	device := database.Device{
		Name: req.Name,
	}
	device.ID = id

	result := h.db.Save(&device)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(device); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {

	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid device ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	device := database.Device{}
	result := h.db.Delete(&device, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		http.Error(w, "Invalid device ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	device := database.Device{}
	result := h.db.First(&device, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(device); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandler) GetDevices(w http.ResponseWriter, r *http.Request) {
	devices := []database.Device{}
	result := h.db.Find(&devices)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(devices); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
