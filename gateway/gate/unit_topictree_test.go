package gateway

import (
	"testing"
)

func eok(e error, t *testing.T) {
	if e != nil {
		t.Fatalf("ERROR %s\n", e)
	}
}

func enok(e error, t *testing.T) {
	if e == nil {
		t.Fatalf("ERROR (NO ERROR)")
	}
}

func chkb(b bool, e bool, t *testing.T) {
	if b != e {
		t.Fatalf("ERROR bool\n")
	}
}

func Test_NewTopicTree(t *testing.T) {
	tt := NewTopicTree()
	if tt == nil {
		t.Fatalf("NewTopicTree was nil")
	}
	if tt.root == nil {
		t.Fatalf("NewTopicTree.Root was nil")
	}
}

func Test_AddSubscription_alpha(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c1", conn, addr)
	tt := NewTopicTree()
	b, e := tt.AddSubscription(c, "alpha")
	eok(e, t)
	chkb(b, true, t)
	if len(tt.root.clients) != 0 {
		t.Fatalf("AddSub root had client")
	}
	if len(tt.root.children["alpha"].clients) != 1 {
		t.Fatalf("AddSub alpha != 1")
	}
}

func Test_AddSubscription_Salpha(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c2", conn, addr)
	tt := NewTopicTree()
	b, e := tt.AddSubscription(c, "/alpha")
	eok(e, t)
	chkb(b, true, t)
}

func Test_AddSubscription_alphaS(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c2.5", conn, addr)
	tt := NewTopicTree()
	b, e := tt.AddSubscription(c, "alpha/")
	enok(e, t)
	chkb(b, false, t)
}

func Test_AddSubscription_aSbScSd(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c3", conn, addr)
	tt := NewTopicTree()
	b, e := tt.AddSubscription(c, "a/b/c/d")
	eok(e, t)
	chkb(b, true, t)
}

func Test_AddSubscription_aSSb(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c4", conn, addr)
	tt := NewTopicTree()
	b, e := tt.AddSubscription(c, "a//b")
	enok(e, t)
	chkb(b, false, t)
}

func Test_AddSubscription_H(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c5", conn, addr)
	tt := NewTopicTree()
	b, e := tt.AddSubscription(c, "#")
	eok(e, t)
	chkb(b, true, t)
}

func Test_AddSubscription_SH(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c6", conn, addr)
	tt := NewTopicTree()
	b, e := tt.AddSubscription(c, "/#")
	eok(e, t)
	chkb(b, true, t)
}

func Test_AddSubscription_aSHSb(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c7", conn, addr)
	tt := NewTopicTree()
	b, e := tt.AddSubscription(c, "a/#/b")
	enok(e, t)
	chkb(b, false, t)
}

func Test_AddSubscription_aSHS(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c8", conn, addr)
	tt := NewTopicTree()
	b, e := tt.AddSubscription(c, "a/#/")
	enok(e, t)
	chkb(b, false, t)
}

func Test_AddSubscription_P(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c9", conn, addr)
	tt := NewTopicTree()
	b, e := tt.AddSubscription(c, "+")
	eok(e, t)
	chkb(b, true, t)
}

func Test_AddSubscription_SPSbSP(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c10", conn, addr)
	tt := NewTopicTree()
	b, e := tt.AddSubscription(c, "/+/b/+")
	eok(e, t)
	chkb(b, true, t)
}

func Test_AddSubscription(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c", conn, addr)
	tt := NewTopicTree()

	topics := []string{
		"a",
		"a/b",
		"a/b/c",
		"a/b/d",
		"a/b/e",
		"a/b/f",
		"a/b/c/d/e/f/g/h/i/j/k",
		"a/b/c/d/e/f/g/h/i/j/+",
		"a/b/c/d/e/f/g/h/i/j/k/+",
		"a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z",
		"a/b/c/d/e/f/g/h/i/j/k",
		"+/b/c",
		"a/+/b",
		"a/b/+",
		"a/b/#",
		"#",
	}

	for _, topic := range topics {
		_, e := tt.AddSubscription(c, topic)
		eok(e, t)
	}
}

func Benchmark_AddSubscription(b *testing.B) {
	var conn uConn
	var addr uAddr
	c := NewClient("b", conn, addr)
	tt := NewTopicTree()

	topics := []string{
		"a",
		"a/b",
		"a/b/c",
		"a/b/d",
		"a/b/e",
		"a/b/f",
		"a/b/c/d/e/f/g/h/i/j/k",
		"a/b/c/d/e/f/g/h/i/j/+",
		"a/b/c/d/e/f/g/h/i/j/k/+",
		"a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z",
		"a/b/c/d/e/f/g/h/i/j/k",
		"+/b/c",
		"a/+/b",
		"a/b/+",
		"a/b/#",
		"#",
	}

	for i := 0; i < b.N; i++ {
		for _, topic := range topics {
			tt.AddSubscription(c, topic)
		}
	}
}

func Test_SubscribersOf_none(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("so_A", conn, addr)
	tt := NewTopicTree()

	sum := 0
	sum += len(tt.SubscribersOf("+"))
	sum += len(tt.SubscribersOf("#"))
	sum += len(tt.SubscribersOf("a"))
	sum += len(tt.SubscribersOf("/+"))
	sum += len(tt.SubscribersOf("/#"))
	sum += len(tt.SubscribersOf("/a"))

	tt.AddSubscription(c, "kappa")
	sum += len(tt.SubscribersOf("+"))
	sum += len(tt.SubscribersOf("#"))
	sum += len(tt.SubscribersOf("a"))
	sum += len(tt.SubscribersOf("/+"))
	sum += len(tt.SubscribersOf("/#"))
	sum += len(tt.SubscribersOf("/a"))

	tt.AddSubscription(c, "/kappa")
	sum += len(tt.SubscribersOf("+"))
	sum += len(tt.SubscribersOf("#"))
	sum += len(tt.SubscribersOf("a"))
	sum += len(tt.SubscribersOf("/+"))
	sum += len(tt.SubscribersOf("/#"))
	sum += len(tt.SubscribersOf("/a"))

	tt.AddSubscription(c, "kappa/#")
	sum += len(tt.SubscribersOf("+"))
	sum += len(tt.SubscribersOf("#"))
	sum += len(tt.SubscribersOf("a"))
	sum += len(tt.SubscribersOf("/+"))
	sum += len(tt.SubscribersOf("/#"))
	sum += len(tt.SubscribersOf("/a"))

	tt.AddSubscription(c, "kappa/+")
	sum += len(tt.SubscribersOf("+"))
	sum += len(tt.SubscribersOf("#"))
	sum += len(tt.SubscribersOf("a"))
	sum += len(tt.SubscribersOf("/+"))
	sum += len(tt.SubscribersOf("/#"))
	sum += len(tt.SubscribersOf("/a"))

	tt.AddSubscription(c, "a/b")
	sum += len(tt.SubscribersOf("+"))
	sum += len(tt.SubscribersOf("#"))
	sum += len(tt.SubscribersOf("a"))
	sum += len(tt.SubscribersOf("/+"))
	sum += len(tt.SubscribersOf("/#"))
	sum += len(tt.SubscribersOf("/a"))

	tt.AddSubscription(c, "b/a")
	sum += len(tt.SubscribersOf("+"))
	sum += len(tt.SubscribersOf("#"))
	sum += len(tt.SubscribersOf("a"))
	sum += len(tt.SubscribersOf("/+"))
	sum += len(tt.SubscribersOf("/#"))
	sum += len(tt.SubscribersOf("/a"))

	tt.AddSubscription(c, "/a/#")
	sum += len(tt.SubscribersOf("+"))
	sum += len(tt.SubscribersOf("#"))
	sum += len(tt.SubscribersOf("a"))
	sum += len(tt.SubscribersOf("/+"))
	sum += len(tt.SubscribersOf("/#"))
	sum += len(tt.SubscribersOf("/a"))
	
	if sum != 0 {
		t.Fatalf("SubscribersOf_none had subscriber")
	}
}

func Test_SubscribersOf_1(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("so_A", conn, addr)
	tt := NewTopicTree()

	tt.AddSubscription(c, "a")
	subs := tt.SubscribersOf("a")
	if len(subs) != 1 {
		t.Fatalf("SubscribersOf_1 bad")
	}
}
