package gateway

import (
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	. "github.com/alsm/gnatt/common/protocol"
	"net"
	"os"
	"time"
)

type AggGate struct {
	mqttclient *MQTT.MqttClient
	stopsig    chan os.Signal
	port       int
}

func NewAggGate(opts *MQTT.ClientOptions, stopsig chan os.Signal, port int) *AggGate {
	client := MQTT.NewClient(opts)
	ag := &AggGate{
		client,
		stopsig,
		port,
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
	listen(ag)
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

func (ag *AggGate) OnPacket(nbytes int, buffer []byte, remote *net.UDPAddr) {
	fmt.Println("OnPacket!")

	var h Header
	h.UnpackHeader(buffer)

	fmt.Printf("h.Length: %d\n", h.Length())
	fmt.Printf("h.msgType: %s\n", h.MsgType())
}

func (ag *AggGate) Port() int {
	return ag.port
}
