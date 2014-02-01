package gateway

import (
	"sync"
)

type TopicTree struct {
	sync.RWMutex
}

type node struct {
	text     string
	clients  []*Client
	children []*node
}

func (tt *TopicTree) AddSubscription(s *Client, topicid uint16) {
}

func (tt *TopicTree) RemoveSubscription(s *Client, topicId ...uint16) error {
	return nil
}

func (tt *TopicTree) ClientsOf(topicId uint16) []*Client {
	var Clients []*Client
	return Clients
}
