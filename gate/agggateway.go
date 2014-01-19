package gateway

import (
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"os"
	"time"
)

type AggGate struct {
	udpListener
	mqttclient *MQTT.MqttClient
	stopsig    chan os.Signal
}

func NewAggGate(opts *MQTT.ClientOptions, stopsig chan os.Signal, port int) *AggGate {
	listener := udpListener{
		port,
	}
	client := MQTT.NewClient(opts)
	ag := &AggGate{
		listener,
		client,
		stopsig,
	}
	return ag
}

func (ag *AggGate) Start() {
	go ag.awaitStop()
	fmt.Println("Aggregating Gateway is starting")
	_, err := ag.mqttclient.Start()
	if err != nil {
		fmt.Println("Aggregating Gateway failed to start")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Aggregating Gateway is started")
	ag.listen()
}

// This does NOT WORK on Windows using Cygwin, however
// it does work using cmd.exe
func (ag *AggGate) awaitStop() {
	<-ag.stopsig
	fmt.Println("Aggregating Gateway is stopping")
	ag.mqttclient.Disconnect(500)
	time.Sleep(500) //give broker some time to process DISCONNECT
	fmt.Println("Aggregating Gateway is stopped")

	// TODO: cleanly close down other goroutines

	os.Exit(0)
}
