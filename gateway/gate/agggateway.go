package gateway

import (
	"fmt"
	"os"
	"sync"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"

	. "github.com/alsm/gnatt/common/protocol"
	"github.com/alsm/gnatt/common/utils"
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
	opts.SetTraceLevel(MQTT.Warn)
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
		nil,
	}

	ag.handler = func(msg MQTT.Message) {
		ag.distribute(msg)
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

func (ag *AggGate) distribute(msg MQTT.Message) {
	topic := msg.Topic()
	fmt.Printf("AG distributing a msg for topic \"%s\"\n", topic)

	// collect a list of clients to which msg should be
	// published
	// then publish msg to those clients (async)

	if clients, e := ag.tTree.SubscribersOf(topic); e != nil {
		fmt.Println(e)
	} else {
		for _, client := range clients {
			go ag.publish(msg, client)
		}
	}
}

func (ag *AggGate) publish(msg MQTT.Message, client *Client) {
	fmt.Printf("publish to client \"%s\"... ", client.ClientId)
	dup := msg.DupFlag()
	qos := QoS(msg.QoS()) // todo: what to do for qos > 0?
	ret := msg.RetainedFlag()
	top := msg.Topic()
	pay := msg.Payload()
	topicid := ag.tIndex.getId(top)
	topicidtype := byte(0x00) // todo: pre-defined (1) and shortname (2)
	msgid := uint16(0x00)     // todo: what should this be??
	pm := NewPublishMessage(dup, ret, qos, topicidtype, topicid, msgid, pay)

	if client.Registered(topicid) {
		fmt.Printf("client \"%s\" already registered to %d, publish ahoy!\n", client, topicid)
		if err := client.Write(pm); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("published a message to \"%s\"\n", client)
		}
	} else {
		fmt.Printf("client \"%s\" is not registered to %d, must REGISTER first\n", client, topicid)
		rm := NewRegisterMessage(topicid, msgid, top)
		client.AddPendingMessage(pm)
		if err := client.Write(rm); err != nil {
			fmt.Printf("error writing REGISTER to \"%s\"\n", client)
		} else {
			fmt.Printf("sent REGISTER to \"%s\" for %d (%d bytes)\n", client, topicid, rm.Length())
		}
	}
}

func (ag *AggGate) OnPacket(nbytes int, buffer []byte, con uConn, addr uAddr) {
	fmt.Println("OnPacket!")
	fmt.Printf("bytes: %s\n", utils.Bytes2str(buffer[0:nbytes]))

	rawmsg := Unpack(buffer[0:nbytes])
	fmt.Printf("rawmsg.MsgType(): %s\n", rawmsg.MsgType())

	switch msg := rawmsg.(type) {
	case *AdvertiseMessage:
		ag.handle_ADVERTISE(msg, addr)
	case *SearchGwMessage:
		ag.handle_SEARCHGW(msg, addr)
	case *GwInfoMessage:
		ag.handle_GWINFO(msg, addr)
	case *ConnectMessage:
		ag.handle_CONNECT(msg, con, addr)
	case *ConnackMessage:
		ag.handle_CONNACK(msg, addr)
	case *WillTopicReqMessage:
		ag.handle_WILLTOPICREQ(msg, addr)
	case *WillTopicMessage:
		ag.handle_WILLTOPIC(msg, addr)
	case *WillMsgReqMessage:
		ag.handle_WILLMSGREQ(msg, addr)
	case *WillMsgMessage:
		ag.handle_WILLMSG(msg, addr)
	case *RegisterMessage:
		ag.handle_REGISTER(msg, con, addr)
	case *RegackMessage:
		ag.handle_REGACK(msg, addr)
	case *PublishMessage:
		ag.handle_PUBLISH(msg, addr)
	case *PubackMessage:
		ag.handle_PUBACK(msg, addr)
	case *PubcompMessage:
		ag.handle_PUBCOMP(msg, addr)
	case *PubrecMessage:
		ag.handle_PUBREC(msg, addr)
	case *PubrelMessage:
		ag.handle_PUBREL(msg, addr)
	case *SubscribeMessage:
		ag.handle_SUBSCRIBE(msg, con, addr)
	case *SubackMessage:
		ag.handle_SUBACK(msg, addr)
	case *UnsubscribeMessage:
		ag.handle_UNSUBSCRIBE(msg, addr)
	case *UnsubackMessage:
		ag.handle_UNSUBACK(msg, addr)
	case *PingreqMessage:
		ag.handle_PINGREQ(msg, con, addr)
	case *PingrespMessage:
		ag.handle_PINGRESP(msg, addr)
	case *DisconnectMessage:
		ag.handle_DISCONNECT(msg, addr)
	case *WillTopicUpdateMessage:
		ag.handle_WILLTOPICUPD(msg, addr)
	case *WillTopicRespMessage:
		ag.handle_WILLTOPICRESP(msg, addr)
	case *WillMsgUpdateMessage:
		ag.handle_WILLMSGUPD(msg, addr)
	case *WillMsgRespMessage:
		ag.handle_WILLMSGRESP(msg, addr)
	default:
		fmt.Printf("Unknown Message Type %T\n", msg)
	}
}

func (ag *AggGate) handle_ADVERTISE(m *AdvertiseMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_SEARCHGW(m *SearchGwMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_GWINFO(m *GwInfoMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_CONNECT(m *ConnectMessage, c uConn, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)

	if clientid, e := validateClientId(m.ClientId()); e != nil {
		fmt.Println(e)
	} else {
		fmt.Printf("clientid: %s\n", clientid)
		fmt.Printf("remoteaddr: %s\n", r.r)
		fmt.Printf("will: %v\n", m.Will())

		if m.Will() {
			// todo: do something about that
		}

		client := NewClient(clientid, c, r)
		ag.clients.AddClient(client)

		ca := NewConnackMessage(0) // todo: 0 ?
		if ioerr := client.Write(ca); ioerr != nil {
			fmt.Println(ioerr)
		} else {
			fmt.Println("CONNACK was sent")
		}
	}
}

func (ag *AggGate) handle_CONNACK(m *ConnackMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLTOPICREQ(m *WillTopicReqMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLTOPIC(m *WillTopicMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLMSGREQ(m *WillMsgReqMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLMSG(m *WillMsgMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_REGISTER(m *RegisterMessage, c uConn, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	topic := string(m.TopicName())
	fmt.Printf("msg id: %d\n", m.MsgId())
	fmt.Printf("topic name: %s\n", topic)

	var topicid uint16
	if !ag.tIndex.containsTopic(topic) {
		topicid = ag.tIndex.putTopic(topic)
	} else {
		topicid = ag.tIndex.getId(topic)
	}

	client := ag.clients.GetClient(r)
	client.Register(topicid)

	fmt.Printf("ag topicid: %d\n", topicid)

	ra := NewRegackMessage(topicid, m.MsgId(), 0)
	fmt.Printf("ra.MsgId: %d\n", ra.MsgId())

	if err := client.Write(ra); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("REGACK sent")
	}
}

func (ag *AggGate) handle_REGACK(m *RegackMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	// the gateway sends a register when there is a message
	// that needs to be published, so we do that now
	topicid := m.TopicId()
	client := ag.clients.GetClient(r)
	pm := client.FetchPendingMessage(topicid)
	if pm == nil {
		fmt.Printf("no pending message for %s id %d\n", client, topicid)
	} else {
		if err := client.Write(pm); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("published a pending message to \"%s\"\n", client)
		}
	}
}

func (ag *AggGate) handle_PUBLISH(m *PublishMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)

	fmt.Printf("m.TopicId: %d\n", m.TopicId())
	fmt.Printf("m.Data: %s\n", string(m.Data()))

	topic := ag.tIndex.getTopic(m.TopicId())

	// TODO: what should the MQTT-QoS be set as? In case of MQTTSN-QoS -1 ?
	receipt := ag.mqttclient.Publish(MQTT.QoS(2), topic, m.Data())
	fmt.Println("published, waiting for receipt")
	<-receipt
	fmt.Println("receipt received")
}

func (ag *AggGate) handle_PUBACK(m *PubackMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_PUBCOMP(m *PubcompMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_PUBREC(m *PubrecMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_PUBREL(m *PubrelMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_SUBSCRIBE(m *SubscribeMessage, c uConn, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	fmt.Printf("m.TopicIdType: %d\n", m.TopicIdType())
	topic := string(m.TopicName())
	var topicid uint16
	if m.TopicIdType() == 0 {
		fmt.Printf("m.TopicName: %s\n", topic)
		if !ContainsWildcard(topic) {
			topicid = ag.tIndex.getId(topic)
			if topicid == 0 {
				topicid = ag.tIndex.putTopic(topic)
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
		client.Register(topicid)
		suba := NewSubackMessage(0, m.QoS(), topicid, m.MsgId())
		if nbytes, err := c.c.WriteToUDP(suba.Pack(), r.r); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("SUBACK sent %d bytes\n", nbytes)
		}
	}
}

func (ag *AggGate) handle_SUBACK(m *SubackMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_UNSUBSCRIBE(m *UnsubscribeMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_UNSUBACK(m *UnsubackMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_PINGREQ(m *PingreqMessage, c uConn, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	resp := NewPingResp()

	if nbytes, err := c.c.WriteToUDP(resp.Pack(), r.r); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("PINGRESP sent %d bytes\n", nbytes)
	}
}

func (ag *AggGate) handle_PINGRESP(m *PingrespMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_DISCONNECT(m *DisconnectMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	fmt.Printf("duration: %d\n", m.Duration())
}

func (ag *AggGate) handle_WILLTOPICUPD(m *WillTopicUpdateMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLTOPICRESP(m *WillTopicRespMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLMSGUPD(m *WillMsgUpdateMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (ag *AggGate) handle_WILLMSGRESP(m *WillMsgRespMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}
