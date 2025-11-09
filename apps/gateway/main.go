package main

import (
	"fmt"
	"time"

	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	kafka "github.com/segmentio/kafka-go"
)

// Connect to the specified topic and partition in the server
func connect(topic string, partition int) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp",
		"localhost:19092", topic, partition)
	if err != nil {
		fmt.Println("failed to dial leader", err)
	}
	return conn, err
} //end connect

func sendViaKafka(conn *kafka.Conn, msg string) {
	// Set write deadline to ensure timely sending
	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))

	_, err := conn.WriteMessages(
		kafka.Message{Value: []byte(msg)})
	if err != nil {
		fmt.Println("failed to write messages:", err)
		return
	}

	// Flush to send messages immediately (real-time delivery)
	if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		fmt.Println("failed to set deadline for flush:", err)
	}
}

func main() {
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
		fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), msg.Payload())
		sendViaKafka(conn, string(msg.Payload()))
	})

	// Keep the client running
	time.Sleep(50000 * time.Second)
	client.Disconnect(250)
}
