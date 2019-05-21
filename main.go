package main

import (
	"crypto/tls"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type mqttMessage struct {
	topic   string
	payload string
}

var ch = make(chan TrafficChange)
var con = NewController(ch)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	if strings.Contains(msg.Topic(), "sensor") {
		fmt.Println("Received: " + msg.Topic() + " - " + string(msg.Payload()))

		sensorName := msg.Topic()[1:]
		con.SetSensorState(sensorName, string(msg.Payload()) == "1")
	}

	if strings.Contains(msg.Topic(), "ondisconnect") {
		con.SetTrafficItemsInitialState()
	}
}

func main() {
	quit := make(chan os.Signal, 1)
	mc := make(chan []mqttMessage)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	opts := MQTT.NewClientOptions().AddBroker("ssl://broker.0f.nl:8883")
	opts.SetClientID("")
	opts.SetDefaultPublishHandler(f)
	opts.SetTLSConfig(&tls.Config{})
	opts.SetWill("5/features/lifecycle/controller/ondisconnect/", "", byte(1), false)
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

	client.Publish("5/features/lifecycle/controller/onconnect", byte(1), false, "")

	go handleTrafficChanges(con, ch, quit)
	go con.Loop(mc)

	for {
		select {
		case msgs := <-mc:
			for _, msg := range msgs {
				fmt.Println("Sent:", msg.topic, "-", msg.payload)
				client.Publish(msg.topic, byte(1), false, msg.payload)
			}
		case <-quit:
			return
		}
	}
}
