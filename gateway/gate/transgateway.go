package gateway

import (
	"os"
	"sync"

	. "github.com/alsm/gnatt/common/protocol"
	"github.com/alsm/gnatt/common/utils"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

type TransGate struct {
	stopsig    chan os.Signal
	port       int
	mqttbroker string
	clients    Clients
	tIndex     topicNames
}

func NewTransGate(gc *GatewayConfig, stopsig chan os.Signal) *TransGate {
	tg := &TransGate{
		stopsig,
		gc.port,
		gc.mqttbroker,
		Clients{
			sync.RWMutex{},
			make(map[string]StorableClient),
		},
		topicNames{
			sync.RWMutex{},
			make(map[uint16]string),
			0,
		},
	}
	return tg
}

func (tg *TransGate) Port() int {
	return tg.port
}

func (tg *TransGate) Start() {
	go tg.awaitStop()
	INFO.Println("Transparent Gataway is started")
	listen(tg)
}

func (tg *TransGate) awaitStop() {
	<-tg.stopsig
	INFO.Println("Transparent Gateway is stopped")
	os.Exit(0)
}

func (tg *TransGate) OnPacket(nbytes int, buffer []byte, con uConn, addr uAddr) {
	INFO.Println("TG OnPacket!")
	INFO.Printf("bytes: %s\n", utils.Bytes2str(buffer[0:nbytes]))
	rawmsg := Unpack(buffer[0:nbytes])
	INFO.Printf("rawmsg.MsgType(): %s\n", rawmsg.MsgType())

	switch msg := rawmsg.(type) {
	case *AdvertiseMessage:
		tg.handle_ADVERTISE(msg, addr)
	case *SearchGwMessage:
		tg.handle_SEARCHGW(msg, addr)
	case *GwInfoMessage:
		tg.handle_GWINFO(msg, addr)
	case *ConnectMessage:
		tg.handle_CONNECT(msg, con, addr)
	case *ConnackMessage:
		tg.handle_CONNACK(msg, addr)
	case *WillTopicReqMessage:
		tg.handle_WILLTOPICREQ(msg, addr)
	case *WillTopicMessage:
		tg.handle_WILLTOPIC(msg, addr)
	case *WillMsgReqMessage:
		tg.handle_WILLMSGREQ(msg, addr)
	case *WillMsgMessage:
		tg.handle_WILLMSG(msg, addr)
	case *RegisterMessage:
		tg.handle_REGISTER(msg, con, addr)
	case *RegackMessage:
		tg.handle_REGACK(msg, addr)
	case *PublishMessage:
		tg.handle_PUBLISH(msg, addr)
	case *PubackMessage:
		tg.handle_PUBACK(msg, addr)
	case *PubcompMessage:
		tg.handle_PUBCOMP(msg, addr)
	case *PubrecMessage:
		tg.handle_PUBREC(msg, addr)
	case *PubrelMessage:
		tg.handle_PUBREL(msg, addr)
	case *SubscribeMessage:
		tg.handle_SUBSCRIBE(msg, addr)
	case *SubackMessage:
		tg.handle_SUBACK(msg, addr)
	case *UnsubackMessage:
		tg.handle_UNSUBACK(msg, addr)
	case *PingreqMessage:
		tg.handle_PINGREQ(msg, con, addr)
	case *DisconnectMessage:
		tg.handle_DISCONNECT(msg, addr)
	case *WillTopicUpdateMessage:
		tg.handle_WILLTOPICUPD(msg, addr)
	case *WillTopicRespMessage:
		tg.handle_WILLTOPICRESP(msg, addr)
	case *WillMsgUpdateMessage:
		tg.handle_WILLMSGUPD(msg, addr)
	case *WillMsgRespMessage:
		tg.handle_WILLMSGRESP(msg, addr)
	default:
		ERROR.Printf("Unknown Message Type %T\n", msg)
	}
}

func (tg *TransGate) handle_ADVERTISE(m *AdvertiseMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_SEARCHGW(m *SearchGwMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_GWINFO(m *GwInfoMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_CONNECT(m *ConnectMessage, c uConn, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	if clientid, err := validateClientId(m.ClientId()); err != nil {
		ERROR.Println(err)
	} else {
		INFO.Printf("clientid: %s\n", clientid)
		INFO.Printf("remoteaddr: %s\n", r.r)
		INFO.Printf("will: %v\n", m.Will())
		if m.Will() {
			// todo: will msg
		}
		if tclient, err := NewTransClient(string(clientid), tg.mqttbroker, c, r); err != nil {
			ERROR.Println(err)
		} else {
			tg.clients.AddClient(tclient)

			// establish connection to mqtt broker

			ca := NewConnackMessage(0) // todo: 0 ?
			if ioerr := tclient.Write(ca); ioerr != nil {
				ERROR.Println(ioerr)
			} else {
				INFO.Println("CONNACK was sent")
			}
		}
	}
}

func (tg *TransGate) handle_CONNACK(m *ConnackMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLTOPICREQ(m *WillTopicReqMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLTOPIC(m *WillTopicMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLMSGREQ(m *WillMsgReqMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLMSG(m *WillMsgMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_REGISTER(m *RegisterMessage, c uConn, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	topic := string(m.TopicName())
	var topicid uint16
	if !tg.tIndex.containsTopic(topic) {
		topicid = tg.tIndex.putTopic(topic)
	} else {
		topicid = tg.tIndex.getId(topic)
	}

	INFO.Printf("tg topicid: %d\n", topicid)

	tclient := tg.clients.GetClient(r).(*TransClient)
	tclient.Register(topicid)

	ra := NewRegackMessage(topicid, m.MsgId(), 0)
	INFO.Printf("ra.Msgid: %d\n", ra.MsgId())

	if err := tclient.Write(ra); err != nil {
		ERROR.Println(err)
	} else {
		INFO.Println("REGACK sent")
	}
}

func (tg *TransGate) handle_REGACK(m *RegackMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_PUBLISH(m *PublishMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	tclient := tg.clients.GetClient(r).(*TransClient)

	topic := tg.tIndex.getTopic(m.TopicId())

	receipt := tclient.mqttclient.Publish(MQTT.QoS(m.QoS()), topic, m.Data())

	<-receipt
	INFO.Println("PUBLISH published")
}

func (tg *TransGate) handle_PUBACK(m *PubackMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_PUBCOMP(m *PubcompMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_PUBREC(m *PubrecMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_PUBREL(m *PubrelMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_SUBSCRIBE(m *SubscribeMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	topic := ""
	if m.TopicIdType() == 0 { // todo: other topic id types, also use enum
		topic = string(m.TopicName())
	} else {
		ERROR.Println("other topic id types not supported yet")
		topic = "not_implemented"
	}
	tclient := tg.clients.GetClient(r).(*TransClient)
	INFO.Printf("subscribe, qos: %d, topic: %s\n", m.QoS(), topic)
	tclient.subscribeMqtt(MQTT.QoS(m.QoS()), topic, &tg.tIndex)

	su := NewSubackMessage(0, m.QoS(), 0, m.MsgId())

	if err := tclient.Write(su); err != nil {
		ERROR.Println(err)
	} else {
		INFO.Println("SUBACK sent")
	}
}

func (tg *TransGate) handle_SUBACK(m *SubackMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_UNSUBSCRIBE(m *UnsubscribeMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_UNSUBACK(m *UnsubackMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_PINGREQ(m *PingreqMessage, c uConn, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	tclient := tg.clients.GetClient(r).(*TransClient)

	resp := NewPingResp()
	if err := tclient.Write(resp); err != nil {
		ERROR.Println(err)
	} else {
		INFO.Println("PINGRESP sent")
	}
}

func (tg *TransGate) handle_PINGRESP(m *PingrespMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_DISCONNECT(m *DisconnectMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	tclient := tg.clients.GetClient(r).(*TransClient)
	tclient.disconnectMqtt()
	tg.clients.RemoveClient(tclient.ClientId)
}

func (tg *TransGate) handle_WILLTOPICUPD(m *WillTopicUpdateMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLTOPICRESP(m *WillTopicRespMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLMSGUPD(m *WillMsgUpdateMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLMSGRESP(m *WillMsgRespMessage, r uAddr) {
	INFO.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}
