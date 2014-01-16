package main

import (
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"os"
)

type aggGate struct {
	udpListener
	mqttclient *MQTT.MqttClient
	stopsig    chan os.Signal
}

func NewAggGate(opts *MQTT.ClientOptions, stopsig chan os.Signal, port int) *aggGate {
	listener := udpListener{
		port,
	}
	client := MQTT.NewClient(opts)
	ag := &aggGate{
		listener,
		client,
		stopsig,
	}
	return ag
}

func (ag *aggGate) start() {
	go ag.awaitStop()
	fmt.Println("Aggregating Gateway is starting")
	ag.mqttclient.Start()
	fmt.Println("Aggregating Gateway is started")
	ag.listen()
}

// This does NOT WORK on Windows using Cygwin, however
// it does work using cmd.exe
func (ag *aggGate) awaitStop() {
	<-ag.stopsig
	fmt.Println("Aggregating Gateway is stopping")

	fmt.Println("Aggregating Gateway is stopped")

	// TODO: cleanly close down other goroutines

	os.Exit(0)
}
