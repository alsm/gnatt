package gateway

import (
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

// return true if this is the first client to be added
// to this node (representing a subscription)
func (n *node) addClient(client *Client) bool {
	isFirst := len(n.clients) == 0
	n.clients = append(n.clients, client)
	return isFirst
}

func NewTopicTree() *TopicTree {
	var t TopicTree
	t.root = newNode()
	return &t
}

// return false if there is already a subscriber of this topic,
// true if this is the first subscriber
func (tt *TopicTree) AddSubscription(client *Client, topic string) (bool, error) {
	defer tt.Unlock()
	tt.Lock()
	INFO.Printf("AddSubscription(\"%s\", \"%s\")\n", client.ClientId, topic)
	if levels, e := ValidateTopicFilter(topic); e != nil {
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

// topic could contain wild cards, however we do only consider the literal
// topic string - (wilds are not evaluated for this)
func (tt *TopicTree) RemoveSubscription(s *Client, topic string) error {
	defer tt.Unlock()
	tt.Lock()
	if levels, e := ValidateTopicFilter(topic); e != nil {
		return e
	} else {
		n := tt.root
		for _, level := range levels {
			if n = n.children[level]; n == nil {
				ERROR.Printf("no subscription exists \"%s\"\n", topic)
				return ErrNoSuchSubscriptionExists
			}
		}
		if len(n.clients) < 1 {
			ERROR.Printf("no clients of subscription \"%s\"\n", topic)
			return ErrNoSubscribers
		}
		for i := 0; i < len(n.clients); i++ {
			if n.clients[i].ClientId == s.ClientId {
				// inexpensive way of removing from a slice
				n.clients[i] = n.clients[len(n.clients)-1]
				n.clients = n.clients[0 : len(n.clients)-1]
				INFO.Printf("deleted subscription of client \"%s\"\n", s.ClientId)
				return nil
			}
		}
		ERROR.Printf("client \"%s\" was not subscribed to \"%s\"\n", s.ClientId, topic)
		return ErrClientNotSubscribed
	}
}

// topic MUST be valid (ie no wild cards, no empty level, no ending slash)
/***! Hey dipstick, read the above comment, !***/
/***! that's where your bug is coming from. !***/
func (tt *TopicTree) SubscribersOf(topic string) ([]*Client, error) {
	defer tt.RUnlock()
	tt.RLock()

	clients := make([]*Client, 0)

	n := tt.root
	if levels, e := ValidateTopicName(topic); e != nil {
		return nil, e
	} else {
		subscribers(n, levels, &clients)
		return clients, nil
	}
}

func subscribers(n *node, levels []string, clients *[]*Client) {
	if len(levels) == 0 {
		return
	}

	if hash := n.children["#"]; hash != nil {
		*clients = append(*clients, hash.clients...)
	}

	if plus := n.children["+"]; plus != nil {
		if len(levels) == 1 {
			*clients = append(*clients, plus.clients...)
		} else {
			subscribers(plus, levels[1:], clients)
		}
	}

	if match := n.children[levels[0]]; match != nil {
		if len(levels) == 1 {
			*clients = append(*clients, match.clients...)
		} else {
			subscribers(match, levels[1:], clients)
		}
	}
}
