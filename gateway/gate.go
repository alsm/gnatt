package main

import (
	"fmt"
	// TODO: if paho switches to github, use that repo instead
	MQTT "github.com/shoenig/go-mqtt"
)

type gate interface {
	run()
}

type aggGate struct {
	mqttclient *MQTT.MqttClient
}

func NewAggGate(opts *MQTT.ClientOptions) *aggGate {
	client := MQTT.NewClient(opts)
	ag := &aggGate{
		client,
	}
	return ag
}

func (ag *aggGate) run() {
	fmt.Println("Aggregating Gateway is running")
	ag.mqttclient.Start()
}

// type transGate struct {
// 	mqttclients []MQTT.MqttClient
// }
