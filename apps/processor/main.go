package main

import (
	"database"

	"context"
	"opentelemetry"

	"time"

	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"gorm.io/gorm"
)

var (
	db              *gorm.DB
	frameIds        *map[int]int
	meter           metric.Meter
	gatewayDuration metric.Int64Histogram
)

func init() {
	frameIds = &map[int]int{}
}

func main() {
	db = database.Connect()
	db.AutoMigrate(&database.ProcessedChannel{})
	opts := mqtt.NewClientOptions().AddBroker("192.168.18.23:31095")
	opts.SetClientID("go_mqtt_client")
	ctx := context.Background()
	meter = otel.Meter("neuro-lab.processor")
	var err error
	gatewayDuration, err = meter.Int64Histogram(
		"processor.duration",
		metric.WithDescription("Duration of processor processing."),
		metric.WithUnit("ms"),
	)
	if err != nil {
		panic(err)
	}
	otelSdkSetup := opentelemetry.NewOtelSdkSetup(opentelemetry.SetupOTelSDKOptions{
		ServiceName:    "processor",
		ServiceVersion: "0.1.0",
	})
	otelShutdown, err := otelSdkSetup.Setup(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := otelShutdown(ctx); err != nil {
			panic(err)
		}
	}()
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go func() {
		client.Subscribe("device/+/raw", 0, func(client mqtt.Client, msg mqtt.Message) {
			start := time.Now()
			processMessage(client, msg)
			duration := int64(time.Since(start).Milliseconds())
			fmt.Printf("Processor duration: %dms\n", duration)
			gatewayDuration.Record(ctx, duration)
		})
	}()

	select {}
}
