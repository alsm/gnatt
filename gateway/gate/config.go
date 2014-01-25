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
	switch key {
	case "type":
		switch value {
		case "aggregating":
			gc.aggregating = true
		case "transparent":
			gc.aggregating = false
		default:
			return fmt.Errorf("Invalid value specified for \"type\": \"%s\"", value)
		}
		return nil
	case "port":
		if p, e := strconv.Atoi(value); e != nil {
			return fmt.Errorf("Invalid value specified for \"port\" (not a number): \"%s\"", value)
		} else {
			gc.port = p
		}
		return nil
	case "mqtt-broker":
		gc.mqttbroker = value
		return nil
	case "mqtt-user":
		gc.mqttuser = value
		return nil
	case "mqtt-password":
		gc.mqttpassword = value
		return nil
	case "mqtt-clientid":
		gc.mqttclientid = value
		return nil
	default:
		return fmt.Errorf("Unknown config option: \"%s\"", key)
	}
}
