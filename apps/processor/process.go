package main

import (
	"encoding/json"
	"fmt"

	"communication"

	"time"

	"database"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	apierrors "github.com/neuro-lab/errors"
)

type SensorDataRaw struct {
	Data       []ChannelData `json:"data"`
	DeviceID   int           `json:"device_id"`
	ScenarioID int           `json:"scenario_id" validate:"required"`
	Timestamp  string        `json:"timestamp" `
	FrameID    int           `json:"frame_id" `
}
type ChannelData struct {
	Values      database.Float8Array `json:"values"`
	ChannelName string               `json:"channel_name"`
}

type ValidationRequest struct {
	ScenarioID int `json:"scenario_id"`
}

func processMessage(mqttClient mqtt.Client, msg mqtt.Message) {
	var sensorData *SensorDataRaw

	err := json.Unmarshal(msg.Payload(), &sensorData)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return
	}

	if err := validateScenarioRaw(sensorData); err != nil {
		fmt.Println("Error validating scenario:", err)
		return
	}

	processData(sensorData)
	fmt.Println("Processed message successfully")
}

func validateScenarioRaw(sensorData *SensorDataRaw) error {
	reqBytes, marshalErr := json.Marshal(ValidationRequest{ScenarioID: int(sensorData.ScenarioID)})
	if marshalErr != nil {
		return fmt.Errorf("error marshalling request: %v", marshalErr)
	}

	resp, httpErr := communication.SendRequest("POST", "http://localhost:3002/api/v1/scenario-validation", reqBytes)
	if httpErr != nil {
		return fmt.Errorf("error sending request: %v", httpErr)
	}
	if resp.StatusCode != 200 {
		var errorResponse apierrors.ErrorResponse
		if err := json.Unmarshal(resp.Body, &errorResponse); err != nil {
			return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, resp.Body)
		}

		if errorResponse.Type == apierrors.TypeValidationFailed {
			return fmt.Errorf("%s: %s\n%v", errorResponse.Title, errorResponse.Detail, errorResponse.Errors)
		}
		return fmt.Errorf("error sending request: %v", resp.Body)
	}

	return nil
}

func getFrameId(scenarioId int) int {
	frameIdValue := (*frameIds)[scenarioId]
	if frameIdValue == 0 {
		frameIdValue = 1
	} else {
		frameIdValue++
	}
	return frameIdValue
}

func setFrameId(scenarioId int, frameId int) {
	(*frameIds)[scenarioId] = frameId
}

func processData(sensorData *SensorDataRaw) error {
	metrics := []database.ProcessedChannel{}
	frameId := getFrameId(sensorData.ScenarioID)
	setFrameId(sensorData.ScenarioID, frameId)
	fmt.Println("Frame ID:", frameId)
	for _, channel := range sensorData.Data {
		metrics = append(metrics, database.ProcessedChannel{
			Values:     channel.Values,
			MetricName: channel.ChannelName,
			Timestamp:  time.Now(),
			FrameID:    uint(frameId),
			DeviceID:   uint(sensorData.DeviceID),
			ScenarioID: uint(sensorData.ScenarioID),
		})
	}
	err := db.Create(&metrics).Error
	if err != nil {
		return fmt.Errorf("could not create samples: %v", err)
	}
	return nil
}
