package gnatt

import (
	"fmt"
	. "github.com/alsm/gnatt/common/utils"
)

type MsgType byte

type Header struct {
	length  []byte
	msgType MsgType
}

func (h *Header) SetLength(length int) {
	var tLength []byte
	if length < 256 && length > 1 {
		tLength = append(tLength, byte(length))
	} else {
		tLength = append(tLength, 0x01)
		tLength = append(tLength, U162b(uint16(length))...)
	}
	h.length = tLength
}

func (h *Header) Length() (length int) {
	if h.length[0] == 0x01 {
		length = int(B2u16(h.length[1:3]))
	} else {
		length = int(h.length[0])
	}
	return
}

func (h Header) MsgType() (m MsgType) {
	m = h.msgType
	return
}

func (h *Header) SetMsgType(m MsgType) {
	h.msgType = m
}

func (h *Header) PackHeader() (packed []byte) {
	packed = append(packed, h.length...)
	packed = append(packed, byte(h.msgType))
	return
}

func (h *Header) UnpackHeader(msg []byte) []byte {
	if msg[0] == 0x01 {
		h.length = msg[0:3]
		h.msgType = MsgType(msg[3])
		return msg[4:]
	} else {
		h.length = msg[0:1]
		h.msgType = MsgType(msg[1])
		return msg[2:]
	}
}

func Unpack(packet []byte) Message {
	fmt.Printf("Unpack: %s\n", Bytes2str(packet))
	var h Header
	h.UnpackHeader(packet)
	fmt.Printf("Unpack h.MsgType() %s\n", h.MsgType())
	m := NewMessage(h.MsgType())
	m.Unpack(packet)
	fmt.Printf("Unpack m.MsgType() %s\n", m.MsgType())
	return m
}

type Message interface {
	MsgType() MsgType
	Pack() []byte
	Unpack([]byte) Message
}

func NewMessage(msgType MsgType) (m Message) {
	switch msgType {
	case ADVERTISE:
		m = new(AdvertiseMessage)
	case SEARCHGW:
		m = new(searchgwMessage)
	case GWINFO:
		m = new(gwInfoMessage)
	case CONNECT:
		m = new(ConnectMessage)
	case CONNACK:
		m = new(ConnackMessage)
	case WILLTOPICREQ:
		m = new(willTopicReqMessage)
	case WILLTOPIC:
		m = new(willTopicMessage)
	case WILLMSGREQ:
		m = new(willMsgReqMessage)
	case WILLMSG:
		m = new(willMsgMessage)
	case REGISTER:
		m = new(RegisterMessage)
	case REGACK:
		m = new(RegackMessage)
	case PUBLISH:
		m = new(PublishMessage)
	case PUBACK:
		m = new(pubackMessage)
	case PUBCOMP:
		m = new(pubcompMessage)
	case PUBREC:
		m = new(pubrecMessage)
	case PUBREL:
		m = new(pubrelMessage)
	case SUBSCRIBE:
		m = new(SubscribeMessage)
	case SUBACK:
		m = new(SubackMessage)
	case UNSUBSCRIBE:
		m = new(unsubscribeMessage)
	case UNSUBACK:
		m = new(unsubackMessage)
	case PINGREQ:
		m = new(pingreqMessage)
	case PINGRESP:
		m = new(pingrespMessage)
	case DISCONNECT:
		m = new(DisconnectMessage)
	case WILLTOPICUPD:
		m = new(willTopicUpdateMessage)
	case WILLTOPICRESP:
		m = new(willTopicRespMessage)
	case WILLMSGUPD:
		m = new(willMsgUpdateMessage)
	case WILLMSGRESP:
		m = new(willMsgRespMessage)
	}
	return
}

type QoS byte

const (
	/* These values are specified in MQTT-SN v1.2 Section 5.3.4 */
	QoS_NegOne QoS = 0x03
	QoS_Zero   QoS = 0x00
	QoS_One    QoS = 0x01
	QoS_Two    QoS = 0x02
)

type qoS struct {
	qos QoS
}

func (q *qoS) QoS() QoS {
	return q.qos
}

func (q *qoS) SetQoS(qos QoS) {
	switch qos {
	case QoS_NegOne:
		q.qos = QoS_NegOne
	case QoS_Zero:
		q.qos = QoS_Zero
	case QoS_One:
		q.qos = QoS_One
	case QoS_Two:
		q.qos = QoS_Two
	default:
		// User is bad at programming, better be safe
		q.qos = QoS_Two
	}
}

type gwId struct {
	gwId byte
}

func (g *gwId) GwId() byte {
	return g.gwId
}

func (g *gwId) SetGwId(gwId byte) {
	g.gwId = gwId
}

// Types are
//  0x00 - topic name
//  0x01 - predefined topic id
//  0x02 - short topic name
type topicIdType struct {
	topicIdType byte
}

func (t *topicIdType) TopicIdType() byte {
	return t.topicIdType
}

func (t topicIdType) SetTopicIdType(TopicIdType byte) {
	switch {
	case (TopicIdType < 0):
		t.topicIdType = 0
	case (TopicIdType > 2):
		t.topicIdType = 2
	default:
		t.topicIdType = TopicIdType
	}
}

type topicId struct {
	topicId uint16
}

func (t *topicId) TopicId() uint16 {
	return t.topicId
}

func (t *topicId) SetTopicId(TopicId uint16) {
	t.topicId = TopicId
}

type topicName struct {
	topicName []byte
}

func (t *topicName) TopicName() []byte {
	return t.topicName
}

func (t *topicName) SetTopicName(topicName []byte) {
	t.topicName = topicName
}

type msgId struct {
	msgId uint16
}

func (m *msgId) MsgId() uint16 {
	return m.msgId
}

func (m *msgId) SetMsgId(msgId uint16) {
	m.msgId = msgId
}

type duration struct {
	duration uint16
}

func (d *duration) Duration() uint16 {
	return d.duration
}

func (d *duration) SetDuration(duration uint16) {
	d.duration = duration
}

type dUP struct {
	dup bool
}

func (d *dUP) SetDUP(dup bool) {
	d.dup = dup
}

func (d *dUP) DUP() bool {
	return d.dup
}

type retain struct {
	retain bool
}

func (r *retain) Retain() bool {
	return r.retain
}

func (r *retain) SetRetain(retain bool) {
	r.retain = retain
}

type willTopic struct {
	willTopic []byte
}

func (w *willTopic) WillTopic() []byte {
	return w.willTopic
}

func (w *willTopic) SetWillTopic(willTopic []byte) {
	w.willTopic = willTopic
}

type willMsg struct {
	willMsg []byte
}

func (w *willMsg) WillMsg() []byte {
	return w.willMsg
}

func (w *willMsg) SetWillMsg(willMsg []byte) {
	w.willMsg = willMsg
}

type msgReturnCode struct {
	returnCode byte
}

func (r *msgReturnCode) MsgReturnCode() byte {
	return r.returnCode
}

func (r *msgReturnCode) SetMsgReturnCode(rc byte) {
	r.returnCode = rc
}

type clientId struct {
	clientId []byte
}

func (c *clientId) ClientId() []byte {
	return c.clientId
}

func (c *clientId) SetClientId(clientId []byte) {
	c.clientId = clientId
}

/*************
 * Advertise *
 *************/

type AdvertiseMessage struct {
	Header
	gwId
	duration
}

func (a *AdvertiseMessage) Pack() (packed []byte) {
	packed = append(packed, a.PackHeader()...)
	packed = append(packed, a.GwId())
	packed = append(packed, U162b(a.Duration())...)
	return
}

func (a *AdvertiseMessage) Unpack(msg []byte) Message {
	msg = a.UnpackHeader(msg)
	a.SetGwId(msg[0])
	a.SetDuration(B2u16(msg[1:3]))
	return a
}

/*************
 * Search GW *
 *************/

type searchgwMessage struct {
	Header
	radius byte
}

func (s *searchgwMessage) Radius() byte {
	return s.radius
}

func (s *searchgwMessage) SetRadius(radius byte) {
	s.radius = radius
}

func (s *searchgwMessage) Pack() (msg []byte) {
	msg = append(msg, s.PackHeader()...)
	msg = append(msg, s.Radius())
	return
}

func (s *searchgwMessage) Unpack(msg []byte) Message {
	msg = s.UnpackHeader(msg)
	s.SetRadius(msg[0])
	return Message(s)
}

/***********
 * GW Info *
 ***********/

type gwInfoMessage struct {
	Header
	gwId
	gwAddress []byte
}

func (g *gwInfoMessage) GwAddress() []byte {
	return g.gwAddress
}

func (g *gwInfoMessage) SetGwAddress(gwAddress []byte) {
	g.gwAddress = gwAddress
}

func (g *gwInfoMessage) Pack() (msg []byte) {
	msg = append(msg, g.PackHeader()...)
	msg = append(msg, g.GwId())
	msg = append(msg, g.GwAddress()...)
	return
}

func (g *gwInfoMessage) Unpack(msg []byte) Message {
	msg = g.UnpackHeader(msg)
	g.SetGwId(msg[0])
	g.SetGwAddress(msg[1:])
	return Message(g)
}

/***********
 * Connect *
 ***********/

type ConnectMessage struct {
	Header
	will         bool
	cleanSession bool
	protocolId   byte
	duration
	clientId
}

func (c *ConnectMessage) Will() bool {
	return c.will
}

func (c *ConnectMessage) SetWill(will bool) *ConnectMessage {
	c.will = will
	return c
}

func (c *ConnectMessage) CleanSession() bool {
	return c.cleanSession
}

func (c *ConnectMessage) SetCleanSession(cleanSession bool) *ConnectMessage {
	c.cleanSession = cleanSession
	return c
}

func (c *ConnectMessage) ProtocolId() byte {
	return c.protocolId
}

func (c *ConnectMessage) SetProtocolId(protocolId byte) *ConnectMessage {
	c.protocolId = protocolId
	return c
}

func (c *ConnectMessage) decodeFlags(b byte) (willed, cleansess bool) {
	if b&0x02 > 0 {
		willed = true
	}
	if b&0x01 > 0 {
		cleansess = true
	}
	return willed, cleansess
}

func (c *ConnectMessage) encodeFlags() byte {
	var b byte
	if c.Will() {
		b |= 0x02
	}
	if c.CleanSession() {
		b |= 0x01
	}
	return b
}

func (c *ConnectMessage) Pack() (bytes []byte) {
	bytes = append(bytes, c.PackHeader()...)
	bytes = append(bytes, c.encodeFlags())
	bytes = append(bytes, c.ProtocolId())
	bytes = append(bytes, U162b(c.duration.Duration())...)
	bytes = append(bytes, c.clientId.ClientId()...)
	return bytes
}

func (c *ConnectMessage) Unpack(bytes []byte) Message {
	bytes = c.UnpackHeader(bytes)
	willed, cleansess := c.decodeFlags(bytes[0])
	c.SetWill(willed)
	c.SetCleanSession(cleansess)
	c.SetProtocolId(bytes[1])
	c.SetDuration(B2u16(bytes[2:4]))
	c.SetClientId(bytes[4:])
	return c
}

/***********
 * Connack *
 ***********/

type ConnackMessage struct {
	Header
	msgReturnCode
}

func NewConnackMessage(rc byte) *ConnackMessage {
	var ca ConnackMessage
	ca.SetLength(3)
	ca.SetMsgType(CONNACK)
	ca.SetMsgReturnCode(rc)
	return &ca
}

func (c *ConnackMessage) Pack() (bytes []byte) {
	bytes = append(bytes, c.PackHeader()...)
	bytes = append(bytes, c.MsgReturnCode())
	return
}

func (c *ConnackMessage) Unpack(msg []byte) Message {
	msg = c.UnpackHeader(msg)
	c.SetMsgReturnCode(msg[0])
	return c
}

/******************
 * Will Toipc Req *
 ******************/

type willTopicReqMessage struct {
	Header
}

func (w *willTopicReqMessage) Pack() []byte {
	return w.PackHeader()
}

func (w *willTopicReqMessage) Unpack(msg []byte) Message {
	_ = w.UnpackHeader(msg)
	return w
}

/**************
 * Will Topic *
 **************/

type willTopicMessage struct {
	Header
	qoS
	willTopic
}

func (w *willTopicMessage) Pack() (msg []byte) {
	msg = append(msg, w.PackHeader()...)
	msg = append(msg, byte(w.QoS()))
	msg = append(msg, w.WillTopic()...)
	return
}

func (w *willTopicMessage) Unpack(msg []byte) Message {
	msg = w.UnpackHeader(msg)
	w.SetQoS(QoS(msg[0]))
	w.SetWillTopic(msg[1:])
	return w
}

/****************
 * Will Msg Req *
 ****************/

type willMsgReqMessage struct {
	Header
}

func (w *willMsgReqMessage) Pack() []byte {
	return w.PackHeader()
}

func (w *willMsgReqMessage) Unpack(msg []byte) Message {
	_ = w.UnpackHeader(msg)
	return w
}

/************
 * Will Msg *
 ************/

type willMsgMessage struct {
	Header
	willMsg
}

func (w *willMsgMessage) Pack() (msg []byte) {
	msg = append(msg, w.PackHeader()...)
	msg = append(msg, w.WillMsg()...)
	return
}

func (w *willMsgMessage) Unpack(msg []byte) Message {
	msg = w.UnpackHeader(msg)
	w.SetWillMsg(msg)
	return w
}

/************
 * Register *
 ************/

type RegisterMessage struct {
	Header
	topicId
	msgId
	topicName
}

func (r *RegisterMessage) Pack() (msg []byte) {
	return r.PackHeader()
}

func (r *RegisterMessage) Unpack(msg []byte) Message {
	msg = r.UnpackHeader(msg)
	r.SetTopicId(B2u16(msg[0:2]))
	r.SetMsgId(B2u16(msg[2:4]))
	r.SetTopicName(msg[4:len(msg)])
	return r
}

/**********
 * Regack *
 **********/

type RegackMessage struct {
	Header
	topicId
	msgId
	msgReturnCode
}

func NewRegackMessage(topicId, msgId uint16, rc byte) *RegackMessage {
	var ra RegackMessage
	ra.SetLength(7)
	ra.SetMsgType(REGACK)
	ra.SetTopicId(topicId)
	ra.SetMsgId(msgId)
	ra.SetMsgReturnCode(rc)
	return &ra
}

func (r *RegackMessage) Pack() (bytes []byte) {
	bytes = append(bytes, r.PackHeader()...)
	bytes = append(bytes, U162b(r.TopicId())...)
	bytes = append(bytes, U162b(r.MsgId())...)
	bytes = append(bytes, r.MsgReturnCode())
	return bytes
}

func (r *RegackMessage) Unpack(msg []byte) Message {
	return r
}

/***********
 * Publish *
 ***********/

type PublishMessage struct {
	Header
	dUP
	qoS
	topicIdType
	topicId
	msgId
	retain
	data []byte
}

func NewPublishMessage(d, r bool, q QoS, t byte, i, m uint16, z []byte) *PublishMessage {
	var pm PublishMessage
	pm.SetLength(7 + len(z))
	pm.SetMsgType(PUBLISH)
	pm.SetDUP(d)
	pm.SetQoS(QoS(q))
	pm.SetTopicIdType(t)
	pm.SetTopicId(i)
	pm.SetMsgId(m)
	pm.SetRetain(r)
	pm.SetData(z)
	return &pm
}

func (p *PublishMessage) Data() []byte {
	return p.data
}

func (p *PublishMessage) SetData(data []byte) *PublishMessage {
	p.data = data
	return p
}

func (p *PublishMessage) encodeFlags() byte {
	var b byte
	if p.dUP.DUP() {
		b = 0x80
	}
	b |= byte(p.qoS.QoS()) << 6
	if p.retain.Retain() {
		b |= (1 << 4)
	}
	b |= p.topicIdType.TopicIdType()
	return b
}

func (p *PublishMessage) decodeFlags(b byte) (d dUP, q qoS, r retain, t topicIdType) {
	d.SetDUP((b & 0x80) == 1)
	q.SetQoS(QoS(b & 0x60))
	r.SetRetain((b & 0x10) == 1)
	t.SetTopicIdType(b & 0x03)
	return d, q, r, t
}

func (p *PublishMessage) Pack() (bytes []byte) {
	bytes = append(bytes, p.PackHeader()...)
	bytes = append(bytes, p.encodeFlags())
	bytes = append(bytes, U162b(p.TopicId())...)
	bytes = append(bytes, U162b(p.MsgId())...)
	bytes = append(bytes, p.Data()...)
	return bytes
}

func (p *PublishMessage) Unpack(bytes []byte) Message {
	bytes = p.UnpackHeader(bytes)
	d, q, r, t := p.decodeFlags(bytes[0])
	p.SetDUP(d.DUP())
	p.SetQoS(q.QoS())
	p.SetRetain(r.Retain())
	p.SetTopicIdType(t.TopicIdType())

	p.SetTopicId(B2u16(bytes[1:3]))
	p.SetMsgId(B2u16(bytes[3:5]))
	p.SetData(bytes[5:])
	return p
}

/**********
 * Puback *
 **********/

type pubackMessage struct {
	Header
	topicId
	msgId
	msgReturnCode
}

func (p *pubackMessage) Pack() []byte {
	return p.PackHeader()
}

func (p *pubackMessage) Unpack(msg []byte) Message {
	return p
}

/**********
 * Pubrec *
 **********/

type pubrecMessage struct {
	Header
	msgId
}

func (p *pubrecMessage) Pack() []byte {
	return p.PackHeader()
}

func (p *pubrecMessage) Unpack(msg []byte) Message {
	return p
}

/**********
 * Pubrel *
 **********/

type pubrelMessage struct {
	Header
	msgId
}

func (p *pubrelMessage) Pack() []byte {
	return p.PackHeader()
}

func (p *pubrelMessage) Unpack(msg []byte) Message {
	return p
}

/***********
 * Pubcomp *
 ***********/

type pubcompMessage struct {
	Header
	msgId
}

func (p *pubcompMessage) Pack() []byte {
	return p.PackHeader()
}

func (p *pubcompMessage) Unpack(msg []byte) Message {
	return p
}

/*************
 * Subscribe *
 *************/

type SubscribeMessage struct {
	Header
	dUP
	qoS
	topicIdType
	topicId
	msgId
	topicName
}

func (s *SubscribeMessage) encodeFlags() byte {
	var b byte
	if s.dUP.DUP() {
		b = 0x80
	}
	b |= byte(s.qoS.QoS()) << 6
	b |= s.topicIdType.TopicIdType()
	return b
}

func (s *SubscribeMessage) decodeFlags(b byte) (d dUP, q qoS, t topicIdType) {
	d.SetDUP((b & 0x80) == 1)
	q.SetQoS(QoS(b & 0x60))
	t.SetTopicIdType(b & 0x03)
	return
}

func (s *SubscribeMessage) Pack() (bytes []byte) {
	bytes = append(bytes, s.PackHeader()...)
	bytes = append(bytes, s.encodeFlags())
	bytes = append(bytes, U162b(s.MsgId())...)
	if s.topicId.TopicId() != 0 {
		bytes = append(bytes, U162b(s.TopicId())...)
	} else {
		bytes = append(bytes, s.TopicName()...)
	}
	return
}

func (s *SubscribeMessage) Unpack(bytes []byte) Message {
	x := len(bytes)
	bytes = s.UnpackHeader(bytes)
	x = x - len(bytes)
	d, q, t := s.decodeFlags(bytes[0])
	s.SetDUP(d.DUP())
	s.SetQoS(q.QoS())
	s.SetTopicIdType(t.TopicIdType())

	if s.TopicIdType() == 0 { // TODO: make an enum
		fmt.Println("topic id type 0 (name)")
		s.SetTopicName(bytes[5-x : len(bytes)])
	} else if s.TopicIdType() == 1 {
		fmt.Println("topic id type 1 (pre-defined)")
		s.SetTopicId(B2u16(bytes[5-x : 7-x]))
	} else {
		fmt.Println("topic id type 2 (short name)")
		s.SetTopicName(bytes[5-x : 7-x])
	}
	s.SetMsgId(B2u16(bytes[3-x : 5-x]))

	return s
}

/**********
 * Suback *
 **********/

type SubackMessage struct {
	Header
	qoS
	topicId
	msgId
	msgReturnCode
}

func NewSubackMessage(rc byte, q QoS, topicId, msgId uint16) *SubackMessage {
	var sa SubackMessage
	sa.SetLength(8)
	sa.SetMsgType(SUBACK)
	sa.SetQoS(QoS(q))
	sa.SetTopicId(topicId)
	sa.SetMsgId(msgId)
	sa.SetMsgReturnCode(rc)
	return &sa
}

func (s *SubackMessage) encodeFlags() byte {
	var b byte
	b |= byte(s.qoS.QoS()) << 6
	return b
}

func (s *SubackMessage) Pack() (bytes []byte) {
	bytes = append(bytes, s.PackHeader()...)
	bytes = append(bytes, s.encodeFlags())
	bytes = append(bytes, U162b(s.TopicId())...)
	bytes = append(bytes, U162b(s.MsgId())...)
	bytes = append(bytes, s.MsgReturnCode())
	return
}

func (s *SubackMessage) Unpack(msg []byte) Message {
	return s
}

/***************
 * Unsubscribe *
 ***************/

type unsubscribeMessage struct {
	Header
	topicIdType
	topicId
	msgId
	topicName
}

func (u *unsubscribeMessage) Pack() []byte {
	return u.PackHeader()
}

func (u *unsubscribeMessage) Unpack(msg []byte) Message {
	return u
}

/************
 * Unsuback *
 ************/

type unsubackMessage struct {
	Header
	msgId
}

func (u *unsubackMessage) Pack() []byte {
	return u.PackHeader()
}

func (u *unsubackMessage) Unpack(msg []byte) Message {
	return u
}

/***********
 * Pingreq *
 ***********/

type pingreqMessage struct {
	Header
	clientId
}

func (p *pingreqMessage) Pack() []byte {
	return p.PackHeader()
}

func (p *pingreqMessage) Unpack(msg []byte) Message {
	return p
}

/************
 * Pingresp *
 ************/

type pingrespMessage struct {
	Header
}

func (p *pingrespMessage) Pack() []byte {
	return p.PackHeader()
}

func (p *pingrespMessage) Unpack(msg []byte) Message {
	return p
}

/**************
 * Disconnect *
 **************/

type DisconnectMessage struct {
	Header
	duration
}

func (d *DisconnectMessage) Pack() (bytes []byte) {
	bytes = append(bytes, d.PackHeader()...)
	bytes = append(bytes, U162b(d.Duration())...)
	return
}

func (d *DisconnectMessage) Unpack(msg []byte) Message {
	msg = d.UnpackHeader(msg)
	if len(msg) == 2 {
		d.SetDuration(B2u16(msg))
	}
	return d
}

/*********************
 * Will Topic Update *
 *********************/

type willTopicUpdateMessage struct {
	Header
	qoS
	retain
	willTopic
}

func (w *willTopicUpdateMessage) Pack() []byte {
	return w.PackHeader()
}

func (w *willTopicUpdateMessage) Unpack(msg []byte) Message {
	return w
}

/*******************
 * Will Msg Update *
 *******************/

type willMsgUpdateMessage struct {
	Header
	willMsg
}

func (w *willMsgUpdateMessage) Pack() []byte {
	return w.PackHeader()
}

func (w *willMsgUpdateMessage) Unpack(msg []byte) Message {
	return w
}

/*******************
 * Will Topic Resp *
 *******************/

type willTopicRespMessage struct {
	Header
	msgReturnCode
}

func (w *willTopicRespMessage) Pack() []byte {
	return w.PackHeader()
}

func (w *willTopicRespMessage) Unpack(msg []byte) Message {
	return w
}

/*****************
 * Will Msg Resp *
 *****************/

type willMsgRespMessage struct {
	Header
	msgReturnCode
}

func (w *willMsgRespMessage) Pack() []byte {
	return w.PackHeader()
}

func (w *willMsgRespMessage) Unpack(msg []byte) Message {
	return w
}

/* MsgType */
const (
	ADVERTISE     MsgType = 0x00
	SEARCHGW      MsgType = 0x01
	GWINFO        MsgType = 0x02
	CONNECT       MsgType = 0x04
	CONNACK       MsgType = 0x05
	WILLTOPICREQ  MsgType = 0x06
	WILLTOPIC     MsgType = 0x07
	WILLMSGREQ    MsgType = 0x08
	WILLMSG       MsgType = 0x09
	REGISTER      MsgType = 0x0A
	REGACK        MsgType = 0x0B
	PUBLISH       MsgType = 0x0C
	PUBACK        MsgType = 0x0D
	PUBCOMP       MsgType = 0x0E
	PUBREC        MsgType = 0x0F
	PUBREL        MsgType = 0x10
	SUBSCRIBE     MsgType = 0x12
	SUBACK        MsgType = 0x13
	UNSUBSCRIBE   MsgType = 0x14
	UNSUBACK      MsgType = 0x15
	PINGREQ       MsgType = 0x16
	PINGRESP      MsgType = 0x17
	DISCONNECT    MsgType = 0x18
	WILLTOPICUPD  MsgType = 0x1A
	WILLTOPICRESP MsgType = 0x1B
	WILLMSGUPD    MsgType = 0x1C
	WILLMSGRESP   MsgType = 0x1D
	/* 0x03 is reserved */
	/* 0x11 is reserved */
	/* 0x19 is reserved */
	/* 0x1E - 0xFD is reserved */
	/* 0xFE - Encapsulated message */
	/* 0xFF is reserved */
)

func (mtype MsgType) String() string {
	switch mtype {
	case ADVERTISE:
		return "ADVERTISE"
	case SEARCHGW:
		return "SEARCHGW"
	case GWINFO:
		return "GWINFO"
	case CONNECT:
		return "CONNECT"
	case CONNACK:
		return "CONNACK"
	case WILLTOPICREQ:
		return "WILLTOPICREQ"
	case WILLTOPIC:
		return "WILLTOPIC"
	case WILLMSGREQ:
		return "WILLMSGREQ"
	case WILLMSG:
		return "WILLMSG"
	case REGISTER:
		return "REGISTER"
	case REGACK:
		return "REGACK"
	case PUBLISH:
		return "PUBLISH"
	case PUBACK:
		return "PUBACK"
	case PUBCOMP:
		return "PUBCOMP"
	case PUBREC:
		return "PUBREC"
	case PUBREL:
		return "PUBREL"
	case SUBSCRIBE:
		return "SUBSCRIBE"
	case SUBACK:
		return "SUBACK"
	case UNSUBSCRIBE:
		return "UNSUBSCRIBE"
	case UNSUBACK:
		return "UNSUBACK"
	case PINGREQ:
		return "PINGREQ"
	case PINGRESP:
		return "PINGRESP"
	case DISCONNECT:
		return "DISCONNECT"
	case WILLTOPICUPD:
		return "WILLTOPICUPD"
	case WILLTOPICRESP:
		return "WILLTOPICRESP"
	case WILLMSGUPD:
		return "WILLMSGUPD"
	case WILLMSGRESP:
		return "WILLMSGRESP"
	}
	return "INVALID"
}
