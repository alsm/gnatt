package gateway

import (
	"fmt"
	"net"
	"sync"
)

type gwStatus byte

const (
	gw_stopped gwStatus = iota
	gw_starting
	gw_running
	gw_stopping
)

type Gateway interface {
	Start()
	Port() int
	OnPacket(int, []byte, *net.UDPConn, *net.UDPAddr)
}

type topicNames struct {
	sync.RWMutex
	contents map[string]uint16
	next     uint16
}

func (repo *topicNames) contains(topic string) bool {
	return repo.get(topic) != 0
}

func (repo *topicNames) get(topic string) uint16 {
	defer repo.RUnlock()
	repo.RLock()
	topicid := repo.contents[topic]
	fmt.Printf("get[%s] -> %d\n", topic, topicid)
	return topicid
}

func (repo *topicNames) put(topic string) uint16 {
	defer repo.Unlock()
	repo.Lock()
	repo.next++
	repo.contents[topic] = repo.next
	fmt.Printf("put[%s] -> %d\n", topic, repo.next)
	return repo.next
}
