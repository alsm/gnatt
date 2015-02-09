package gateway

import (
	"net"
	"sync"
)

type Clients struct {
	sync.RWMutex
	// indexed by "address:port" => StorableClient
	clients map[string]SNClient
}

func (c *Clients) GetClient(addr *net.UDPAddr) SNClient {
	defer c.RUnlock()
	c.RLock()
	return c.clients[addr.String()]
}

// Return true if this is a new client, false otherwise
// Clients are indexed by their address:port b/c
// that's the only indentifying information we have
// outside of a CONNECT packet
func (c *Clients) AddClient(client SNClient) bool {
	defer c.Unlock()
	c.Lock()
	addr := client.AddrString()
	INFO.Println("AddClient(%s - %s)", client, addr)
	isNew := false
	if c.clients[addr] == nil {
		isNew = true
	}
	//todo: what to do if clientid is in use?
	//     is there some cleanup involved in topictree?
	c.clients[addr] = client
	return isNew
}

func (c *Clients) RemoveClient(id string) {
	defer c.Unlock()
	c.Lock()
	INFO.Println("RemoveClient(%s)", id)
	delete(c.clients, id)
}
