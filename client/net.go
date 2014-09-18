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
		mt := <-c.outgoing
		DEBUG.Println(NET, "sending message", MessageNames[mt.m.MessageType()])
		switch mt.m.MessageType() {
		case REGISTER:
			mt.m.(*RegisterMessage).MessageId = c.MessageIds.getId(mt.t)
		case PUBLISH:
			mt.m.(*PublishMessage).MessageId = c.MessageIds.getId(mt.t)
		case SUBSCRIBE:
			mt.m.(*SubscribeMessage).MessageId = c.MessageIds.getId(mt.t)
		case UNSUBSCRIBE:
			mt.m.(*UnsubscribeMessage).MessageId = c.MessageIds.getId(mt.t)
		}
		mt.m.Write(c.conn)
		if mt.m.MessageType() == DISCONNECT {
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
			case REGISTER:
				r := m.(*RegisterMessage)
				c.RegisteredTopics[string(r.TopicName)] = r.TopicId
			case REGACK:
				ra := m.(*RegackMessage)
				t := c.MessageIds.getToken(ra.MessageId).(*RegisterToken)
				t.ReturnCode = ra.ReturnCode
				switch ra.ReturnCode {
				case ACCEPTED:
					DEBUG.Println(NET, t.TopicName, "registered as", ra.TopicId)
					c.RegisteredTopics[t.TopicName] = ra.TopicId
					t.TopicId = ra.TopicId
				default:
					ERROR.Println(NET, ra.ReturnCode, "for REGISTER for", string(t.TopicName))
				}
				t.flowComplete()
			case SUBACK:
				sa := m.(*SubackMessage)
				t := c.MessageIds.getToken(sa.MessageId).(*SubscribeToken)
				t.ReturnCode = sa.ReturnCode
				switch sa.ReturnCode {
				case ACCEPTED:
					t.Qos = sa.Qos
					t.TopicId = sa.TopicId
				default:
					ERROR.Println(NET, sa.ReturnCode, "for SUBSCRIBE to", t.TopicName)
				}
				t.flowComplete()
			case PUBLISH:
			case PUBACK:
			case DISCONNECT:
			}
		}
	}
}
