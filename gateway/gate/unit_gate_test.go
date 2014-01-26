package gateway

import (
	"sync"
	"testing"
)

func new_topicNames() *topicNames {
	t := &topicNames{
		sync.RWMutex{},
		make(map[uint16]string),
		0,
	}
	return t
}

func Test_topicName_contains(t *testing.T) {
	topics := new_topicNames()
	topics.containsTopic("notinthere")

	if topics.containsTopic("notinthere") != false {
		t.Errorf("new topicNames contains something")
	}

	topics.putTopic("alpha")
	if topics.containsTopic("beta") != false {
		t.Errorf("topicNames contains invalid topic")
	}

	if topics.containsTopic("alpha") == false {
		t.Errorf("topicNames missing topic")
	}

	if topics.containsId(333) {
		t.Errorf("topicNames contains rediculous topicId")
	}

	if !topics.containsId(1) {
		t.Errorf("topicNames is missing a topicId")
	}

	if topics.containsId(0) {
		t.Errorf("topicNames contains impossible topicId")
	}
}

func Test_topicName_putTopic(t *testing.T) {
	topics := new_topicNames()

	i := topics.putTopic("foo")
	if !topics.containsTopic("foo") {
		t.Errorf("topicNames missing topic")
	}
	if i != 1 {
		t.Errorf("topicNames put unexpected topicId")
	}

	i = topics.putTopic("bar")
	if !topics.containsTopic("bar") {
		t.Errorf("topicNames missing 2nd topic")
	}
	if i != 2 {
		t.Errorf("topicNames put unexpected topicId")
	}

	if !topics.containsTopic("foo") {
		t.Errorf("topicNames lost topic")
	}
}

func Test_topicName_get(t *testing.T) {
	topics := new_topicNames()

	a := topics.putTopic("/a/b")
	b := topics.getId("/a/b")
	if a != b {
		t.Errorf("topicNames did not get the same topic id as assigned by put")
	}
	b2 := topics.getTopic(1)
	if b2 != "/a/b" {
		t.Errorf("getTopic got wrong topic")
	}

	c := topics.putTopic("d/c")
	d := topics.getId("d/c")
	e := topics.getId("/a/b")
	if c != d {
		t.Errorf("topicNames did not get the same topic id as assigned by put")
	}
	if c == e {
		t.Errorf("topicName assigned same topic id to different topics")
	}
}
