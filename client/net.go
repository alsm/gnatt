package gnatt

import (
	. "github.com/alsm/gnatt/packets"
)

func (c *SNClient) receive() {
	var err error
	var m Message

	DEBUG.Println(NET, "started receive()")
	for {
		if m, err = ReadPacket(c.conn); err != nil {
			break
		}
		DEBUG.Println(NET, "Received", MessageNames[m.MessageType()])
		c.incoming <- m
	}

	select {
	case <-c.stop:
		DEBUG.Println(NET, "stopped receive()")
	default:
		ERROR.Println(NET, "stopped receive() due to error")
	}
	return
}

func (c *SNClient) send() {
	DEBUG.Println(NET, "started send()")
	for {
		m := <-c.outgoing
		DEBUG.Println(NET, "sending message", MessageNames[m.MessageType()])
		switch m.MessageType() {
		case REGISTER:
			m.(*RegisterMessage).MessageId = 0x01
		case PUBLISH:
			m.(*PublishMessage).MessageId = 0x01
		case SUBSCRIBE:
			m.(*SubscribeMessage).MessageId = 0x01
		case UNSUBSCRIBE:
			m.(*UnsubscribeMessage).MessageId = 0x01
		}
		m.Write(c.conn)
		if m.MessageType() == DISCONNECT {
			DEBUG.Println(NET, "Sent DISCONNECT, closing connection")
			c.conn.Close()
			return
		}
	}
}

func (c *SNClient) handle() {
	DEBUG.Println(NET, "started handle()")
	for {
		select {
		case m := <-c.incoming:
			DEBUG.Println(NET, "got message off <-incoming", MessageNames[m.MessageType()])
			switch m.MessageType() {
			case CONNACK:
			case REGACK:
			case SUBACK:
			case PUBLISH:
			case PUBACK:
			case DISCONNECT:
			}
		}
	}
}
