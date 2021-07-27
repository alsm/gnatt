package gnatt

import (
	"time"

	. "github.com/petrue/gnatt/packets"
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
	err      error
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
	case CONNECT:
		return &ConnectToken{baseToken: baseToken{complete: make(chan struct{})}}
	case REGISTER:
		return &RegisterToken{baseToken: baseToken{complete: make(chan struct{})}}
	case SUBSCRIBE:
		return &SubscribeToken{baseToken: baseToken{complete: make(chan struct{})}}
	case PUBLISH:
		return &PublishToken{baseToken: baseToken{complete: make(chan struct{})}}
	case UNSUBSCRIBE:
		return &UnsubscribeToken{baseToken: baseToken{complete: make(chan struct{})}}
	case WILLTOPIC, WILLMSG, WILLMSGUPD, WILLTOPICUPD:
		return &WillToken{baseToken: baseToken{complete: make(chan struct{})}}
	}
	return nil
}

type ConnectToken struct {
	baseToken
	ReturnCode byte
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
	ReturnCode byte
}

type SubscribeToken struct {
	baseToken
	handler    MessageHandler
	topicType  byte
	Qos        byte
	TopicName  string
	TopicId    uint16
	ReturnCode byte
}

type UnsubscribeToken struct {
	baseToken
}

type WillToken struct {
	baseToken
}
