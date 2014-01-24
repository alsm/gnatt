package gateway

import (
	"bufio"
	"fmt"
	"io/ioutil"
)

type GatewayConfig struct {
	aggregating bool
	udpport     int
	mqttbroker  string
	mqttuser    string
	mqttpass    string
	mqttcid     string
}

func ParseConfigFile(file string) (*GatewayConfig, error) {
	if bytes, rerr := ioutil.ReadFile(file); rerr != nil {
		return nil, rerr
	} else {
		if config, perr := parseConfig(string(bytes)); perr != nil {
			return nil, perr
		} else {
			return config, nil
		}
	}
}

func parseConfig(config string) (*GatewayConfig, error) {
	i, lines, err := bufio.ScanLines([]byte(config), true)
	if err != nil {
		return nil, err
	}

	fmt.Printf("i: %d\n", i)
	fmt.Printf("len(lines): %d\n", len(lines))

	return nil, nil
}
