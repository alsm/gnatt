package gateway

import (
	"sync"
)

type subscriber int // uhmmm

type TopicTree struct {
	sync.RWMutex
}

type node struct {
	text        string
	subscribers []*subscriber
	children    []*node
}

func (tt *TopicTree) AddSubscription(s *subscriber, topicid uint16) {
}

func (tt *TopicTree) RemoveSubscription(s *subscriber, topicId ...uint16) error {
	return nil
}

func (tt *TopicTree) SubscribersOf(topicId uint16) []*subscriber {
	var subscribers []*subscriber
	return subscribers
}
