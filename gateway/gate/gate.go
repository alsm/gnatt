package gateway

import (
	"fmt"
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
	OnPacket(int, []byte, uConn, uAddr)
}

// This needs to be efficient for indexing by topicId.
// However, it is necessary when adding a new topic to index
// by topic name (to check if it already exists). We optimze
// for the former case.
type topicNames struct {
	sync.RWMutex
	contents map[uint16]string
	next     uint16
}

// O(n)
func (repo *topicNames) containsTopic(topic string) bool {
	return repo.getId(topic) != 0
}

// O(1)
func (repo *topicNames) containsId(id uint16) bool {
	return repo.getTopic(id) != ""
}

// O(n)
func (repo *topicNames) getId(topic string) uint16 {
	defer repo.RUnlock()
	repo.RLock()
	var topicid uint16
	for id, topicVal := range repo.contents {
		if topicVal == topic {
			topicid = id
			break
		}
	}
	fmt.Printf("get[%s] -> %d\n", topic, topicid)
	return topicid
}

// O(1)
func (repo *topicNames) getTopic(id uint16) string {
	defer repo.RUnlock()
	repo.RLock()
	topic := repo.contents[id]
	fmt.Printf("getTopic[%d] -> %s\n", id, topic)
	return topic
}

// O(1)
func (repo *topicNames) putTopic(topic string) uint16 {
	defer repo.Unlock()
	repo.Lock()
	repo.next++
	repo.contents[repo.next] = topic
	fmt.Printf("put[%d] -> %s\n", repo.next, topic)
	return repo.next
}
