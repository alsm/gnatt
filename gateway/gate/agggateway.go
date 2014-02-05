package gateway

import (
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	. "github.com/alsm/gnatt/common/protocol"
	"github.com/alsm/gnatt/common/utils"
	"os"
	"sync"
	"time"
)

type AggGate struct {
	mqttclient *MQTT.MqttClient
	stopsig    chan os.Signal
	port       int
	tIndex     topicNames
	tTree      *TopicTree
	clients    Clients
	handler    MQTT.MessageHandler
}

func NewAggGate(gc *GatewayConfig, stopsig chan os.Signal) *AggGate {
	opts := MQTT.NewClientOptions()
	opts.SetBroker(gc.mqttbroker)
	if gc.mqttuser != "" {
		opts.SetUsername(gc.mqttuser)
	}
	if gc.mqttpassword != "" {
		opts.SetPassword(gc.mqttpassword)
	}
	if gc.mqttclientid != "" {
		opts.SetClientId(gc.mqttclientid)
	}
	if gc.mqtttimeout > 0 {
		opts.SetTimeout(uint(gc.mqtttimeout))
	}
	client := MQTT.NewClient(opts)
	ag := &AggGate{
		client,
		stopsig,
		gc.port,
		topicNames{
			sync.RWMutex{},
			make(map[uint16]string),
			0,
		},
		NewTopicTree(),
		Clients{
			sync.RWMutex{},
			make(map[string]*Client),
		},
		func(msg MQTT.Message) {
			fmt.Println("subscribeHandler!")
		},
	}
	return ag
}

func (ag *AggGate) Port() int {
	return ag.port
}

func (ag *AggGate) Start() {
	go ag.awaitStop()
	fmt.Println("Aggregating Gateway is starting")
	_, err := ag.mqttclient.Start()
	if err != nil {
		fmt.Println("Aggregating Gateway failed to start")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Aggregating Gateway is started")
	listen(ag)
}

// This does NOT WORK on Windows using Cygwin, however
// it does work using cmd.exe
func (ag *AggGate) awaitStop() {
	<-ag.stopsig
	fmt.Println("Aggregating Gateway is stopping")
	ag.mqttclient.Disconnect(500)
	time.Sleep(500) //give broker some time to process DISCONNECT
	fmt.Println("Aggregating Gateway is stopped")

	// TODO: cleanly close down other goroutines

	os.Exit(0)
}

func (ag *AggGate) OnPacket(nbytes int, buffer []byte, conn uConn, remote uAddr) {
	fmt.Println("OnPacket!")
	fmt.Printf("bytes: %s\n", utils.Bytes2str(buffer[0:nbytes]))

	m := Unpack(buffer[0:nbytes])
	fmt.Printf("m.MsgType(): %s\n", m.MsgType())

	switch m.MsgType() {
	case ADVERTISE:
		ag.handle_ADVERTISE(m, remote)
	case SEARCHGW:
		ag.handle_SEARCHGW(m, remote)
	case GWINFO:
		ag.handle_GWINFO(m, remote)
	case CONNECT:
		ag.handle_CONNECT(m, conn, remote)
	case CONNACK:
		ag.handle_CONNACK(m, remote)
	case WILLTOPICREQ:
		ag.handle_WILLTOPICREQ(m, remote)
	case WILLTOPIC:
		ag.handle_WILLTOPIC(m, remote)
	case WILLMSGREQ:
		ag.handle_WILLMSGREQ(m, remote)
	case WILLMSG:
		ag.handle_WILLMSG(m, remote)
	case REGISTER:
		ag.handle_REGISTER(m, conn, remote)
	case REGACK:
		ag.handle_REGACK(m, remote)
	case PUBLISH:
		ag.handle_PUBLISH(m, remote)
	case PUBACK:
		ag.handle_PUBACK(m, remote)
	case PUBCOMP:
		ag.handle_PUBCOMP(m, remote)
	case PUBREC:
		ag.handle_PUBREC(m, remote)
	case PUBREL:
		ag.handle_PUBREL(m, remote)
	case SUBSCRIBE:
		ag.handle_SUBSCRIBE(m, conn, remote)
	case SUBACK:
		ag.handle_SUBACK(m, remote)
	case UNSUBSCRIBE:
		ag.handle_UNSUBSCRIBE(m, remote)
	case UNSUBACK:
		ag.handle_UNSUBACK(m, remote)
	case PINGREQ:
		ag.handle_PINGREQ(m, remote)
	case PINGRESP:
		ag.handle_PINGRESP(m, remote)
	case DISCONNECT:
		ag.handle_DISCONNECT(m, remote)
	case WILLTOPICUPD:
		ag.handle_WILLTOPICUPD(m, remote)
	case WILLTOPICRESP:
		ag.handle_WILLTOPICRESP(m, remote)
	case WILLMSGUPD:
		ag.handle_WILLMSGUPD(m, remote)
	case WILLMSGRESP:
		ag.handle_WILLMSGRESP(m, remote)
	}
}

func (ag *AggGate) handle_ADVERTISE(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_SEARCHGW(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_GWINFO(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_CONNECT(m Message, c uConn, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	cm, _ := m.(*ConnectMessage)
	clientid := string(cm.ClientId())
	fmt.Printf("clientid: %s\n", clientid)
	fmt.Printf("remoteaddr: %s\n", r.r)
	fmt.Printf("will: %v\n", cm.Will())

	if cm.Will() {
		// todo: do something about that
	}

	client := NewClient(string(cm.ClientId()), c, r)
	ag.clients.AddClient(client)

	ca := NewConnackMessage(0)

	if err := client.Write(ca); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("CONNACK was sent")
	}
}

func (ag *AggGate) handle_CONNACK(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLTOPICREQ(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLTOPIC(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLMSGREQ(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLMSG(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_REGISTER(m Message, c uConn, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	rm := m.(*RegisterMessage)
	topic := string(rm.TopicName())
	fmt.Printf("msg id: %d\n", rm.MsgId())
	fmt.Printf("topic name: %s\n", topic)

	var topicid uint16
	if !ag.tIndex.containsTopic(topic) {
		topicid = ag.tIndex.putTopic(topic)
	} else {
		topicid = ag.tIndex.getId(topic)
	}

	fmt.Printf("ag topicid: %d\n", topicid)

	ra := NewRegackMessage(topicid, rm.MsgId(), 0)
	fmt.Printf("ra.MsgId: %d\n", ra.MsgId())

	if nbytes, err := c.c.WriteToUDP(ra.Pack(), r.r); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("REGACK sent %d bytes\n", nbytes)
	}
}

func (ag *AggGate) handle_REGACK(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_PUBLISH(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	pm := m.(*PublishMessage)

	fmt.Printf("pm.TopicId: %d\n", pm.TopicId())
	fmt.Printf("pm.Data: %s\n", string(pm.Data()))

	topic := ag.tIndex.getTopic(pm.TopicId())

	// TODO: what should the MQTT-QoS be set as? In case of MQTTSN-QoS -1 ?
	receipt := ag.mqttclient.Publish(MQTT.QoS(2), topic, pm.Data())
	fmt.Println("published, waiting for receipt")
	<-receipt
	fmt.Println("receipt received")
}

func (ag *AggGate) handle_PUBACK(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_PUBCOMP(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_PUBREC(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_PUBREL(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_SUBSCRIBE(m Message, c uConn, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	sm := m.(*SubscribeMessage)
	fmt.Printf("sm.TopicIdType: %d\n", sm.TopicIdType())
	topic := string(sm.TopicName())
	var topid uint16
	if sm.TopicIdType() == 0 {
		fmt.Printf("sm.TopicName: %s\n", topic)
		if !ContainsWildcard(topic) {
			topid = ag.tIndex.getId(topic)
			if topid == 0 {
				topid = ag.tIndex.putTopic(topic)
			}
		} else {
			// todo: if topic contains wildcard, something about REGISTER
			// at a later time, but send topic id 0x0000 for now
		}
	} // todo: other topic id types

	client := ag.clients.GetClient(r)
	if first, err := ag.tTree.AddSubscription(client, topic); err != nil {
		fmt.Println("error adding subscription: %v\n", err)
		// todo: suback an error message?
	} else {
		if first {
			fmt.Println("first subscriber of subscription, subscribbing via MQTT")
			if receipt, sserr := ag.mqttclient.StartSubscription(ag.handler, topic, MQTT.QOS_TWO); sserr != nil {
				fmt.Printf("StartSubscription error: %v\n", sserr)
			} else {
				<-receipt
			}
		}
		// AG is subscribed at this point
		suba := NewSubackMessage(0, sm.QoS(), topid, sm.MsgId())
		if nbytes, err := c.c.WriteToUDP(suba.Pack(), r.r); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("SUBACK sent %d bytes\n", nbytes)
		}
	}
}

func (ag *AggGate) handle_SUBACK(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_UNSUBSCRIBE(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_UNSUBACK(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_PINGREQ(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_PINGRESP(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_DISCONNECT(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	dm := m.(*DisconnectMessage)
	fmt.Printf("duration: %d\n", dm.Duration())
}

func (ag *AggGate) handle_WILLTOPICUPD(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLTOPICRESP(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLMSGUPD(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLMSGRESP(m Message, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}
