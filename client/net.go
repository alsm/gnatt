package gnatt

import (
	. "github.com/petrue/gnatt/packets"
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
				ca := m.(*ConnackMessage)
				c.setState(CONNECTED)
				ct := c.suTokens[CONNECT].(*ConnectToken)
				ct.ReturnCode = ca.ReturnCode
				ct.flowComplete()
			case REGISTER:
				r := m.(*RegisterMessage)
				c.RegisteredTopics[string(r.TopicName)] = r.TopicId
				ra := NewMessage(REGACK).(*RegackMessage)
				ra.MessageId = r.MessageId
				ra.TopicId = r.TopicId
				ra.ReturnCode = ACCEPTED
				c.outgoing <- &MessageAndToken{m: ra, t: nil}
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
					if sa.TopicId > 0x00 {
						t.TopicId = sa.TopicId
						switch t.topicType {
						case 0x01:
							c.PredefinedMessageHandlers[sa.TopicId] = t.handler
						case 0x00, 0x02:
							c.RegisteredTopics[t.TopicName] = sa.TopicId
							c.MessageHandlers[sa.TopicId] = t.handler
						}
					}
				default:
					ERROR.Println(NET, sa.ReturnCode, "for SUBSCRIBE to", t.TopicName)
				}
				t.flowComplete()
			case PUBLISH:
				p := m.(*PublishMessage)
				switch p.TopicIdType {
				case 0x00:
					if handler, ok := c.MessageHandlers[p.TopicId]; ok {
						go handler(c, p)
					}
				default:
					if c.DefaultMessageHandler != nil {
						go c.DefaultMessageHandler(c, p)
					}
				}
			case PUBACK:
				pa := m.(*PubackMessage)
				t := c.MessageIds.getToken(pa.MessageId).(*PublishToken)
				t.ReturnCode = pa.ReturnCode
				t.TopicId = pa.TopicId
				t.flowComplete()
			case DISCONNECT:
				c.setState(UNCONNECTED)
				DEBUG.Println(NET, "Received DISCONNECT, closing socket")
				c.conn.Close()
			case WILLTOPICREQ:
				wm := NewMessage(WILLTOPIC).(*WillTopicMessage)
				wm.WillTopic = []byte(c.will.Topic)
				wm.Qos = c.will.Qos
				wm.Retain = c.will.Retain
				c.outgoing <- &MessageAndToken{m: wm, t: nil}
			case WILLMSGREQ:
				wm := NewMessage(WILLMSG).(*WillMsgMessage)
				wm.WillMsg = c.will.Data
				c.outgoing <- &MessageAndToken{m: wm, t: nil}
			}
		}
	}
}
