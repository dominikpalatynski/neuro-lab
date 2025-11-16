/*
References:

	git: https://github.com/segmentio/kafka-go
	doc: https://pkg.go.dev/github.com/segmentio/kafka-go#section-readme
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"database"

	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"gorm.io/gorm"
)

const (
	// Sample interval calculation:
	// 50 samples per batch, batches arrive every 80ms
	// samples_per_second = 50 / 0.08 = 625 Hz
	// interval = 1000ms / 625 = 1.6ms per sample (rounded to 2ms)
	sampleIntervalMs = 2
)

var (
	meter                         metric.Meter
	dbSaveDuration                metric.Int64Histogram
	transformerProcessingDuration metric.Int64Histogram
)

func getTimestamp(index float64, startTime time.Time, frequency float64) time.Time {
	return startTime.Add(time.Duration(index*frequency) * time.Millisecond)
}

func processAll(dbConnection *gorm.DB, metrics *[]ProcessedSample) {
	err := dbConnection.Create(&metrics).Error
	if err != nil {
		fmt.Println("could not create samples:", err)
		return
	}
	fmt.Println("batch created samples:", len(*metrics))
}

func main() {
	ctx := context.Background()
	shutdown, err := setupOTelSDK(ctx)
	if err != nil {
		log.Fatalf("failed to setup opentelemetry: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown opentelemetry: %v", err)
		}
	}()
	topic := "gateway.raw"
	db := database.Connect()
	db.AutoMigrate(&ProcessedSample{})
	readWithReader(db, topic, "transformer-group")
	select {}
}

// Read from the topic using kafka.Reader
// Readers can use consumer groups (but are not required to)
func readWithReader(db *gorm.DB, topic string, groupID string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:19092"},
		GroupID:  groupID,
		Topic:    topic,
		MaxBytes: 100, //per message
		// more options are available
	})
	var err error
	meter = otel.Meter("neuro-lab.transformer")

	transformerProcessingDuration, err = meter.Int64Histogram(
		"transformer.processing.duration",
		metric.WithDescription("Duration of transformer processing."),
		metric.WithUnit("ms"),
	)
	if err != nil {
		log.Fatalf("failed to create transformer processing duration histogram: %v", err)
	}

	dbSaveDuration, err = meter.Int64Histogram(
		"transformer.db.save.duration",
		metric.WithDescription("Duration of database saving."),
		metric.WithUnit("ms"),
	)
	if err != nil {
		log.Fatalf("failed to create database save duration histogram: %v", err)
	}

	fmt.Println("Consumer is running, waiting for messages...")
	for {
		transformerProcessingDurationStart := time.Now()
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("could not read message:", err)
			break
		}

		var rawData RawData
		err = json.Unmarshal(msg.Value, &rawData)
		if err != nil {
			fmt.Println("could not unmarshal message:", err)
			continue
		}

		deviceID := rawData.DeviceID
		scenarioID := rawData.ScenarioID
		frameID := rawData.FrameID

		timestamp, err := time.Parse("2006-01-02 15:04:05.000000", rawData.Timestamp)
		if err != nil {
			fmt.Println("could not parse timestamp:", err)
			continue
		}
		// Process all channels
		metrics := []ProcessedSample{}

		for index, value := range rawData.Data.AccX {
			metrics = append(metrics, ProcessedSample{
				DeviceID:   deviceID,
				ScenarioID: scenarioID,
				FrameID:    frameID,
				MetricName: "acc_x",
				Value:      float64(value),
				Timestamp:  getTimestamp(float64(index), timestamp, sampleIntervalMs),
			})
		}
		for index, value := range rawData.Data.AccY {
			metrics = append(metrics, ProcessedSample{
				DeviceID:   deviceID,
				ScenarioID: scenarioID,
				FrameID:    frameID,
				MetricName: "acc_y",
				Value:      float64(value),
				Timestamp:  getTimestamp(float64(index), timestamp, sampleIntervalMs),
			})
		}
		for index, value := range rawData.Data.AccZ {
			metrics = append(metrics, ProcessedSample{
				DeviceID:   deviceID,
				ScenarioID: scenarioID,
				FrameID:    frameID,
				MetricName: "acc_z",
				Value:      float64(value),
				Timestamp:  getTimestamp(float64(index), timestamp, sampleIntervalMs),
			})
		}
		for index, value := range rawData.Data.GyroX {
			metrics = append(metrics, ProcessedSample{
				DeviceID:   deviceID,
				ScenarioID: scenarioID,
				FrameID:    frameID,
				MetricName: "gyro_x",
				Value:      float64(value),
				Timestamp:  getTimestamp(float64(index), timestamp, sampleIntervalMs),
			})
		}
		for index, value := range rawData.Data.GyroY {
			metrics = append(metrics, ProcessedSample{
				DeviceID:   deviceID,
				ScenarioID: scenarioID,
				FrameID:    frameID,
				MetricName: "gyro_y",
				Value:      float64(value),
				Timestamp:  getTimestamp(float64(index), timestamp, sampleIntervalMs),
			})
		}
		for index, value := range rawData.Data.GyroZ {
			metrics = append(metrics, ProcessedSample{
				DeviceID:   deviceID,
				ScenarioID: scenarioID,
				FrameID:    frameID,
				MetricName: "gyro_z",
				Value:      float64(value),
				Timestamp:  getTimestamp(float64(index), timestamp, sampleIntervalMs),
			})
		}
		for index, value := range rawData.Data.CurrV {
			metrics = append(metrics, ProcessedSample{
				DeviceID:   deviceID,
				ScenarioID: scenarioID,
				FrameID:    frameID,
				MetricName: "curr_v",
				Value:      float64(value),
				Timestamp:  getTimestamp(float64(index), timestamp, sampleIntervalMs),
			})
		}
		for index, value := range rawData.Data.Temp {
			metrics = append(metrics, ProcessedSample{
				DeviceID:   deviceID,
				ScenarioID: scenarioID,
				FrameID:    frameID,
				MetricName: "temp",
				Value:      float64(value),
				Timestamp:  getTimestamp(float64(index), timestamp, sampleIntervalMs),
			})
		}

		dbSaveDurationStart := time.Now()
		processAll(db, &metrics)
		dbSaveDuration.Record(context.Background(), int64(time.Since(dbSaveDurationStart).Milliseconds()))
		transformerProcessingDuration.Record(context.Background(), int64(time.Since(transformerProcessingDurationStart).Milliseconds()))
	}

	if err := r.Close(); err != nil {
		fmt.Println("failed to close reader:", err)
	}
}
