package handlers

import (
	"config/utils"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"net/http"

	"database"
	"encoding/csv"
	"os"
	"path/filepath"

	apierrors "github.com/neuro-lab/errors"
	kafka "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

var channelNames = []string{"acc_x", "acc_y", "acc_z", "gyro_x", "gyro_y", "gyro_z", "curr_v", "temp"}

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

	data := []database.ProcessedChannel{}
	err = h.db.Where("scenario_id = ?", scenarioID).Find(&data).Error
	if err != nil {
		apierrors.WriteError(w, apierrors.NewInternalError(r.URL.Path))
		return
	}

	err = h.exportDataToCSV(data, 1, scenarioID)
	if err != nil {
		apierrors.WriteError(w, apierrors.NewInternalError(r.URL.Path))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ExportHandler) exportDataToCSV(data []database.ProcessedChannel, deviceID, scenarioID uint) error {
	if len(data) == 0 {
		return fmt.Errorf("no samples to export")
	}

	outputDir := "./exports"
	partitionPath := filepath.Join(
		outputDir,
		fmt.Sprintf("device_id=%d", deviceID),
		fmt.Sprintf("scenario_id=%d", scenarioID),
	)

	if err := os.MkdirAll(partitionPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Group data by FrameID
	frameData := make(map[uint]map[string][]float64)
	for _, channel := range data {
		if frameData[channel.FrameID] == nil {
			frameData[channel.FrameID] = make(map[string][]float64)
		}
		frameData[channel.FrameID][channel.MetricName] = channel.Values
	}

	// Create output file path
	outputPath := filepath.Join(partitionPath, "data.csv")
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := append([]string{"frame_id"}, channelNames...)
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Sort frame IDs to ensure consistent ordering
	frameIDs := make([]uint, 0, len(frameData))
	for frameID := range frameData {
		frameIDs = append(frameIDs, frameID)
	}
	sort.Slice(frameIDs, func(i, j int) bool {
		return frameIDs[i] < frameIDs[j]
	})

	// Write data rows in sorted order
	for _, frameID := range frameIDs {
		channels := frameData[frameID]

		// Determine the number of values (assumes all channels have same length)
		var numValues int
		for _, values := range channels {
			numValues = len(values)
			break
		}

		// Write one row per value index
		for i := 0; i < numValues; i++ {
			row := []string{fmt.Sprintf("%d", frameID)}
			for _, channelName := range channelNames {
				if values, ok := channels[channelName]; ok && i < len(values) {
					row = append(row, fmt.Sprintf("%f", values[i]))
				} else {
					row = append(row, "")
				}
			}
			if err := writer.Write(row); err != nil {
				return fmt.Errorf("failed to write row: %w", err)
			}
		}
	}

	return nil
}
