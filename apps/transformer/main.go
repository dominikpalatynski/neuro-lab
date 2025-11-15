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
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	kafka "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

var (
	// Histogram to measure operation duration
	operationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kafka_event_processing_duration_seconds",
			Help:    "Duration of Kafka event processing in seconds",
			Buckets: prometheus.DefBuckets, // Default buckets: 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
		},
		[]string{"topic", "status"}, // Labels for filtering
	)

	// Counter for total processed events
	eventsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_events_processed_total",
			Help: "Total number of Kafka events processed",
		},
		[]string{"topic", "status"},
	)
)

const (
	// Sample interval calculation:
	// 50 samples per batch, batches arrive every 80ms
	// samples_per_second = 50 / 0.08 = 625 Hz
	// interval = 1000ms / 625 = 1.6ms per sample (rounded to 2ms)
	sampleIntervalMs = 2
)

func getTimestamp(index float64, startTime time.Time, frequency float64) time.Time {
	return startTime.Add(time.Duration(index*frequency) * time.Millisecond)
}

func StartMetricsServer(port string) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("Starting metrics server on :%s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()
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
	StartMetricsServer("9091")
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

	fmt.Println("Consumer is running, waiting for messages...")
	for {
		msg, err := r.ReadMessage(context.Background())
		start := time.Now()
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

		processAll(db, &metrics)
		duration := float64(time.Since(start).Milliseconds())
		operationDuration.WithLabelValues(topic, "success").Observe(duration)
		eventsProcessed.WithLabelValues(topic, "success").Inc()
	}

	if err := r.Close(); err != nil {
		fmt.Println("failed to close reader:", err)
	}
}
