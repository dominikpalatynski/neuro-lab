package handlers

import (
	"database"
	"encoding/json"
	"net/http"

	"config/types"
	"config/utils"
	apierrors "github.com/neuro-lab/errors"

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
		apierrors.WriteError(w, apierrors.NewBadRequestError("Invalid request body: "+err.Error(), r.URL.Path))
		return
	}

	if err := validate.Struct(req); err != nil {
		apierrors.WriteError(w, apierrors.NewValidationError(err, r.URL.Path))
		return
	}

	device := database.Device{
		Name: req.Name,
	}

	result := h.db.Create(&device)
	if result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(device)
}

func (h *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		apierrors.WriteError(w, apierrors.NewBadRequestError("Invalid device ID: "+err.Error(), r.URL.Path))
		return
	}

	var req types.UpdateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, apierrors.NewBadRequestError("Invalid request body: "+err.Error(), r.URL.Path))
		return
	}
	req.ID = id

	if err := validate.Struct(req); err != nil {
		apierrors.WriteError(w, apierrors.NewValidationError(err, r.URL.Path))
		return
	}

	device := database.Device{
		Name: req.Name,
	}
	device.ID = id

	result := h.db.Save(&device)
	if result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(device)
}

func (h *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		apierrors.WriteError(w, apierrors.NewBadRequestError("Invalid device ID: "+err.Error(), r.URL.Path))
		return
	}

	// First check if the device exists
	device := database.Device{}
	result := h.db.First(&device, id)
	if result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	// Delete the device
	result = h.db.Delete(&device, id)
	if result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		apierrors.WriteError(w, apierrors.NewBadRequestError("Invalid device ID: "+err.Error(), r.URL.Path))
		return
	}

	device := database.Device{}
	result := h.db.First(&device, id)
	if result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(device)
}

func (h *DeviceHandler) GetDevices(w http.ResponseWriter, r *http.Request) {
	devices := []database.Device{}
	result := h.db.Find(&devices)
	if result.Error != nil {
		apierrors.WriteError(w, apierrors.NewDatabaseError(result.Error, r.URL.Path))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(devices)
}
