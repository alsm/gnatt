package gateway

import (
	"fmt"
	"os"
	"sync"

	. "github.com/alsm/gnatt/common/protocol"
	"github.com/alsm/gnatt/common/utils"
)

type TransGate struct {
	stopsig    chan os.Signal
	port       int
	mqttbroker string
	clients    Clients
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
	}
	return tg
}

func (tg *TransGate) Port() int {
	return tg.port
}

func (tg *TransGate) Start() {
	go tg.awaitStop()
	fmt.Println("Transparent Gateway is starting")
	fmt.Println("Transparent Gataway is started")
	listen(tg)
}

func (tg *TransGate) awaitStop() {
	<-tg.stopsig
	fmt.Println("Transparent Gateway is stopping")
	fmt.Println("Transparent Gateway is stopped")
	os.Exit(0)
}

func (tg *TransGate) OnPacket(nbytes int, buffer []byte, con uConn, addr uAddr) {
	fmt.Println("TG OnPacket!")
	fmt.Printf("bytes: %s\n", utils.Bytes2str(buffer[0:nbytes]))
	rawmsg := Unpack(buffer[0:nbytes])
	fmt.Printf("rawmsg.MsgType(): %s\n", rawmsg.MsgType())

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
		tg.handle_SUBSCRIBE(msg, con, addr)
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
		fmt.Printf("Unknown Message Type %T\n", msg)
	}
}

func (tg *TransGate) handle_ADVERTISE(m *AdvertiseMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_SEARCHGW(m *SearchGwMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_GWINFO(m *GwInfoMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_CONNECT(m *ConnectMessage, c uConn, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
	if clientid, err := validateClientId(m.ClientId()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("clientid: %s\n", clientid)
		fmt.Printf("remoteaddr: %s\n", r.r)
		fmt.Printf("will: %v\n", m.Will())
		if m.Will() {
			// todo: will msg
		}
		tclient := NewTransClient(string(clientid), c, r)
		tg.clients.AddClient(tclient)

		ca := NewConnackMessage(0) // todo: 0 ?
		if ioerr := tclient.Write(ca); ioerr != nil {
			fmt.Println(ioerr)
		} else {
			fmt.Println("CONNACK was sent")
		}
	}
}

func (tg *TransGate) handle_CONNACK(m *ConnackMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLTOPICREQ(m *WillTopicReqMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLTOPIC(m *WillTopicMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLMSGREQ(m *WillMsgReqMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLMSG(m *WillMsgMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_REGISTER(m *RegisterMessage, c uConn, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)

}

func (tg *TransGate) handle_REGACK(m *RegackMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_PUBLISH(m *PublishMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_PUBACK(m *PubackMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_PUBCOMP(m *PubcompMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_PUBREC(m *PubrecMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_PUBREL(m *PubrelMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_SUBSCRIBE(m *SubscribeMessage, c uConn, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_SUBACK(m *SubackMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_UNSUBSCRIBE(m *UnsubscribeMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_UNSUBACK(m *UnsubackMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_PINGREQ(m *PingreqMessage, c uConn, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_PINGRESP(m *PingrespMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_DISCONNECT(m *DisconnectMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLTOPICUPD(m *WillTopicUpdateMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLTOPICRESP(m *WillTopicRespMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLMSGUPD(m *WillMsgUpdateMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}

func (tg *TransGate) handle_WILLMSGRESP(m *WillMsgRespMessage, r uAddr) {
	fmt.Printf("handle_%s from %v\n", m.MsgType(), r.r)
}
