package gnatt

import (
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
	ClientId string
	conn     net.Conn
	incoming chan Message
	outgoing chan Message
	stop     chan struct{}
}

func NewClient(server string) *SNClient {
	c := &SNClient{}
	serverURI, _ := url.Parse(server)
	c.conn, _ = net.Dial("udp", serverURI.Host)
	c.incoming = make(chan Message)
	c.outgoing = make(chan Message)
	c.stop = make(chan struct{})
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
	c.outgoing <- cp
}
