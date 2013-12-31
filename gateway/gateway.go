package main

import (
	"flag"
	"fmt"
	"os"
	// TODO: if paho switches to github, use that repo instead
	MQTT "github.com/shoenig/go-mqtt"
)

func main() {
	aggregating, broker := setup()
	var gateway gate
	if aggregating {
		fmt.Println("GNATT Gateway starting in aggregating mode")
		gateway = initAggregating(broker)
	} else {
		fmt.Println("GNATT Transparent gateway not yet implemented")
		os.Exit(0)
		//fmt.Println("GNATT Gateway starting in transparent mode")
		//gateway = initTransparent(broker)
	}

	gateway.run()
}

func setup() (bool, string) {
	var aggregating bool
	var broker string
	flag.BoolVar(&aggregating, "aggregating", false, "Transparent or Aggregating")
	flag.StringVar(&broker, "broker", "", "MQTT Broker")
	flag.Parse()
	if broker == "" {
		fmt.Println("Must specify broker")
		os.Exit(1)
	}
	return aggregating, broker
}

func initAggregating(broker string) *aggGate {
	fmt.Printf("ag broker: %s\n", broker)
	opts := MQTT.NewClientOptions()
	opts.SetBroker(broker)
	ag := NewAggGate(opts)
	return ag
}

// func initTransparent(broker string) *transGate {
// 	fmt.Printf("trans broker: %s\n", broker)
// }
