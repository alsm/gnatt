package gnatt

import (
	"errors"
	"net"
	"sync"

	. "github.com/petrue/gnatt/packets"
)

type MessageHandler func(client *SNClient, message *PublishMessage)

type Will struct {
	Topic  string
	Data   []byte
	Qos    byte
	Retain bool
}

type Client interface {
	Connect()
}

const (
	UNCONNECTED byte = iota
	CONNECTING
	CONNECTED
)

type SNClient struct {
	sync.RWMutex
	ClientId                  string
	OutstandingMessages       map[uint16]Message
	RegisteredTopics          map[string]uint16
	MessageHandlers           map[uint16]MessageHandler
	PredefinedTopics          map[string]uint16
	PredefinedMessageHandlers map[uint16]MessageHandler
	DefaultMessageHandler     MessageHandler
	MessageIds                mids
	will                      Will
	suTokens                  map[int]Token
	conn                      net.Conn
	incoming                  chan Message
	outgoing                  chan *MessageAndToken
	stop                      chan struct{}
	state                     byte
}

func NewClient(conn net.Conn, clientid string) (*SNClient, error) {
	c := &SNClient{}
	c.ClientId = clientid

	c.conn = conn
	c.RegisteredTopics = make(map[string]uint16)
	c.MessageHandlers = make(map[uint16]MessageHandler)
	c.PredefinedTopics = make(map[string]uint16)
	c.PredefinedMessageHandlers = make(map[uint16]MessageHandler)
	c.suTokens = make(map[int]Token)
	c.incoming = make(chan Message)
	c.outgoing = make(chan *MessageAndToken)
	c.stop = make(chan struct{})
	c.MessageIds.index = make(map[uint16]Token)
	go c.receive()
	go c.send()
	go c.handle()

	return c, nil
}

func (c *SNClient) setState(s byte) {
	c.Lock()
	defer c.Unlock()
	c.state = s
}

func (c *SNClient) Connect() *ConnectToken {
	c.setState(CONNECTING)
	ct := newToken(CONNECT).(*ConnectToken)
	c.suTokens[CONNECT] = ct
	cp := NewMessage(CONNECT).(*ConnectMessage)
	cp.CleanSession = true
	cp.ClientId = []byte(c.ClientId)
	cp.Duration = 30
	if c.will.Topic != "" {
		cp.Will = true
	}
	cp.Write(c.conn)
	return ct
}

func (c *SNClient) Predefine(topic string, topicid uint16) error {
	for t, id := range c.RegisteredTopics {
		if t == topic || id == topicid {
			return errors.New("Topic or TopicId already in use")
		}
	}
	c.RegisteredTopics[topic] = topicid
	return nil
}

func (c *SNClient) Register(topic string) *RegisterToken {
	t := newToken(REGISTER).(*RegisterToken)
	t.TopicName = topic
	r := NewMessage(REGISTER).(*RegisterMessage)
	r.TopicName = []byte(topic)
	c.outgoing <- &MessageAndToken{m: r, t: t}
	return t
}

func (c *SNClient) Subscribe(topic string, qos byte, mh MessageHandler) *SubscribeToken {
	t := newToken(SUBSCRIBE).(*SubscribeToken)
	t.handler = mh
	t.TopicName = topic
	s := NewMessage(SUBSCRIBE).(*SubscribeMessage)
	if len(topic) > 2 {
		s.TopicIdType = 0x00
		t.topicType = 0x00
	} else {
		s.TopicIdType = 0x02
		t.topicType = 0x02
	}
	s.TopicName = []byte(topic)
	s.Qos = qos
	c.outgoing <- &MessageAndToken{m: s, t: t}
	return t
}

func (c *SNClient) SubscribePredefined(topicid uint16, qos byte, mh MessageHandler) *SubscribeToken {
	t := newToken(SUBSCRIBE).(*SubscribeToken)
	t.handler = mh
	t.topicType = 0x01
	s := NewMessage(SUBSCRIBE).(*SubscribeMessage)
	s.TopicIdType = 0x01
	s.TopicId = topicid
	s.Qos = qos
	c.outgoing <- &MessageAndToken{m: s, t: t}
	return t
}

func (c *SNClient) Publish(topic string, qos byte, retain bool, data []byte) *PublishToken {
	t := newToken(PUBLISH).(*PublishToken)
	p := NewMessage(PUBLISH).(*PublishMessage)
	p.TopicId = c.RegisteredTopics[topic]
	p.Qos = qos
	p.Retain = retain
	p.Data = data
	c.outgoing <- &MessageAndToken{m: p, t: t}
	return t
}

func (c *SNClient) PublishPredefined(topicid uint16, qos byte, retain bool, data []byte) *PublishToken {
	t := newToken(PUBLISH).(*PublishToken)
	p := NewMessage(PUBLISH).(*PublishMessage)
	p.TopicIdType = 0x01
	p.TopicId = topicid
	p.Qos = qos
	p.Retain = retain
	p.Data = data
	c.outgoing <- &MessageAndToken{m: p, t: t}
	return t
}

func (c *SNClient) SetWill(topic string, qos byte, retain bool, data []byte) {
	c.will.Topic = topic
	c.will.Qos = qos
	c.will.Retain = retain
	c.will.Data = make([]byte, len(data))
	copy(c.will.Data, data)
}

func (c *SNClient) SetWillTopic(t string) *WillToken {
	c.will.Topic = t
	wt := newToken(WILLTOPIC).(*WillToken)
	if c.state != UNCONNECTED {
		wm := NewMessage(WILLTOPICUPD).(*WillTopicUpdateMessage)
		wm.Qos = c.will.Qos
		wm.Retain = c.will.Retain
		c.outgoing <- &MessageAndToken{m: wm, t: wt}
		return wt
	}
	wt.flowComplete()
	return wt
}

func (c *SNClient) SetWillQos(q byte) *WillToken {
	c.will.Qos = q
	wt := newToken(WILLTOPIC).(*WillToken)
	if c.state != UNCONNECTED {
		wm := NewMessage(WILLTOPICUPD).(*WillTopicUpdateMessage)
		wm.Qos = c.will.Qos
		wm.Retain = c.will.Retain
		c.outgoing <- &MessageAndToken{m: wm, t: wt}
		return wt
	}
	wt.flowComplete()
	return wt
}

func (c *SNClient) SetWillRetain(r bool) *WillToken {
	c.will.Retain = r
	wt := newToken(WILLTOPIC).(*WillToken)
	if c.state != UNCONNECTED {
		wm := NewMessage(WILLTOPICUPD).(*WillTopicUpdateMessage)
		wm.Qos = c.will.Qos
		wm.Retain = c.will.Retain
		c.outgoing <- &MessageAndToken{m: wm, t: wt}
		return wt
	}
	wt.flowComplete()
	return wt
}

func (c *SNClient) SetWillData(d []byte) *WillToken {
	c.will.Data = make([]byte, len(d))
	copy(c.will.Data, d)
	wt := newToken(WILLMSGUPD).(*WillToken)
	if c.state != UNCONNECTED {
		wm := NewMessage(WILLMSGUPD).(*WillMsgUpdateMessage)
		c.outgoing <- &MessageAndToken{m: wm, t: wt}
		return wt
	}
	wt.flowComplete()
	return wt
}
