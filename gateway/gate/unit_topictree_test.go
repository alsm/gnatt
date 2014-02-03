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
	e := tt.AddSubscription(c, "alpha")
	eok(e, t)
}

func Test_AddSubscription_Salpha(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c2", conn, addr)
	tt := NewTopicTree()
	e := tt.AddSubscription(c, "/alpha")
	eok(e, t)
}

func Test_AddSubscription_alphaS(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c2.5", conn, addr)
	tt := NewTopicTree()
	e := tt.AddSubscription(c, "alpha/")
	enok(e, t)
}

func Test_AddSubscription_aSbScSd(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c3", conn, addr)
	tt := NewTopicTree()
	e := tt.AddSubscription(c, "a/b/c/d")
	eok(e, t)
}

func Test_AddSubscription_aSSb(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c4", conn, addr)
	tt := NewTopicTree()
	e := tt.AddSubscription(c, "a//b")
	enok(e, t)
}

func Test_AddSubscription_H(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c5", conn, addr)
	tt := NewTopicTree()
	e := tt.AddSubscription(c, "#")
	eok(e, t)
}

func Test_AddSubscription_SH(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c6", conn, addr)
	tt := NewTopicTree()
	e := tt.AddSubscription(c, "/#")
	eok(e, t)
}

func Test_AddSubscription_aSHSb(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c7", conn, addr)
	tt := NewTopicTree()
	e := tt.AddSubscription(c, "a/#/b")
	enok(e, t)
}

func Test_AddSubscription_aSHS(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c8", conn, addr)
	tt := NewTopicTree()
	e := tt.AddSubscription(c, "a/#/")
	enok(e, t)
}

func Test_AddSubscription_P(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c9", conn, addr)
	tt := NewTopicTree()
	e := tt.AddSubscription(c, "+")
	eok(e, t)
}

func Test_AddSubscription_SPSbSP(t *testing.T) {
	var conn uConn
	var addr uAddr
	c := NewClient("c10", conn, addr)
	tt := NewTopicTree()
	e := tt.AddSubscription(c, "/+/b/+")
	eok(e, t)
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
		e := tt.AddSubscription(c, topic)
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
