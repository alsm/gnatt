package gateway

import (
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

type TransCon struct {
	mqttclient *MQTT.MqttClient
}
