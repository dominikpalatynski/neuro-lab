package main

import (
	"time"

	"gorm.io/gorm"
)

type ProcessedSample struct {
	gorm.Model
	DeviceID   string    `json:"device_id"`
	ScenarioID string    `json:"scenario_id"`
	MetricName string    `json:"metric_name"`
	Value      float64   `json:"value"`
	Timestamp  time.Time `json:"timestamp"`
}

type RawData struct {
	Data SensorData `json:"data"`
}

type SensorData struct {
	AccX  []float64 `json:"acc_x"`
	AccY  []float64 `json:"acc_y"`
	AccZ  []float64 `json:"acc_z"`
	GyroX []float64 `json:"gyro_x"`
	GyroY []float64 `json:"gyro_y"`
	GyroZ []float64 `json:"gyro_z"`
	CurrV []float64 `json:"curr_v"`
	Temp  []float64 `json:"temp"`
}
