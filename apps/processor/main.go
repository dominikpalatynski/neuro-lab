package main

import (
	"database"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

var (
	db       *gorm.DB
	frameIds *map[int]int
)

func init() {
	frameIds = &map[int]int{}
}

func main() {
	db = database.Connect()
	db.AutoMigrate(&database.ProcessedChannel{})
	opts := mqtt.NewClientOptions().AddBroker("localhost:1884")
	opts.SetClientID("go_mqtt_client")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go func() {
		client.Subscribe("device/+/raw", 0, processMessage)
	}()

	select {}
}
