package gateway

import (
	"fmt"
	"sync"

	. "github.com/alsm/gnatt/common/protocol"
)

type Client struct {
	sync.RWMutex
	ClientId         string
	Socket           uConn
	Address          uAddr
	registeredTopics map[uint16]bool
	pendingMessages  map[uint16]*PublishMessage
}

func NewClient(id string, c uConn, a uAddr) *Client {
	fmt.Printf("NewClient, id: \"%s\"\n", id)
	return &Client{
		sync.RWMutex{},
		id,
		c,
		a,
		make(map[uint16]bool),
		make(map[uint16]*PublishMessage),
	}
}

func (c *Client) Write(m Message) error {
	_, e := c.Socket.c.WriteToUDP(m.Pack(), c.Address.r)
	return e
}

func (c *Client) Register(topicId uint16) {
	defer c.Unlock()
	c.Lock()
	fmt.Printf("client %s registered topicId %d\n", c.ClientId, topicId)
	c.registeredTopics[topicId] = true
}

func (c *Client) Registered(topicId uint16) bool {
	defer c.RUnlock()
	c.RLock()
	return c.registeredTopics[topicId]
}

func (c *Client) AddPendingMessage(p *PublishMessage) {
	defer c.Unlock()
	c.Lock()
	c.pendingMessages[p.TopicId()] = p
}

func (c *Client) FetchPendingMessage(topicid uint16) *PublishMessage {
	defer c.Unlock()
	c.Lock()
	pm := c.pendingMessages[topicid]
	delete(c.pendingMessages, topicid)
	return pm
}

func (c *Client) AddrStr() string {
	return c.Address.r.String()
}

func (c *Client) String() string {
	return c.ClientId
}
