package main

import (
	"database"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func exportDataToCSV(data []database.ProcessedChannel, deviceID, scenarioID uint) error {
	if len(data) == 0 {
		return fmt.Errorf("no samples to export")
	}

	// Create partition directory structure
	outputDir := "./exports"
	partitionPath := filepath.Join(
		outputDir,
		fmt.Sprintf("device_id=%d", deviceID),
		fmt.Sprintf("scenario_id=%d", scenarioID),
	)

	if err := os.MkdirAll(partitionPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
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

	writer.Write([]string{"timestamp", "frame_id", "device_id", "scenario_id", "metric_name", "values"})
	for _, channel := range data {
		writer.Write([]string{channel.Timestamp.Format(time.RFC3339), fmt.Sprintf("%d", channel.FrameID), fmt.Sprintf("%d", channel.DeviceID), fmt.Sprintf("%d", channel.ScenarioID), channel.MetricName, fmt.Sprintf("%f", channel.Values)})
	}

	return nil
}
