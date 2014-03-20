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
	mqttbroker string
	username   string
	password   string
}

// Do not allow the creation of an MQTT-SN client if
// a connection to the MQTT broker cannot be established
func NewTransClient(id, mqttbroker string, c uConn, a uAddr) (*TransClient, error) {
	fmt.Printf("NewTransClient, id: \"%s\"\n", id)
	tc := &TransClient{
		Client{
			sync.RWMutex{},
			id,
			c,
			a,
			make(map[uint16]bool),
			make(map[uint16]*PublishMessage),
		},
		nil,
		mqttbroker,
		"",
		"",
	}
	if err := tc.connectMqtt(id, mqttbroker); err != nil {
		return nil, err
	}
	return tc, nil
}

func (tc *TransClient) connectMqtt(id, mqttbroker string) error {
	opts := MQTT.NewClientOptions()
	opts.SetBroker(mqttbroker)
	opts.SetClientId(id)
	if tc.username != "" {
		opts.SetUsername(tc.username)
		opts.SetPassword(tc.password)
	}
	opts.SetTraceLevel(MQTT.Warn)
	tc.mqttclient = MQTT.NewClient(opts)

	if _, err := tc.mqttclient.Start(); err != nil {
		fmt.Printf("SN client \"%s\" failed to connect to mqtt broker", tc.ClientId)
		return err
	}
	fmt.Println("TransClient connected to mqtt broker")
	return nil
}

func (tc *TransClient) disconnectMqtt() {
	tc.mqttclient.Disconnect(100)
}

func (tc *TransClient) subscribeMqtt(qos MQTT.QoS, topic string) {
	var handler MQTT.MessageHandler = func(msg MQTT.Message) {
		fmt.Printf("publish handler\n")

		// how to get topicid ?
		pm := NewPublishMessage(msg.DupFlag(), msg.RetainedFlag(), QoS(msg.QoS()), 0, 0, 0, msg.Payload())

		if err := tc.Write(pm); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("incoming mqtt published to mqtt-sn")
		}
	}

	if r, e := tc.mqttclient.StartSubscription(handler, topic, qos); e != nil {
		fmt.Printf("subscribe to \"%s\" failed: %s\n", topic, e)
	} else {
		<-r
		fmt.Printf("subscribe to \"%s\" succeeded\n", topic)
	}
}
