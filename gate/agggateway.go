package gateway

import (
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	. "github.com/alsm/gnatt/common/protocol"
	"github.com/alsm/gnatt/common/utils"
	"net"
	"os"
	"sync"
	"time"
)

type AggGate struct {
	mqttclient *MQTT.MqttClient
	stopsig    chan os.Signal
	port       int
	topics     topicNames
}

func NewAggGate(opts *MQTT.ClientOptions, stopsig chan os.Signal, port int) *AggGate {
	client := MQTT.NewClient(opts)
	ag := &AggGate{
		client,
		stopsig,
		port,
		topicNames{
			sync.RWMutex{},
			make(map[string]uint16),
			0,
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

func (ag *AggGate) OnPacket(nbytes int, buffer []byte, conn *net.UDPConn, remote *net.UDPAddr) {
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
		ag.handle_SUBSCRIBE(m, remote)
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

func (ag *AggGate) handle_ADVERTISE(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_SEARCHGW(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_GWINFO(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_CONNECT(m Message, c *net.UDPConn, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
	wm, _ := m.(*ConnectMessage)
	fmt.Printf("will: %v\n", wm.Will())

	if wm.Will() {
		// do something about that
	}

	ca := NewConnackMessage(0)

	if nbytes, err := c.WriteToUDP(ca.Pack(), r); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("CONNACK sent %d bytes\n", nbytes)
	}
}

func (ag *AggGate) handle_CONNACK(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_WILLTOPICREQ(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_WILLTOPIC(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_WILLMSGREQ(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_WILLMSG(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_REGISTER(m Message, c *net.UDPConn, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
	rm := m.(*RegisterMessage)
	topic := string(rm.TopicName())
	fmt.Printf("msg id: %d\n", rm.MsgId())
	fmt.Printf("topic name: %s\n", topic)

	var topicid uint16
	if !ag.topics.contains(topic) {
		topicid = ag.topics.put(topic)
	} else {
		topicid = ag.topics.get(topic)
	}

	fmt.Printf("ag topicid: %d\n", topicid)

	ra := NewRegackMessage(topicid, rm.MsgId(), 0)
	fmt.Printf("ra.MsgId: %d\n", ra.MsgId())

	if nbytes, err := c.WriteToUDP(ra.Pack(), r); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("REGACK sent %d bytes\n", nbytes)
	}
}

func (ag *AggGate) handle_REGACK(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_PUBLISH(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_PUBACK(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_PUBCOMP(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_PUBREC(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_PUBREL(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_SUBSCRIBE(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_SUBACK(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_UNSUBSCRIBE(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_UNSUBACK(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_PINGREQ(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_PINGRESP(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_DISCONNECT(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
	dm := m.(*DisconnectMessage)
	fmt.Printf("duration: %d\n", dm.Duration())
}

func (ag *AggGate) handle_WILLTOPICUPD(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_WILLTOPICRESP(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_WILLMSGUPD(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}

func (ag *AggGate) handle_WILLMSGRESP(m Message, r *net.UDPAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r)
}
