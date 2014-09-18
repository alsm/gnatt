package gnatt

import (
	. "github.com/alsm/gnatt/packets"
	"time"
)

type MessageAndToken struct {
	m Message
	t Token
}

type Token interface {
	Wait()
	WaitTimeout(time.Duration)
	flowComplete()
}

type baseToken struct {
	complete chan struct{}
	ready    bool
}

// Wait will wait indefinitely for the Token to complete, ie the Publish
// to be sent and confirmed receipt from the broker
func (b *baseToken) Wait() {
	if !b.ready {
		<-b.complete
		b.ready = true
	}
}

// WaitTimeout takes a time in ms
func (b *baseToken) WaitTimeout(d time.Duration) {
	if !b.ready {
		select {
		case <-b.complete:
			b.ready = true
		case <-time.After(d):
		}
	}
}

func (b *baseToken) flowComplete() {
	close(b.complete)
}

func newToken(tType byte) Token {
	switch tType {
	case REGISTER:
		return &RegisterToken{baseToken: baseToken{complete: make(chan struct{})}}
	case SUBSCRIBE:
		return &SubscribeToken{baseToken: baseToken{complete: make(chan struct{})}}
	case PUBLISH:
		return &PublishToken{baseToken: baseToken{complete: make(chan struct{})}}
	case UNSUBSCRIBE:
		return &UnsubscribeToken{baseToken: baseToken{complete: make(chan struct{})}}
	}
	return nil
}

type RegisterToken struct {
	baseToken
	TopicName  string
	TopicId    uint16
	ReturnCode byte
}

type PublishToken struct {
	baseToken
	TopicId    uint16
	ReturnCode uint16
}

type SubscribeToken struct {
	baseToken
	Qos        byte
	TopicName  string
	TopicId    uint16
	ReturnCode byte
}

type UnsubscribeToken struct {
	baseToken
}
