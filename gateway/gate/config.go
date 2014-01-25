package gateway

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type GatewayConfig struct {
	aggregating  bool
	port         int
	mqttbroker   string
	mqttuser     string
	mqttpassword string
	mqttclientid string
	mqtttimeout  int
}

func (gc *GatewayConfig) IsAggregating() bool {
	return gc.aggregating
}

func ParseConfigFile(file string) (*GatewayConfig, error) {
	gc := &GatewayConfig{}
	if bytes, rerr := ioutil.ReadFile(file); rerr != nil {
		return nil, rerr
	} else {
		if perr := gc.parseConfig(string(bytes)); perr != nil {
			return nil, perr
		}
	}
	return gc, nil
}

func (gc *GatewayConfig) parseConfig(config string) error {
	scanner := bufio.NewScanner(bytes.NewReader([]byte(config)))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if k, v, e := gc.parseLine(line); e != nil {
			return e
		} else if k == "" && v == "" {
			// skipping comment or blank line
		} else {
			if e = gc.setOption(k, v); e != nil {
				return e
			}
		}
	}
	return nil
}

func (gc *GatewayConfig) parseLine(line string) (string, string, error) {
	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] == '#' {
		return "", "", nil
	}
	fields := strings.Fields(line)
	if len(fields) == 1 {
		return "", "", fmt.Errorf("Missing value for config option: \"%s\"", fields[0])
	} else if len(fields) > 2 {
		return "", "", fmt.Errorf("Too many values supplied for config option: \"%s\"", fields[0])
	}
	return fields[0], fields[1], nil
}

func (gc *GatewayConfig) setOption(key, value string) error {
	var e error
	switch key {
	case "type":
		gc.aggregating, e = checkType(value)
	case "port":
		gc.port, e = checkNum("port", value)
	case "mqtt-broker":
		gc.mqttbroker, e = checkURI(value)
	case "mqtt-user":
		gc.mqttuser = value
	case "mqtt-password":
		gc.mqttpassword = value
	case "mqtt-clientid":
		gc.mqttclientid = value
	case "mqtt-timeout":
		gc.mqtttimeout, e = checkNum("mqtt-timeout", value)
	default:
		return fmt.Errorf("Unknown config option: \"%s\"", key)
	}
	return e
}

func checkURI(value string) (string, error) {
	if value[0:6] != "tcp://" &&
		value[0:6] != "ssl://" &&
		value[0:7] != "tcps://" {
		return "", fmt.Errorf("Invalid URI, must specify transport (ex: \"tcp://\"): \"%s\"", value)
	}
	// todo: check that a port is provided
	// also there is probably a library way to verify a URI
	return value, nil
}

func checkType(value string) (bool, error) {
	var isAggregating bool
	switch value {
	case "aggregating":
		isAggregating = true
	case "transparent":
		isAggregating = false
	default:
		return false, fmt.Errorf("Invalid value specified for \"type\": \"%s\"", value)
	}
	return isAggregating, nil
}

func checkNum(label, value string) (int, error) {
	if p, e := strconv.Atoi(value); e != nil {
		return 0, fmt.Errorf("Invalid value specified for \"%s\" (not a number): \"%s\"", label, value)
	} else {
		return p, nil
	}
}
