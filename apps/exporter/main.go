package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"database"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	parquet "github.com/parquet-go/parquet-go"
	kafka "github.com/segmentio/kafka-go"
)

var (
	bucketName = "neuro-lab"
	location   = "us-east-1"
)

type NotificationMessage struct {
	ScenarioID uint `json:"scenario_id"`
}

type ExportKey struct {
	FrameID   uint      `json:"frame_id"`
	Timestamp time.Time `json:"timestamp"`
}

// ParquetRow represents a single row in the wide-format parquet file
// Each row contains all metrics for a given timestamp
type ParquetRow struct {
	Timestamp time.Time `parquet:"timestamp,timestamp(microsecond)"`
	FrameID   uint      `parquet:"frame_id"`
	AccX      *float64  `parquet:"acc_x,optional"`
	AccY      *float64  `parquet:"acc_y,optional"`
	AccZ      *float64  `parquet:"acc_z,optional"`
	GyroX     *float64  `parquet:"gyro_x,optional"`
	GyroY     *float64  `parquet:"gyro_y,optional"`
	GyroZ     *float64  `parquet:"gyro_z,optional"`
	CurrV     *float64  `parquet:"curr_v,optional"`
	Temp      *float64  `parquet:"temp,optional"`
}

func createMinioClient() (*minio.Client, error) {
	endpoint := "localhost:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println("could not create minio client:", err)
		return nil, err
	}
	fmt.Println("Minio client created successfully!")
	return minioClient, nil
}

func CreateBucket(minioClient *minio.Client, bucketName string) error {
	err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			fmt.Println("Bucket already exists")
			return nil
		} else {
			fmt.Println("error while creating a bucket: ", err)
			return err
		}
	}
	fmt.Println("Bucket created successfully!")
	return nil
}

func PutObject(minioClient *minio.Client, bucketName string, objectName string, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	uploadInfo, err := minioClient.PutObject(context.Background(), bucketName, objectName, file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}
	fmt.Println("Successfully uploaded bytes: ", uploadInfo)
	return nil
}

// pivotSamplesToRows converts long-format samples to wide-format parquet rows
// Groups samples by timestamp - typically 8 metrics share the same timestamp
func pivotSamplesToRows(samples []database.ProcessedSample) []ParquetRow {
	// Group by timestamp
	timestampMap := make(map[ExportKey]map[string]float64)

	for _, sample := range samples {
		if _, exists := timestampMap[ExportKey{FrameID: sample.FrameID, Timestamp: sample.Timestamp}]; !exists {
			timestampMap[ExportKey{FrameID: sample.FrameID, Timestamp: sample.Timestamp}] = make(map[string]float64)
		}
		timestampMap[ExportKey{FrameID: sample.FrameID, Timestamp: sample.Timestamp}][sample.MetricName] = sample.Value
	}

	// Convert to ParquetRow slice
	rows := make([]ParquetRow, 0, len(timestampMap))

	for exportKey, metrics := range timestampMap {
		row := ParquetRow{Timestamp: exportKey.Timestamp, FrameID: exportKey.FrameID}

		// Map each metric to its corresponding field
		if val, ok := metrics["acc_x"]; ok {
			row.AccX = &val
		}
		if val, ok := metrics["acc_y"]; ok {
			row.AccY = &val
		}
		if val, ok := metrics["acc_z"]; ok {
			row.AccZ = &val
		}
		if val, ok := metrics["gyro_x"]; ok {
			row.GyroX = &val
		}
		if val, ok := metrics["gyro_y"]; ok {
			row.GyroY = &val
		}
		if val, ok := metrics["gyro_z"]; ok {
			row.GyroZ = &val
		}
		if val, ok := metrics["curr_v"]; ok {
			row.CurrV = &val
		}
		if val, ok := metrics["temp"]; ok {
			row.Temp = &val
		}

		rows = append(rows, row)
	}

	// Sort rows by FrameID and then by Timestamp for consistent order from lowest to biggest
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].FrameID != rows[j].FrameID {
			return rows[i].FrameID < rows[j].FrameID
		}
		return rows[i].Timestamp.Before(rows[j].Timestamp)
	})

	return rows
}

// exportData exports scenario samples to a parquet file with partitioned naming
// File path format: device_id=X/scenario_id=Y/data.parquet
func exportData(minioClient *minio.Client, samples []database.ProcessedSample, deviceID, scenarioID uint) error {
	if len(samples) == 0 {
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
	outputPath := filepath.Join(partitionPath, "data.parquet")

	// Pivot samples from long format to wide format
	rows := pivotSamplesToRows(samples)

	// Write to parquet file
	err := parquet.WriteFile(outputPath, rows)
	if err != nil {
		return fmt.Errorf("failed to write parquet file: %w", err)
	}

	err = PutObject(minioClient, bucketName, fmt.Sprintf("device_id=%d/scenario_id=%d/data.parquet", deviceID, scenarioID), outputPath)
	if err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}

	os.Remove(outputPath)
	fmt.Printf("Successfully exported %d rows to %s\n", len(rows), outputPath)
	return nil
}

func main() {
	db := database.Connect()
	minioClient, err := createMinioClient()
	if err != nil {
		fmt.Println("could not create minio client:", err)
		return
	}
	if err := CreateBucket(minioClient, bucketName); err != nil {
		return
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:19092"},
		GroupID:  "exporter-group",
		Topic:    "export.notification",
		MaxBytes: 100,
	})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("could not read message:", err)
			break
		}

		var notification NotificationMessage
		err = json.Unmarshal(msg.Value, &notification)
		if err != nil {
			fmt.Println("could not unmarshal message:", err)
			continue
		}

		// Get scenario with test session to retrieve device_id
		var scenario database.Scenario
		if err := db.Preload("TestSession").First(&scenario, notification.ScenarioID).Error; err != nil {
			fmt.Println("could not get scenario:", err)
			continue
		}

		if scenario.TestSession == nil {
			fmt.Println("scenario has no test session")
			continue
		}

		deviceID := scenario.TestSession.DeviceID

		// Get all processed samples for this scenario, ordered by timestamp and ID for stable sorting
		processedSamples := []database.ProcessedSample{}
		if err := db.Where("scenario_id = ?", notification.ScenarioID).Order("timestamp ASC, id ASC").Find(&processedSamples).Error; err != nil {
			fmt.Println("could not get processed samples:", err)
			continue
		}

		// Export to parquet file
		fmt.Printf("Exporting %d samples for scenario %d (device %d)\n", len(processedSamples), scenario.ID, deviceID)
		if err := exportData(minioClient, processedSamples, deviceID, scenario.ID); err != nil {
			fmt.Println("could not export data:", err)
			continue
		}

		fmt.Printf("Successfully exported scenario %d\n", scenario.ID)
	}
}
