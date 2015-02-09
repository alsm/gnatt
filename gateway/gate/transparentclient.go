package gateway

import (
	"net"
	"sync"

	. "github.com/alsm/gnatt/packets"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

type TClient struct {
	Client
	mqttClient *MQTT.Client
	mqttBroker string
	username   string
	password   string
}

// Do not allow the creation of an MQTT-SN client if
// a connection to the MQTT broker cannot be established
func NewTClient(ClientId, Broker string, Connection *net.UDPConn, Address *net.UDPAddr) (*TClient, error) {
	INFO.Println("NewTClient, id: %s", ClientId)
	t := &TClient{
		Client{
			sync.RWMutex{},
			ClientId,
			Connection,
			Address,
			make(map[uint16]string),
			make(map[uint16]*PublishMessage),
		},
		nil,
		Broker,
		"",
		"",
	}
	if err := t.connectMQTT(ClientId, Broker); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *TClient) connectMQTT(ClientId, Broker string) error {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(Broker)
	opts.SetClientID(ClientId)
	if t.username != "" {
		opts.SetUsername(t.username)
		opts.SetPassword(t.password)
	}
	t.mqttClient = MQTT.NewClient(opts)

	if token := t.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	INFO.Println("TClient connected to mqtt broker")
	return nil
}

func (t *TClient) disconnectMQTT() {
	t.mqttClient.Disconnect(100)
}

func (t *TClient) subscribeMQTT(qos byte, topic string, tIndex *topicNames) {
	var handler MQTT.MessageHandler = func(client *MQTT.Client, msg MQTT.Message) {
		INFO.Println("publish handler")

		tid := tIndex.getId(msg.Topic())
		// is topicid type always 0 coming out of tIndex ?
		// todo: msgid is not always 0
		pm := NewPublishMessage(tid, 0x00, msg.Payload(), msg.Qos(), 0x00, msg.Retained(), msg.Duplicate())

		if err := t.Write(pm); err != nil {
			ERROR.Println(err)
		} else {
			INFO.Println("incoming mqtt published to mqtt-sn")
		}
	}

	if token := t.mqttClient.Subscribe(topic, qos, handler); token.WaitTimeout(2000) && token.Error() != nil {
		ERROR.Println("Error subscribing,", token.Error())
	}
	INFO.Println(t.ClientId, "subscribed to", topic)
}
