package gateway

import (
	//"fmt"
	"sync"
)

type TopicTree struct {
	sync.RWMutex
	root *node
}

type node struct {
	clients  []*Client
	children map[string]*node
}

func newNode() *node {
	return &node{
		make([]*Client, 0),
		make(map[string]*node),
	}
}

// return true if level needed to be created, false otherwise
func (n *node) goTo(level string) (*node, bool) {
	created := false
	if n.children[level] == nil {
		n.children[level] = newNode()
		created = true
	}
	return n.children[level], created
}

func (n *node) addClient(client *Client) {
	n.clients = append(n.clients, client)
}

func NewTopicTree() *TopicTree {
	var t TopicTree
	// Annoyingly, the spec allows for an empty
	// string root (but no empty string levels after that)
	// ex: "/a/b" and "a/b"
	// do NOT match, and "a//b" is not allowed.
	// Consequently, "/#" and "b" do NOT match because
	// they contain a different number of levels
	t.root = newNode()
	return &t
}

func (tt *TopicTree) AddSubscription(client *Client, topic string) error {
	defer tt.Unlock()
	tt.Lock()
	//fmt.Printf("AddSubscription(\"%s\", \"%s\")\n", client.ClientId, topic)
	if levels, e := ValidateSubscribeTopicName(topic); e != nil {
		return e
	} else {
		n := tt.root
		// walk the tree following the path of topic, creating new
		// nodes as necessary
		for _, level := range levels {
			n, _ = n.goTo(level)
		}
		n.addClient(client)
	}
	return nil
}

func (tt *TopicTree) RemoveSubscription(s *Client, topicId ...uint16) error {
	defer tt.Unlock()
	tt.Lock()
	return nil
}

func (tt *TopicTree) ClientsOf(topic string) []*Client {
	defer tt.RUnlock()
	tt.RLock()
	var Clients []*Client
	return Clients
}
