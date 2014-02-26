package gateway

import (
	"fmt"
	"sync"

	. "github.com/alsm/gnatt/common/protocol"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

type TransClient struct {
	Client
	mqttclient *MQTT.MqttClient
}

func NewTransClient(id string, c uConn, a uAddr) *TransClient {
	fmt.Printf("NewTransClient, id: \"%s\"\n", id)
	return &TransClient{
		Client{
			sync.RWMutex{},
			id,
			c,
			a,
			make(map[uint16]bool),
			make(map[uint16]*PublishMessage),
		},
		nil,
	}
}
