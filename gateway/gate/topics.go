package gateway

import (
	"fmt"
	"strings"
	"sync"
)

func ContainsWildcard(topic string) bool {
	if len(topic) == 1 && (topic == "+" || topic == "#") {
		return true
	}
	if len(topic) > 1 && (topic[len(topic)-2:] == "/#" || topic[len(topic)-2:] == "/+") {
		return true
	}
	return strings.Contains(topic, "/+/")
}

func ValidateSubscribeTopicName(topic string) ([]string, error) {
	if len(topic) == 0 {
		return nil, fmt.Errorf("topic may not be empty string")
	}
	if topic[len(topic)-1] == '/' {
		return nil, fmt.Errorf("topic may not end with level seperator")
	}

	levels := strings.Split(topic, "/")
	for i, level := range levels {
		if level == "" && i != 0 {
			return nil, fmt.Errorf("non-root topic level may not be empty string")
		}
		if level == "#" && i != len(levels)-1 {
			return nil, fmt.Errorf("multi-level wild card must be last level in topic")
		}
	}
	return levels, nil
}

func ValidatePublishTopicName(topic string) ([]string, error) {
	if len(topic) == 0 {
		return nil, fmt.Errorf("topic may not be empty string")
	}
	if topic[len(topic)-1] == '/' {
		return nil, fmt.Errorf("topic may not end with level seperator")
	}

	levels := strings.Split(topic, "/")
	for i, level := range levels {
		if level == "#" || level == "+" {
			return nil, fmt.Errorf("topic may not contain wild card character")
		}
		if level == "" && i != 0 {
			return nil, fmt.Errorf("non-root topic level may not be empty string")
		}
	}
	return levels, nil
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
