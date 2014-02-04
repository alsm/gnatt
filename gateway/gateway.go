package main

import (
	"flag"
	"fmt"
	G "github.com/alsm/gnatt/gateway/gate"
	"os"
	"os/signal"
)

func main() {
	var gateway G.Gateway
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

func setup() *G.GatewayConfig {
	var configFile string
	var port int

	flag.StringVar(&configFile, "c", "", "Configuration File")
	flag.IntVar(&port, "port", 0, "MQTT-G UDP Listening Port")
	flag.Parse()

	if configFile != "" {
		if gc, err := G.ParseConfigFile(configFile); err != nil {
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

func initAggregating(gc *G.GatewayConfig, stopsig chan os.Signal) *G.AggGate {
	ag := G.NewAggGate(gc, stopsig)
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
