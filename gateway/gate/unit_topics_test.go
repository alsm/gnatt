package gateway

import (
	"testing"
)

func Test_ContainsWildcard(t *testing.T) {
	topics := map[string]bool{
		"a":       false,
		"a/b":     false,
		"/a/b":    false,
		"+":       true,
		"/+":      true,
		"a/+":     true,
		"a/b/c/+": true,
		"a/+/+/d": true,
		"a/++/+d": false,
		"#":       true,
		"/#":      true,
		"a/b/c/#": true,
		"a/##/b":  false,
	}

	for topic, exp := range topics {
		res := ContainsWildcard(topic)
		if res != exp {
			t.Errorf("ContainsWildcard expected \"%s\", got \"%s\"", exp, res)
		}
	}
}
