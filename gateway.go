package main

import (
	"flag"
	"fmt"
	SN "github.com/alsm/gnatt/gate"
	"os"
	"os/signal"
)

func main() {
	var gateway SN.Gateway
	stopsig := registerSignals()
	gatewayconf := setup()

	if gatewayconf.IsAggregating() {
		fmt.Println("GNATT Gateway starting in aggregating mode")
		gateway = initAggregating(gatewayconf, stopsig)
	} else {
		fmt.Println("GNATT Transparent gateway not yet implemented")
		os.Exit(0)
		//fmt.Println("GNATT Gateway starting in transparent mode")
		//gateway = initTransparent(broker)
	}

	gateway.Start()
}

func setup() *SN.GatewayConfig {
	var configFile string
	var port int

	flag.StringVar(&configFile, "configuration", "", "Configuration File")
	flag.IntVar(&port, "port", 0, "MQTT-SN UDP Listening Port")
	flag.Parse()

	if configFile != "" {
		if gc, err := SN.ParseConfigFile(configFile); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			return gc
		}
	}

	fmt.Println("-configuration <file> must be specified")
	os.Exit(1)
	return nil
}

func initAggregating(gc *SN.GatewayConfig, stopsig chan os.Signal) *SN.AggGate {
	ag := SN.NewAggGate(gc, stopsig)
	return ag
}

// func initTransparent(broker string) *transGate {
// 	fmt.Printf("trans broker: %s\n", broker)
// }

func registerSignals() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}
