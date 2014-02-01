package gateway

import (
	"fmt"
	. "github.com/alsm/gnatt/common/protocol"
	"net"
	"sync"
)

type Client struct {
	ClientId string
	Socket   *net.UDPConn
	Address  *net.UDPAddr
}

func NewClient(id string, c *net.UDPConn, a *net.UDPAddr) *Client {
	return &Client{
		id,
		c,
		a,
	}
}

func (c *Client) Write(m Message) error {
	_, e := c.Socket.WriteToUDP(m.Pack(), c.Address)
	return e
}

type Clients struct {
	sync.RWMutex
	clients map[string]*Client
}

// Return true if this is a new clientid, false otherwise
func (c *Clients) AddClient(client *Client) bool {
	defer c.Unlock()
	c.Lock()
	fmt.Printf("AddClient(%s)\n", client.ClientId)
	isNew := false
	if c.clients[client.ClientId] == nil {
		isNew = true
	}
	//todo: what to do if clientid is in use?
	//     is there some cleanup involved in topictree?
	c.clients[client.ClientId] = client
	return isNew
}

func (c *Clients) RemoveClient(id string) {
	defer c.Unlock()
	c.Lock()
	fmt.Printf("RemoveClient(%s)\n", id)
	delete(c.clients, id)
}
