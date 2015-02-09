package gateway

import (
	"bytes"
	"log"
	"net"
	"os"
	"sync"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"

	. "github.com/alsm/gnatt/packets"
)

type AGateway struct {
	mqttclient *MQTT.Client
	stopsig    chan os.Signal
	port       int
	tIndex     topicNames
	tTree      *TopicTree
	clients    Clients
	handler    MQTT.MessageHandler
}

func NewAGateway(gc *GatewayConfig, stopsig chan os.Signal) *AGateway {
	MQTT.WARN = log.New(os.Stdout, "", 0)
	MQTT.DEBUG = log.New(os.Stdout, "", 0)
	MQTT.CRITICAL = log.New(os.Stdout, "", 0)
	MQTT.ERROR = log.New(os.Stdout, "", 0)
	opts := MQTT.NewClientOptions()
	opts.AddBroker(gc.mqttbroker)
	if gc.mqttuser != "" {
		opts.SetUsername(gc.mqttuser)
	}
	if gc.mqttpassword != "" {
		opts.SetPassword(gc.mqttpassword)
	}
	if gc.mqttclientid != "" {
		opts.SetClientID(gc.mqttclientid)
	}
	if gc.mqtttimeout > 0 {
		opts.SetKeepAlive(time.Duration(gc.mqtttimeout))
	}
	client := MQTT.NewClient(opts)
	ag := &AGateway{
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
			make(map[string]SNClient),
		},
		nil,
	}

	ag.handler = func(client *MQTT.Client, msg MQTT.Message) {
		ag.distribute(msg)
	}

	return ag
}

func (ag *AGateway) Port() int {
	return ag.port
}

func (ag *AGateway) Start() {
	go ag.awaitStop()
	INFO.Println("Aggregating Gateway is starting")
	if token := ag.mqttclient.Connect(); token.Wait() && token.Error() != nil {
		ERROR.Println(token.Error())
		return
	}
	INFO.Println("Aggregating Gateway is started")
	listen(ag)
}

// This does NOT WORK on Windows using Cygwin, however
// it does work using cmd.exe
func (ag *AGateway) awaitStop() {
	<-ag.stopsig
	INFO.Println("Aggregating Gateway is stopping")
	ag.mqttclient.Disconnect(500)
	time.Sleep(500) //give broker some time to process DISCONNECT
	INFO.Println("Aggregating Gateway is stopped")

	// TODO: cleanly close down other goroutines

	os.Exit(0)
}

func (ag *AGateway) distribute(msg MQTT.Message) {
	topic := msg.Topic()
	INFO.Printf("AG distributing a msg for topic \"%s\"\n", topic)

	// collect a list of clients to which msg should be
	// published
	// then publish msg to those clients (async)

	if clients, e := ag.tTree.SubscribersOf(topic); e != nil {
		ERROR.Println(e)
	} else {
		for _, client := range clients {
			go ag.publish(msg, client)
		}
	}
}

func (ag *AGateway) publish(msg MQTT.Message, client *Client) {
	INFO.Printf("publish to client \"%s\"... ", client.ClientId)
	topicid := ag.tIndex.getId(msg.Topic())
	// topicidtype := byte(0x00) // todo: pre-defined (1) and shortname (2)
	// msgid := uint16(0x00) // todo: what should this be??
	pm := NewPublishMessage(topicid, 0x00, msg.Payload(), msg.Qos(), 0x00, msg.Retained(), msg.Duplicate())

	if client.Registered(topicid) {
		INFO.Printf("client \"%s\" already registered to %d, publish ahoy!\n", client, topicid)
		if err := client.Write(pm); err != nil {
			ERROR.Println(err)
		} else {
			INFO.Printf("published a message to \"%s\"\n", client)
		}
	} else {
		INFO.Printf("client \"%s\" is not registered to %d, must REGISTER first\n", client, topicid)
		rm := NewRegisterMessage(topicid, 0x00, []byte(msg.Topic()))
		client.AddPendingMessage(pm)
		if err := client.Write(rm); err != nil {
			ERROR.Printf("error writing REGISTER to \"%s\"\n", client)
		} else {
			INFO.Printf("sent REGISTER to \"%s\" for %d (%d bytes)\n", client, topicid, rm.Length)
		}
	}
}

func (ag *AGateway) OnPacket(nbytes int, buffer []byte, con *net.UDPConn, addr *net.UDPAddr) {
	INFO.Printf("OnPacket!  - bytes: %s\n", string(buffer[0:nbytes]))

	buf := bytes.NewBuffer(buffer)
	rawmsg, _ := ReadPacket(buf)
	INFO.Printf("rawmsg.MessageType(): %s\n", rawmsg.MessageType())

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
		ERROR.Printf("Unknown Message Type %T\n", msg)
	}
}

func (ag *AGateway) handle_ADVERTISE(m *AdvertiseMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_SEARCHGW(m *SearchGwMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_GWINFO(m *GwInfoMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_CONNECT(m *ConnectMessage, c *net.UDPConn, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)

	if clientid, e := validateClientId(m.ClientId); e != nil {
		ERROR.Println(e)
	} else {
		INFO.Printf("clientid: %s\n", clientid)
		INFO.Printf("remoteaddr: %s\n", r)
		INFO.Printf("will: %v\n", m.Will)

		if m.Will {
			// todo: do something about that
		}

		client := NewClient(clientid, c, r)
		ag.clients.AddClient(client)

		ca := NewMessage(CONNACK).(*ConnackMessage) // todo: 0 ?
		ca.ReturnCode = 0
		if ioerr := client.Write(ca); ioerr != nil {
			ERROR.Println(ioerr)
		} else {
			INFO.Println("CONNACK was sent")
		}
	}
}

func (ag *AGateway) handle_CONNACK(m *ConnackMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_WILLTOPICREQ(m *WillTopicReqMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_WILLTOPIC(m *WillTopicMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_WILLMSGREQ(m *WillMsgReqMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_WILLMSG(m *WillMsgMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_REGISTER(m *RegisterMessage, c *net.UDPConn, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
	topic := string(m.TopicName)
	INFO.Printf("msg id: %d\n", m.MessageId)
	INFO.Printf("topic name: %s\n", topic)

	var topicid uint16
	if !ag.tIndex.containsTopic(topic) {
		topicid = ag.tIndex.putTopic(topic)
	} else {
		topicid = ag.tIndex.getId(topic)
	}

	client := ag.clients.GetClient(r).(*Client)
	client.Register(topicid, topic)

	INFO.Printf("ag topicid: %d\n", topicid)

	ra := NewRegackMessage(topicid, m.MessageId, 0)
	INFO.Printf("ra.MsgId: %d\n", ra.MessageId)

	if err := client.Write(ra); err != nil {
		ERROR.Println(err)
	} else {
		INFO.Println("REGACK sent")
	}
}

func (ag *AGateway) handle_REGACK(m *RegackMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
	// the gateway sends a register when there is a message
	// that needs to be published, so we do that now
	topicid := m.TopicId
	client := ag.clients.GetClient(r).(*Client)
	pm := client.FetchPendingMessage(topicid)
	if pm == nil {
		ERROR.Printf("no pending message for %s id %d\n", client, topicid)
	} else {
		if err := client.Write(pm); err != nil {
			ERROR.Println(err)
		} else {
			INFO.Printf("published a pending message to \"%s\"\n", client)
		}
	}
}

func (ag *AGateway) handle_PUBLISH(m *PublishMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)

	INFO.Printf("m.TopicId: %d\n", m.TopicId)
	INFO.Printf("m.Data: %s\n", string(m.Data))

	topic := ag.tIndex.getTopic(m.TopicId)

	// TODO: what should the MQTT-QoS be set as? In case of MQTTSN-QoS -1 ?
	if token := ag.mqttclient.Publish(topic, m.Qos, m.Retain, m.Data); token.WaitTimeout(2000) && token.Error() != nil {
		ERROR.Println("Error publishing message", token.Error())
	}
	INFO.Println("Message Published")
}

func (ag *AGateway) handle_PUBACK(m *PubackMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_PUBCOMP(m *PubcompMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_PUBREC(m *PubrecMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_PUBREL(m *PubrelMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_SUBSCRIBE(m *SubscribeMessage, c *net.UDPConn, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
	INFO.Printf("m.TopicIdType: %d\n", m.TopicIdType)
	topic := string(m.TopicName)
	var topicid uint16
	if m.TopicIdType == 0 {
		INFO.Printf("m.TopicName: %s\n", topic)
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

	client := ag.clients.GetClient(r).(*Client)
	if first, err := ag.tTree.AddSubscription(client, topic); err != nil {
		INFO.Println("error adding subscription: %v\n", err)
		// todo: suback an error message?
	} else {
		if first {
			INFO.Println("first subscriber of subscription, subscribbing via MQTT")
			if token := ag.mqttclient.Subscribe(topic, 2, ag.handler); token.WaitTimeout(2000) && token.Error() != nil {
				ERROR.Println("Error subscribing,", token.Error())
			}
		}
		// AG is subscribed at this point
		client.Register(topicid, topic)
		suba := NewSubackMessage(topicid, m.MessageId, m.Qos, 0)
		var buf bytes.Buffer
		suba.Write(&buf)
		if nbytes, err := c.WriteToUDP(buf.Bytes(), r); err != nil {
			ERROR.Println(err)
		} else {
			INFO.Printf("SUBACK sent %d bytes\n", nbytes)
		}
	}
}

func (ag *AGateway) handle_SUBACK(m *SubackMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_UNSUBSCRIBE(m *UnsubscribeMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_UNSUBACK(m *UnsubackMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_PINGREQ(m *PingreqMessage, c *net.UDPConn, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
	resp := NewMessage(PINGRESP)

	var buf bytes.Buffer
	resp.Write(&buf)

	if nbytes, err := c.WriteToUDP(buf.Bytes(), r); err != nil {
		ERROR.Println(err)
	} else {
		INFO.Printf("PINGRESP sent %d bytes\n", nbytes)
	}
}

func (ag *AGateway) handle_PINGRESP(m *PingrespMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_DISCONNECT(m *DisconnectMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
	INFO.Printf("duration: %d\n", m.Duration)
	// todo: cleanup the client
}

func (ag *AGateway) handle_WILLTOPICUPD(m *WillTopicUpdateMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_WILLTOPICRESP(m *WillTopicRespMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_WILLMSGUPD(m *WillMsgUpdateMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}

func (ag *AGateway) handle_WILLMSGRESP(m *WillMsgRespMessage, r *net.UDPAddr) {
	INFO.Printf("handle_%s from %v\n", m.MessageType(), r)
}
