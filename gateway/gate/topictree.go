package gateway

import (
	"fmt"
	"sync"
	"strings"
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
	fmt.Printf("goTo(\"%s\")\n", level)
	created := false
	if n.children[level] == nil {
		n.children[level] = newNode()
		created = true
	}
	return n.children[level], created
}

// return true if this is the first client to be added
// to this node (representing a subscription)
func (n *node) addClient(client *Client) bool {
	isFirst := len(n.clients) == 0
	n.clients = append(n.clients, client)
	fmt.Printf("addClient(\"%s\")\n", client.ClientId)
	return isFirst
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

// return false if there is already a subscriber of this topic,
// true if this is the first subscriber
func (tt *TopicTree) AddSubscription(client *Client, topic string) (bool, error) {
	defer tt.Unlock()
	tt.Lock()
	fmt.Printf("AddSubscription(\"%s\", \"%s\")\n", client.ClientId, topic)
	if levels, e := ValidateSubscribeTopicName(topic); e != nil {
		return false, e
	} else {
		n := tt.root
		// walk the tree following the path of topic, creating new
		// nodes as necessary
		for _, level := range levels {
			n, _ = n.goTo(level)
		}
		return n.addClient(client), nil
	}
}

func (tt *TopicTree) RemoveSubscription(s *Client, topicId ...uint16) error {
	defer tt.Unlock()
	tt.Lock()
	return nil
}

func (tt *TopicTree) SubscribersOf(topic string) (clients []*Client) {
	defer tt.RUnlock()
	tt.RLock()

	n := tt.root
	levels := strings.Split(topic, "/")
	subscribers(n, levels, clients)
	return
}

// good luck understanding this
//    levels is the topic string (tokenized). It will not have wild cards.
//    The TopicTree must be walked down, as long as one of the children in the
//    next level is a wildcard or a level match.
//      If a child is "#", we add the clients of that child. The recursive stepping continues.
//      If a child is "+", we then start a recursive tree at the next level.
//      If a child == levels[0], we recursivly step to the next level.
//      Otherwise we just return and pop back up.
func subscribers(n *node, levels []string, clients []*Client) {
	fmt.Printf("subscribers: levels: %v\n", levels)
	fmt.Printf("len n.clients: %d\n", len(n.clients))
	fmt.Printf("len n.children: %d\n", len(n.children))

	if len(levels) == 0 {
		return
	}

	if hash := n.children["#"]; hash != nil {
		clients = append(clients, hash.clients...)
	}

	if plus := n.children["+"]; plus != nil {
		subscribers(plus, levels[1:], clients)
	}

	if len(levels) == 1 {
		fmt.Println("len(levels)==1")
		if last := n.children[levels[0]]; last != nil {
			fmt.Println("adding last level childrens")
			clients = append(clients, last.clients...)
		} else if lastH := n.children["#"]; lastH != nil {
			fmt.Println("adding last level hashes")
			clients = append(clients, lastH.clients...)
		} else if lastP := n.children["+"]; lastP != nil {
			fmt.Println("adding last level pluses")
			clients = append(clients, lastP.clients...)
		}
		return // needed?
	}

	if next := n.children[levels[0]]; next != nil {
		subscribers(next, levels[1:], clients)
	}
}
