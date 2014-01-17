package main

import (
	"flag"
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	SN "github.com/alsm/gnatt/gate"
	"os"
	"os/signal"
)

func main() {
	stopsig := registerSignals()
	aggregating, port, mqttopts := setup()
	var gateway SN.Gateway
	if aggregating {
		fmt.Println("GNATT Gateway starting in aggregating mode")
		gateway = initAggregating(port, stopsig, mqttopts)
	} else {
		fmt.Println("GNATT Transparent gateway not yet implemented")
		os.Exit(0)
		//fmt.Println("GNATT Gateway starting in transparent mode")
		//gateway = initTransparent(broker)
	}

	gateway.Start()
}

func setup() (bool, int, [4]string) {
	var aggregating bool
	var udpport int
	var mqttbroker string
	var mqttuser string
	var mqttpass string
	var mqttcid string

	flag.BoolVar(&aggregating, "aggregating", false, "Transparent or Aggregating")
	flag.IntVar(&udpport, "port", 0, "MQTT-SN UDP Listening Port")
	flag.StringVar(&mqttbroker, "mqtt-broker", "", "MQTT Broker URI")
	flag.StringVar(&mqttuser, "mqtt-user", "", "MQTT User")
	flag.StringVar(&mqttpass, "mqtt-password", "", "MQTT Password")
	flag.StringVar(&mqttcid, "mqtt-clientid", "gnatt-gateway", "MQTT Client ID")

	flag.Parse()
	if mqttbroker == "" {
		fmt.Println("Must specify -mqtt-broker <tcp://host:port>")
		os.Exit(1)
	}

	if udpport == 0 {
		fmt.Println("Must specify -port <port>")
		os.Exit(1)
	}

	opts := [4]string{mqttbroker, mqttuser, mqttpass, mqttcid}

	return aggregating, udpport, opts
}

func registerSignals() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}

func initAggregating(port int, stopsig chan os.Signal, mqttopts [4]string) *SN.AggGate {
	opts := MQTT.NewClientOptions()
	opts.SetBroker(mqttopts[0])
	if mqttopts[1] != "" {
		opts.SetUsername(mqttopts[1])
	}
	if mqttopts[2] != "" {
		opts.SetPassword(mqttopts[2])
	}
	if mqttopts[3] != "" {
		opts.SetClientId(mqttopts[3])
	}
	ag := SN.NewAggGate(opts, stopsig, port)
	return ag
}

// func initTransparent(broker string) *transGate {
// 	fmt.Printf("trans broker: %s\n", broker)
// }
