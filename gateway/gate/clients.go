package gateway

import (
	"fmt"
	"sync"
)

type StorableClient interface {
	String() string
	AddrStr() string
}

type Clients struct {
	sync.RWMutex
	// indexed by "address:port"
	clients map[string]StorableClient
}

func (c *Clients) GetClient(addr uAddr) StorableClient {
	defer c.RUnlock()
	c.RLock()
	return c.clients[addr.r.String()]
}

// Return true if this is a new client, false otherwise
// Clients are indexed by their address:port b/c
// that's the only indentifying information we have
// outside of a CONNECT packet
func (c *Clients) AddClient(client StorableClient) bool {
	defer c.Unlock()
	c.Lock()
	fmt.Printf("AddClient(%s - %s)\n", client, client.AddrStr())
	isNew := false
	if c.clients[client.AddrStr()] == nil {
		isNew = true
	}
	//todo: what to do if clientid is in use?
	//     is there some cleanup involved in topictree?
	c.clients[client.AddrStr()] = client
	return isNew
}

func (c *Clients) RemoveClient(id string) {
	defer c.Unlock()
	c.Lock()
	fmt.Printf("RemoveClient(%s)\n", id)
	delete(c.clients, id)
}
