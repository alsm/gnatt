package gnatt

import (
	"errors"
	. "github.com/alsm/gnatt/packets"
	"net"
	"net/url"
	"sync"
)

type MessageHandler func(client Client, message *PublishMessage)

type Client interface {
	Connect()
}

type SNClient struct {
	sync.RWMutex
	ClientId            string
	RegisteredTopics    map[string]uint16
	OutstandingMessages map[uint16]Message
	MessageIds          mids
	conn                net.Conn
	incoming            chan Message
	outgoing            chan *MessageAndToken
	stop                chan struct{}
}

func NewClient(server string) *SNClient {
	c := &SNClient{}
	serverURI, _ := url.Parse(server)
	c.conn, _ = net.Dial("udp", serverURI.Host)
	c.RegisteredTopics = make(map[string]uint16)
	c.incoming = make(chan Message)
	c.outgoing = make(chan *MessageAndToken)
	c.stop = make(chan struct{})
	c.MessageIds.index = make(map[uint16]Token)
	go c.receive()
	go c.send()
	go c.handle()

	return c
}

func (c *SNClient) Connect() {
	cp := NewMessage(CONNECT).(*ConnectMessage)
	cp.CleanSession = true
	cp.ClientId = []byte(c.ClientId)
	cp.Duration = 30
	cp.Write(c.conn)
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

func (c *SNClient) Subscribe(topic string, qos byte) *SubscribeToken {
	t := newToken(SUBSCRIBE).(*SubscribeToken)
	t.TopicName = topic
	s := NewMessage(SUBSCRIBE).(*SubscribeMessage)
	if len(topic) > 2 {
		s.TopicIdType = 0x00
	} else {
		s.TopicIdType = 0x02
	}
	s.TopicName = []byte(topic)
	s.Qos = qos
	c.outgoing <- &MessageAndToken{m: s, t: t}
	return t
}

func (c *SNClient) SubscribePredefined(topicid uint16, qos byte) *SubscribeToken {
	t := newToken(SUBSCRIBE).(*SubscribeToken)
	s := NewMessage(SUBSCRIBE).(*SubscribeMessage)
	s.TopicIdType = 0x01
	s.TopicId = topicid
	s.Qos = qos
	c.outgoing <- &MessageAndToken{m: s, t: t}
	return t
}
