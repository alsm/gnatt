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

func (c *Client) String() string {
	return c.ClientId
}

type Clients struct {
	sync.RWMutex
	// indexed by "address:port"
	clients map[string]*Client
}

func (c *Clients) GetClient(r uAddr) *Client {
	defer c.RUnlock()
	c.RLock()
	return c.clients[r.r.String()]
}

// Return true if this is a new client, false otherwise
// Clients are indexed by their address:port b/c
// that's the only indentifying information we have
// outside of a CONNECT packet
func (c *Clients) AddClient(client *Client) bool {
	defer c.Unlock()
	c.Lock()
	fmt.Printf("AddClient(%s - %s)\n", client.ClientId, client.Address.r)
	isNew := false
	if c.clients[client.Address.r.String()] == nil {
		isNew = true
	}
	//todo: what to do if clientid is in use?
	//     is there some cleanup involved in topictree?
	c.clients[client.Address.r.String()] = client
	return isNew
}

func (c *Clients) RemoveClient(id string) {
	defer c.Unlock()
	c.Lock()
	fmt.Printf("RemoveClient(%s)\n", id)
	delete(c.clients, id)
}
