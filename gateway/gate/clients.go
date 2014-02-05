package gateway

import (
	"fmt"
	. "github.com/alsm/gnatt/common/protocol"
	"sync"
)

type Client struct {
	ClientId string
	Socket   uConn // i honestly don't know what this actually is
	Address  uAddr
}

func NewClient(id string, c uConn, a uAddr) *Client {
	return &Client{
		id,
		c,
		a,
	}
}

func (c *Client) Write(m Message) error {
	_, e := c.Socket.c.WriteToUDP(m.Pack(), c.Address.r)
	return e
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
