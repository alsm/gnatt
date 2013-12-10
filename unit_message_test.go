package gnatt

import (
	//"fmt"
	"testing"
)

func Test_MessageTypeValues(t *testing.T) {
	if ADVERTISE != 0x00 {
		t.Errorf("ADVERTISE is %x instead of %x", ADVERTISE, 0x00)
	}
	if SEARCHGW != 0x01 {
		t.Errorf("SEARCHGW is %x instead of %x", SEARCHGW, 0x01)
	}
	if GWINFO != 0x02 {
		t.Errorf("GWINFO is %x instead of %x", GWINFO, 0x02)
	}
	if CONNECT != 0x04 {
		t.Errorf("CONNECT is %x instead of %x", CONNECT, 0x04)
	}
	if CONNACK != 0x05 {
		t.Errorf("CONNACK is %x instead of %x", CONNACK, 0x05)
	}
	if WILLTOPICREQ != 0x06 {
		t.Errorf("WILLTOPICREQ is %x instead of %x", WILLTOPICREQ, 0x06)
	}
	if WILLTOPIC != 0x07 {
		t.Errorf("WILLTOPIC is %x instead of %x", WILLTOPIC, 0x07)
	}
	if WILLMSGREQ != 0x08 {
		t.Errorf("WILLMSGREQ is %x instead of %x", WILLMSGREQ, 0x08)
	}
	if WILLMSG != 0x09 {
		t.Errorf("WILLMSG is %x instead of %x", WILLMSG, 0x09)
	}
	if REGISTER != 0x0A {
		t.Errorf("REGISTER is %x instead of %x", REGISTER, 0x0A)
	}
	if REGACK != 0x0B {
		t.Errorf("REGACK is %x instead of %x", REGACK, 0x0B)
	}
	if PUBLISH != 0x0C {
		t.Errorf("PUBLISH is %x instead of %x", PUBLISH, 0x0C)
	}
	if PUBACK != 0x0D {
		t.Errorf("PUBACK is %x instead of %x", PUBACK, 0x0D)
	}
	if PUBCOMP != 0x0E {
		t.Errorf("PUBCOMP is %x instead of %x", PUBCOMP, 0x0E)
	}
	if PUBREC != 0x0F {
		t.Errorf("PUBREC is %x instead of %x", PUBREC, 0x0F)
	}
	if PUBREL != 0x10 {
		t.Errorf("PUBREL is %x instead of %x", PUBREL, 0x10)
	}
	if SUBSCRIBE != 0x12 {
		t.Errorf("SUBSCRIBE is %x instead of %x", SUBSCRIBE, 0x12)
	}
	if SUBACK != 0x13 {
		t.Errorf("SUBACK is %x instead of %x", SUBACK, 0x13)
	}
	if UNSUBSCRIBE != 0x14 {
		t.Errorf("UNSUBSCRIBE is %x instead of %x", UNSUBSCRIBE, 0x14)
	}
	if UNSUBACK != 0x15 {
		t.Errorf("UNSUBACK is %x instead of %x", UNSUBACK, 0x15)
	}
	if PINGREQ != 0x16 {
		t.Errorf("PINGREQ is %x instead of %x", PINGREQ, 0x16)
	}
	if PINGRESP != 0x17 {
		t.Errorf("PINGRESP is %x intead of %x", PINGRESP, 0x17)
	}
	if DISCONNECT != 0x18 {
		t.Errorf("DISCONNECT is %x intead of %x", DISCONNECT, 0x18)
	}
	if WILLTOPICUPD != 0x1A {
		t.Errorf("WILLTOPICUPD is %x intead of %x", WILLTOPICUPD, 0x1A)
	}
	if WILLTOPICRESP != 0x1B {
		t.Errorf("WILLTOPICRESPis %x intead of %x", WILLTOPICRESP, 0x1B)
	}
	if WILLMSGUPD != 0x1C {
		t.Errorf("WILLMSGUPD is %x intead of %x", WILLMSGUPD, 0x1C)
	}
	if WILLMSGRESP != 0x1D {
		t.Errorf("WILLMSGRESP is %x intead of %x", WILLMSGRESP, 0x1D)
	}
}

func Test_Length(t *testing.T) {
	var h Header
	m := map[int]int{
		0:     0x0000,
		1:     0x0001,
		2:     0x0002,
		3:     0x0003,
		8:     0x0008,
		15:    0x000F,
		16:    0x0010,
		127:   0x007F,
		256:   0x0100,
		16384: 0x4000,
		65535: 0xFFFF,
	}

	for in, exp := range m {
		h.SetLength(in)
		if h.Length() != exp {
			t.Errorf("encodeLength failed, input %d expected 0x%X, got 0x%X", in, exp, h.Length())
		}
	}
}

func Test_NewAdvertiseMessage(t *testing.T) {
	am := NewMessage(ADVERTISE)
	switch m := am.(type) {
	case *advertiseMessage:

	default:
		t.Errorf("msg should be *advertiseMessage, is %s", m)
	}
}

func Test_NewSearchGwMessage(t *testing.T) {
	sm := NewMessage(SEARCHGW)
	switch m := sm.(type) {
	case *searchgwMessage:

	default:
		t.Errorf("msg should be *searchgwMessage, is %s", m)
	}
}

func Test_NewGwInfoMessage(t *testing.T) {
	gm := NewMessage(GWINFO)
	switch m := gm.(type) {
	case *gwInfoMessage:

	default:
		t.Errorf("msg should be *gwInfoMessage, is %s", m)
	}
}

func Test_NewConnectMessage(t *testing.T) {
	cm := NewMessage(CONNECT)
	switch m := cm.(type) {
	case *connectMessage:

	default:
		t.Errorf("msg should be *connectMessage, is %s", m)
	}
}

func Test_NewConnackMessage(t *testing.T) {
	cm := NewMessage(CONNACK)
	switch m := cm.(type) {
	case *connackMessage:

	default:
		t.Errorf("msg should be *connackMessage, is %s", m)
	}
}

func Test_NewWillTopicReqMessage(t *testing.T) {
	wm := NewMessage(WILLTOPICREQ)
	switch m := wm.(type) {
	case *willTopicReqMessage:

	default:
		t.Errorf("msg should be *willTopicReqMessage, is %s", m)
	}
}

func Test_NewWillTopicMessage(t *testing.T) {
	wm := NewMessage(WILLTOPIC)
	switch m := wm.(type) {
	case *willTopicMessage:

	default:
		t.Errorf("msg should be *willTopicMessage, is %s", m)
	}
}

func Test_NewWillMsgReqMessage(t *testing.T) {
	wm := NewMessage(WILLMSGREQ)
	switch m := wm.(type) {
	case *willMsgReqMessage:

	default:
		t.Errorf("msg should be *willMsgReqMessage, is %s", m)
	}
}

func Test_NewWillMsgMessage(t *testing.T) {
	wm := NewMessage(WILLMSG)
	switch m := wm.(type) {
	case *willMsgMessage:

	default:
		t.Errorf("msg should be *willMsgMessage, is %s", m)
	}
}

func Test_NewRegisterMessage(t *testing.T) {
	rm := NewMessage(REGISTER)
	switch m := rm.(type) {
	case *registerMessage:

	default:
		t.Errorf("msg should be *registerMessage, is %s", m)
	}
}

func Test_NewRegackMessage(t *testing.T) {
	rm := NewMessage(REGACK)
	switch m := rm.(type) {
	case *regackMessage:

	default:
		t.Errorf("msg should be *regackMessage, is %s", m)
	}
}

func Test_NewPublishMessage(t *testing.T) {
	pm := NewMessage(PUBLISH)
	switch m := pm.(type) {
	case *publishMessage:

	default:
		t.Errorf("msg should be *publishMessage, is %s", m)
	}
}

func Test_NewPubackMessage(t *testing.T) {
	pm := NewMessage(PUBACK)
	switch m := pm.(type) {
	case *pubackMessage:

	default:
		t.Errorf("msg should be *pubackMessage, is %s", m)
	}
}

func Test_NewPubcompMessage(t *testing.T) {
	pm := NewMessage(PUBCOMP)
	switch m := pm.(type) {
	case *pubcompMessage:

	default:
		t.Errorf("msg should be *pubcompMessage, is %s", m)
	}
}

func Test_NewPubrecMessage(t *testing.T) {
	pm := NewMessage(PUBREC)
	switch m := pm.(type) {
	case *pubrecMessage:

	default:
		t.Errorf("msg should be *pubrecMessage, is %s", m)
	}
}

func Test_NewPubrelMessage(t *testing.T) {
	pm := NewMessage(PUBREL)
	switch m := pm.(type) {
	case *pubrelMessage:

	default:
		t.Errorf("msg should be *pubrelMessage, is %s", m)
	}
}

func Test_NewSubscribeMessage(t *testing.T) {
	sm := NewMessage(SUBSCRIBE)
	switch m := sm.(type) {
	case *subscribeMessage:

	default:
		t.Errorf("msg should be *subscribeMessage, is %s", m)
	}
}

func Test_NewSubackMessage(t *testing.T) {
	sm := NewMessage(SUBACK)
	switch m := sm.(type) {
	case *subackMessage:

	default:
		t.Errorf("msg should be *subackMessage, is %s", m)
	}
}

func Test_NewUnsubscribeMessage(t *testing.T) {
	sm := NewMessage(UNSUBSCRIBE)
	switch m := sm.(type) {
	case *unsubscribeMessage:

	default:
		t.Errorf("msg should be *unsubscribeMessage, is %s", m)
	}
}

func Test_NewUnsubackMessage(t *testing.T) {
	sm := NewMessage(UNSUBACK)
	switch m := sm.(type) {
	case *unsubackMessage:

	default:
		t.Errorf("msg should be *unsubackMessage, is %s", m)
	}
}

func Test_NewPingreqMessage(t *testing.T) {
	pm := NewMessage(PINGREQ)
	switch m := pm.(type) {
	case *pingreqMessage:

	default:
		t.Errorf("msg should be *pingreqMessage, is %s", m)
	}
}

func Test_NewPingrespMessage(t *testing.T) {
	pm := NewMessage(PINGRESP)
	switch m := pm.(type) {
	case *pingrespMessage:

	default:
		t.Errorf("msg should be *pingrespMessage, is %s", m)
	}
}

func Test_NewDisconnectMessage(t *testing.T) {
	dm := NewMessage(DISCONNECT)
	switch m := dm.(type) {
	case *disconnectMessage:

	default:
		t.Errorf("msg should be *disconnectMessage, is %s", m)
	}
}

func Test_NewWillTopicUpdateMessage(t *testing.T) {
	wm := NewMessage(WILLTOPICUPD)
	switch m := wm.(type) {
	case *willTopicUpdateMessage:

	default:
		t.Errorf("msg should be *willTopicUpdateMessage, is %s", m)
	}
}

func Test_NewWillTopicRespMessage(t *testing.T) {
	wm := NewMessage(WILLTOPICRESP)
	switch m := wm.(type) {
	case *willTopicRespMessage:

	default:
		t.Errorf("msg should be *willTopicRespMessage, is %s", m)
	}
}

func Test_NewWillMsgUpdateMessage(t *testing.T) {
	wm := NewMessage(WILLMSGUPD)
	switch m := wm.(type) {
	case *willMsgUpdateMessage:

	default:
		t.Errorf("msg should be *willMsgUpdateMessage, is %s", m)
	}
}

func Test_NewWillMsgRespMessage(t *testing.T) {
	wm := NewMessage(WILLMSGRESP)
	switch m := wm.(type) {
	case *willMsgRespMessage:

	default:
		t.Errorf("msg should be *willMsgRespMessage, is %s", m)
	}
}
