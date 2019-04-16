package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)


type mqttMessage struct {
	topic string
	payload string
}


var con = NewController()

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	if strings.Contains(msg.Topic(), "sensor") {
		// fmt.Println("Received: " + msg.Topic() + " - " + string(msg.Payload()))

		sensorName := msg.Topic()[1:]
		con.SetSensorState(sensorName, string(msg.Payload()) == "1")
	}
}

func main() {
	quit := make(chan os.Signal, 1)
	mc := make(chan []mqttMessage)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	opts := MQTT.NewClientOptions().AddBroker("ssl://broker.0f.nl:8883")
	opts.SetClientID("go-test")
	opts.SetDefaultPublishHandler(f)
	opts.SetTLSConfig(&tls.Config{})
	topic := "5/#"

	opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}

	go con.Loop(mc)

	for {
		select {
			case msgs := <- mc:
				for _, msg := range msgs {
					client.Publish(msg.topic, byte(1), false, msg.payload)
				}
			case <- quit:
				return
		}
	}
}