package handlers

import (
	"config/utils"
	"encoding/json"
	"fmt"
	"time"

	"net/http"

	apierrors "github.com/neuro-lab/errors"
	kafka "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type NotificationMessage struct {
	ScenarioID uint `json:"scenario_id"`
}

type ExportHandler struct {
	db    *gorm.DB
	kafka *kafka.Conn
}

func NewExportHandler(kafka *kafka.Conn, db *gorm.DB) *ExportHandler {
	return &ExportHandler{kafka: kafka, db: db}
}

func (h *ExportHandler) sendToKafka(notification NotificationMessage) error {
	h.kafka.SetWriteDeadline(time.Now().Add(1 * time.Second))

	jsonData, marshalErr := json.Marshal(notification)
	if marshalErr != nil {
		return fmt.Errorf("error marshalling notification: %v", marshalErr)
	}

	_, err := h.kafka.WriteMessages(
		kafka.Message{Value: jsonData})
	if err != nil {
		return fmt.Errorf("failed to write messages: %v", err)
	}

	return nil
}

func (h *ExportHandler) ExportData(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := utils.ParseID(r)
	if err != nil {
		apierrors.WriteError(w, apierrors.NewBadRequestError("Invalid scenario ID: "+err.Error(), r.URL.Path))
		return
	}

	notification := NotificationMessage{ScenarioID: scenarioID}
	kafkaErr := h.sendToKafka(notification)
	if kafkaErr != nil {
		apierrors.WriteError(w, apierrors.NewInternalError(r.URL.Path))
		return
	}

	w.WriteHeader(http.StatusOK)
}
