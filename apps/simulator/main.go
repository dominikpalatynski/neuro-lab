package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.yaml.in/yaml/v2"
)

var channelNames = []string{"acc_x", "acc_y", "acc_z", "gyro_x", "gyro_y", "gyro_z", "curr_v", "temp"}

type Config struct {
	Interval int            `yaml:"interval"`
	Devices  []DeviceConfig `yaml:"devices"`
}

type DeviceConfig struct {
	ScenarioID int `yaml:"scenario_id"`
	DeviceID   int `yaml:"device_id"`
}

type SensorData struct {
	ScenarioID int           `json:"scenario_id"`
	DeviceID   int           `json:"device_id"`
	Data       []ChannelData `json:"data"`
}

type Data struct {
	AccX  []float64 `json:"acc_x"`
	AccY  []float64 `json:"acc_y"`
	AccZ  []float64 `json:"acc_z"`
	GyroX []float64 `json:"gyro_x"`
	GyroY []float64 `json:"gyro_y"`
	GyroZ []float64 `json:"gyro_z"`
	CurrV []float64 `json:"curr_v"`
	Temp  []float64 `json:"temp"`
}

type ChannelData struct {
	Values      []float64 `json:"values"`
	ChannelName string    `json:"channel_name"`
}

func generateSensorData(data []float64) []ChannelData {
	channelData := []ChannelData{}
	for _, channelName := range channelNames {
		channelData = append(channelData, ChannelData{
			Values:      data,
			ChannelName: channelName,
		})
	}
	return channelData
}

func initMQTTClient() (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().AddBroker("localhost:1884")
	opts.SetClientID("simulator_mqtt_client")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

func main() {

	client, err := initMQTTClient()
	if err != nil {
		log.Fatalf("Failed to initialize MQTT client: %v", err)
	}

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	cfg := Config{}
	yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %v", err)
	}

	for _, device := range cfg.Devices {
		go func(device DeviceConfig) {
			start := 0
			for {
				data := []float64{}
				for i := 0; i < 50; i++ {
					data = append(data, float64(start+i))
				}
				fmt.Printf("Sending data for Device ID: %d, Scenario ID: %d\n", device.DeviceID, device.ScenarioID)
				channelData := generateSensorData(data)
				jsonData, err := json.Marshal(SensorData{
					ScenarioID: device.ScenarioID,
					DeviceID:   device.DeviceID,
					Data:       channelData,
				})
				if err != nil {
					log.Fatalf("Failed to marshal data: %v", err)
				}

				token := client.Publish("device/"+strconv.Itoa(device.DeviceID)+"/raw", 0, false, jsonData)
				token.Wait()
				if token.Error() != nil {
					log.Fatalf("Failed to publish data: %v", token.Error().Error())
				}
				start += 50
				fmt.Printf("Sent %d samples for Device ID: %d, Scenario ID: %d\n", start, device.DeviceID, device.ScenarioID)
				time.Sleep(time.Duration(cfg.Interval) * time.Millisecond)
			}
		}(device)
	}
	select {}
}
