package gateway

import (
	"bytes"
	"net"
	"os"
	"sync"

	. "github.com/alsm/gnatt/packets"

	//MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

type TGateway struct {
	stopsig    chan os.Signal
	port       int
	mqttBroker string
	clients    Clients
	tIndex     topicNames
}

func NewTGateway(gc *GatewayConfig, stopsig chan os.Signal) *TGateway {
	t := &TGateway{
		stopsig,
		gc.port,
		gc.mqttbroker,
		Clients{
			sync.RWMutex{},
			make(map[string]SNClient),
		},
		topicNames{
			sync.RWMutex{},
			make(map[uint16]string),
			0,
		},
	}
	return t
}

func (t *TGateway) Port() int {
	return t.port
}

func (t *TGateway) Start() {
	go t.awaitStop()
	INFO.Println("Transparent Gataway is started")
	listen(t)
}

func (t *TGateway) awaitStop() {
	<-t.stopsig
	INFO.Println("Transparent Gateway is stopped")
	os.Exit(0)
}

func (t *TGateway) OnPacket(nbytes int, buffer []byte, con *net.UDPConn, addr *net.UDPAddr) {
	INFO.Println("TG OnPacket!")
	INFO.Printf("bytes: %s\n", string(buffer[0:nbytes]))

	buf := bytes.NewBuffer(buffer)
	rawmsg, _ := ReadPacket(buf)

	INFO.Printf("rawmsg.MessageType(): %s\n", MessageNames[rawmsg.MessageType()])

	switch msg := rawmsg.(type) {
	case *AdvertiseMessage:
		t.handle_ADVERTISE(msg, addr)
	case *SearchGwMessage:
		t.handle_SEARCHGW(msg, addr)
	case *GwInfoMessage:
		t.handle_GWINFO(msg, addr)
	case *ConnectMessage:
		t.handle_CONNECT(msg, con, addr)
	case *ConnackMessage:
		t.handle_CONNACK(msg, addr)
	case *WillTopicReqMessage:
		t.handle_WILLTOPICREQ(msg, addr)
	case *WillTopicMessage:
		t.handle_WILLTOPIC(msg, addr)
	case *WillMsgReqMessage:
		t.handle_WILLMSGREQ(msg, addr)
	case *WillMsgMessage:
		t.handle_WILLMSG(msg, addr)
	case *RegisterMessage:
		t.handle_REGISTER(msg, con, addr)
	case *RegackMessage:
		t.handle_REGACK(msg, addr)
	case *PublishMessage:
		t.handle_PUBLISH(msg, addr)
	case *PubackMessage:
		t.handle_PUBACK(msg, addr)
	case *PubcompMessage:
		t.handle_PUBCOMP(msg, addr)
	case *PubrecMessage:
		t.handle_PUBREC(msg, addr)
	case *PubrelMessage:
		t.handle_PUBREL(msg, addr)
	case *SubscribeMessage:
		t.handle_SUBSCRIBE(msg, addr)
	case *SubackMessage:
		t.handle_SUBACK(msg, addr)
	case *UnsubackMessage:
		t.handle_UNSUBACK(msg, addr)
	case *PingreqMessage:
		t.handle_PINGREQ(msg, con, addr)
	case *DisconnectMessage:
		t.handle_DISCONNECT(msg, addr)
	case *WillTopicUpdateMessage:
		t.handle_WILLTOPICUPD(msg, addr)
	case *WillTopicRespMessage:
		t.handle_WILLTOPICRESP(msg, addr)
	case *WillMsgUpdateMessage:
		t.handle_WILLMSGUPD(msg, addr)
	case *WillMsgRespMessage:
		t.handle_WILLMSGRESP(msg, addr)
	default:
		ERROR.Printf("Unknown Message Type %T\n", msg)
	}
}

func (t *TGateway) handle_ADVERTISE(m *AdvertiseMessage, a *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], a)
}

func (t *TGateway) handle_SEARCHGW(m *SearchGwMessage, a *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], a)
}

func (t *TGateway) handle_GWINFO(m *GwInfoMessage, a *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], a)
}

func (t *TGateway) handle_CONNECT(m *ConnectMessage, c *net.UDPConn, a *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], a)
	INFO.Println(m.ProtocolId, m.Duration, m.ClientId)
	if clientid, err := validateClientId(m.ClientId); err != nil {
		ERROR.Println(err)
	} else {
		INFO.Printf("clientid: %s\n", clientid)
		INFO.Printf("remoteaddr: %s\n", a)
		INFO.Printf("will: %v\n", m.Will)
		if m.Will {
			// todo: will msg
		}
		if tClient, err := NewTClient(string(clientid), t.mqttBroker, c, a); err != nil {
			ERROR.Println(err)
		} else {
			t.clients.AddClient(tClient)

			// establish connection to mqtt broker

			ca := NewMessage(CONNACK).(*ConnackMessage)
			ca.ReturnCode = 0
			if err = tClient.Write(ca); err != nil {
				ERROR.Println(err)
			} else {
				INFO.Println("CONNACK was sent")
			}
		}
	}
}

func (t *TGateway) handle_CONNACK(m *ConnackMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_WILLTOPICREQ(m *WillTopicReqMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_WILLTOPIC(m *WillTopicMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_WILLMSGREQ(m *WillMsgReqMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_WILLMSG(m *WillMsgMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_REGISTER(m *RegisterMessage, c *net.UDPConn, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
	topic := string(m.TopicName)
	var topicid uint16
	if !t.tIndex.containsTopic(topic) {
		topicid = t.tIndex.putTopic(topic)
	} else {
		topicid = t.tIndex.getId(topic)
	}

	INFO.Printf("t topicid: %d\n", topicid)

	tclient := t.clients.GetClient(r).(*TClient)
	tclient.Register(topicid, topic)

	ra := NewRegackMessage(topicid, m.MessageId, 0)
	INFO.Printf("ra.Msgid: %d\n", ra.MessageId)

	if err := tclient.Write(ra); err != nil {
		ERROR.Println(err)
	} else {
		INFO.Println("REGACK sent")
	}
}

func (t *TGateway) handle_REGACK(m *RegackMessage, a *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], a)
}

func (t *TGateway) handle_PUBLISH(m *PublishMessage, a *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], a)
	tclient := t.clients.GetClient(a).(*TClient)

	topic := t.tIndex.getTopic(m.TopicId)

	INFO.Println(topic, m.Qos, m.Retain, m.Data)
	if token := tclient.mqttClient.Publish(topic, m.Qos, m.Retain, m.Data); token.WaitTimeout(2000) && token.Error() != nil {
		ERROR.Println("Error publishing message", token.Error())
		return
	}
	INFO.Println("PUBLISH published")
}

func (t *TGateway) handle_PUBACK(m *PubackMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_PUBCOMP(m *PubcompMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_PUBREC(m *PubrecMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_PUBREL(m *PubrelMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_SUBSCRIBE(m *SubscribeMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
	topic := ""
	if m.TopicIdType == 0 { // todo: other topic id types, also use enum
		topic = string(m.TopicName)
	} else {
		ERROR.Println("other topic id types not supported yet")
		topic = "not_implemented"
	}
	tclient := t.clients.GetClient(r).(*TClient)
	INFO.Printf("subscribe, qos: %d, topic: %s\n", m.Qos, topic)
	tclient.subscribeMQTT(m.Qos, topic, &t.tIndex)

	suba := NewSubackMessage(0, m.MessageId, m.Qos, 0)

	if err := tclient.Write(suba); err != nil {
		ERROR.Println(err)
	} else {
		INFO.Println("SUBACK sent")
	}
}

func (t *TGateway) handle_SUBACK(m *SubackMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_UNSUBSCRIBE(m *UnsubscribeMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_UNSUBACK(m *UnsubackMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_PINGREQ(m *PingreqMessage, c *net.UDPConn, a *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], a)
	tclient := t.clients.GetClient(a).(*TClient)

	resp := NewMessage(PINGRESP)

	if err := tclient.Write(resp); err != nil {
		ERROR.Println(err)
	} else {
		INFO.Println("PINGRESP sent")
	}
}

func (t *TGateway) handle_PINGRESP(m *PingrespMessage, a *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], a)
}

func (t *TGateway) handle_DISCONNECT(m *DisconnectMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
	tclient := t.clients.GetClient(r).(*TClient)
	tclient.disconnectMQTT()
	t.clients.RemoveClient(tclient.ClientId)
}

func (t *TGateway) handle_WILLTOPICUPD(m *WillTopicUpdateMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_WILLTOPICRESP(m *WillTopicRespMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_WILLMSGUPD(m *WillMsgUpdateMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}

func (t *TGateway) handle_WILLMSGRESP(m *WillMsgRespMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", MessageNames[m.MessageType()], r)
}
