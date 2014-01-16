package main

import (
	"flag"
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"os"
	"os/signal"
)

func main() {
	stopsig := registerSignals()
	aggregating, mqttBroker, port := setup()
	var gateway gateway
	if aggregating {
		fmt.Println("GNATT Gateway starting in aggregating mode")
		gateway = initAggregating(mqttBroker, stopsig, port)
	} else {
		fmt.Println("GNATT Transparent gateway not yet implemented")
		os.Exit(0)
		//fmt.Println("GNATT Gateway starting in transparent mode")
		//gateway = initTransparent(broker)
	}

	gateway.start()
}

func setup() (bool, string, int) {
	var aggregating bool
	var broker string
	var udpport int
	flag.BoolVar(&aggregating, "aggregating", false, "Transparent or Aggregating")
	flag.StringVar(&broker, "broker", "", "MQTT Broker URI")
	flag.IntVar(&udpport, "port", 0, "MQTT-SN UDP Listening Port")
	flag.Parse()
	if broker == "" {
		fmt.Println("Must specify -broker <tcp://host:port>")
		os.Exit(1)
	}
	if udpport == 0 {
		fmt.Println("Must specify -port <port>")
		os.Exit(1)
	}

	return aggregating, broker, udpport
}

func registerSignals() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}

func initAggregating(broker string, stopsig chan os.Signal, port int) *aggGate {
	opts := MQTT.NewClientOptions()
	opts.SetBroker(broker)
	ag := NewAggGate(opts, stopsig, port)
	return ag
}

// func initTransparent(broker string) *transGate {
// 	fmt.Printf("trans broker: %s\n", broker)
// }
