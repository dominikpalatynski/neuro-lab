package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"communication"
	"context"
	"encoding/json"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	apierrors "github.com/neuro-lab/errors"
	kafka "github.com/segmentio/kafka-go"
)

// Global atomic counter for frame IDs
var frameCounter atomic.Uint64
var (
	meter              metric.Meter
	messageCounter     metric.Int64Counter
	gatewayDuration    metric.Int64Histogram
	validationDuration metric.Int64Histogram
)

type SensorData struct {
	Data       Data   `json:"data"`
	DeviceID   int    `json:"device_id"`
	ScenarioID int    `json:"scenario_id"`
	Timestamp  string `json:"timestamp"`
	FrameID    int    `json:"frame_id"`
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

type ValidationRequest struct {
	ScenarioID int `json:"scenario_id"`
}

// Connect to the specified topic and partition in the server
func connect(topic string, partition int) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp",
		"localhost:19092", topic, partition)
	if err != nil {
		fmt.Println("failed to dial leader", err)
	}
	return conn, err
} //end connect

func sendViaKafka(conn *kafka.Conn, sensorData *SensorData) {
	// Set write deadline to ensure timely sending
	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))

	jsonData, marshalErr := json.Marshal(sensorData)
	if marshalErr != nil {
		fmt.Println("Error marshalling sensor data:", marshalErr)
		return
	}

	_, err := conn.WriteMessages(
		kafka.Message{Value: jsonData})
	if err != nil {
		fmt.Println("failed to write messages:", err)
		return
	}

	// Flush to send messages immediately (real-time delivery)
	if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		fmt.Println("failed to set deadline for flush:", err)
	}
}

func processMessage(msg []byte) (*SensorData, error) {
	var sensorData SensorData

	err := json.Unmarshal(msg, &sensorData)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return nil, fmt.Errorf("error unmarshalling message: %v", err)
	}

	reqBytes, marshalErr := json.Marshal(ValidationRequest{ScenarioID: int(sensorData.ScenarioID)})
	if marshalErr != nil {
		fmt.Println("Error marshalling request:", marshalErr)
		return nil, fmt.Errorf("error marshalling request: %v", marshalErr)
	}

	start := time.Now()
	resp, httpErr := communication.SendRequest("POST", "http://localhost:3002/api/v1/scenario-validation", reqBytes)
	if httpErr != nil {
		fmt.Println("Error sending request:", httpErr)
		return nil, fmt.Errorf("error sending request: %v", httpErr)
	}
	if resp.StatusCode != 200 {
		var errorResponse apierrors.ErrorResponse
		if err := json.Unmarshal(resp.Body, &errorResponse); err != nil {
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, resp.Body)
		}

		if errorResponse.Type == apierrors.TypeValidationFailed {
			return nil, fmt.Errorf("%s: %s\n%v", errorResponse.Title, errorResponse.Detail, errorResponse.Errors)
		}
		return nil, fmt.Errorf("error sending request: %v", resp.Body)
	}

	if err != nil {
		fmt.Println("Error converting duration:", err)
		return nil, fmt.Errorf("error converting duration: %v", err)
	}
	validationDuration.Record(context.Background(), int64(time.Duration(time.Since(start).Milliseconds())))

	// Use nanosecond precision and format to microseconds to ensure unique timestamps
	now := time.Now()
	sensorData.Timestamp = now.Format("2006-01-02 15:04:05.") + fmt.Sprintf("%06d", now.Nanosecond()/1000)
	return &sensorData, nil
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

	meter = otel.Meter("neuro-lab.gateway")
	messageCounter, err = meter.Int64Counter(
		"mqtt.messages.received",
		metric.WithDescription("Number of MQTT messages received."),
		metric.WithUnit("{message}"),
	)
	if err != nil {
		log.Fatalf("failed to create message counter: %v", err)
	}

	gatewayDuration, err = meter.Int64Histogram(
		"gateway.duration",
		metric.WithDescription("Duration of gateway processing."),
		metric.WithUnit("ms"),
	)
	if err != nil {
		log.Fatalf("failed to create gateway duration histogram: %v", err)
	}

	validationDuration, err = meter.Int64Histogram(
		"validation.duration",
		metric.WithDescription("Duration of validation processing."),
		metric.WithUnit("ms"),
	)
	if err != nil {
		log.Fatalf("failed to create validation duration histogram: %v", err)
	}

	opts := mqtt.NewClientOptions().AddBroker("localhost:1884")
	opts.SetClientID("go_mqtt_client")

	topic := "gateway.raw"
	partition := 0
	conn, err := connect(topic, partition)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to Kafka: %v", err))
	}
	defer conn.Close()

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe
	client.Subscribe("device/+/raw", 0, func(client mqtt.Client, msg mqtt.Message) {
		start := time.Now()
		messageCounter.Add(ctx, 1)
		sensorData, err := processMessage(msg.Payload())
		if err != nil {
			fmt.Println("Error processing message:", err)
			return
		}
		fmt.Println("Sensor data:", sensorData)
		sensorData.FrameID = int(frameCounter.Add(1))
		sendViaKafka(conn, sensorData)
		duration := int64(time.Since(start).Milliseconds())
		fmt.Printf("Gateway duration: %dms\n", duration)
		gatewayDuration.Record(ctx, duration)
	})

	// Keep the client running
	time.Sleep(50000 * time.Second)
	client.Disconnect(250)
}
