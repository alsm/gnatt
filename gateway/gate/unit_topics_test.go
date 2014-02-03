package gateway

/**  HEY YOU! **/
/* Remember benchmarks can be run with the command
 *  go test -bench=".*" ./...
 *
 * I get that a lot of these algorithms have more elegant
 * recursive solutions, but all I care about is performance.
 */

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

// 450 ns/op
func Benchmark_ContainsWildcard(b *testing.B) {
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
		"a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z":   false,
		"a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/+/z":   true,
		"a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/#":   true,
		"+/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z":   true,
		"a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/+/s/+/u/v/w/x/+/z":   true,
		"a/b/c/d/e/f/g/+/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/+/z":   true,
		"a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/+/y/+":   true,
		"a/b/+/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/+/z":   true,
		"a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/+/w/x/+/z":   true,
		"+/+/+/d/+/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/+/z":   true,
		"a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/+/t/u/v/w/x/+/#":   true,
		"a/b/c/d/e/f/g/h/i/j/k/+/+/n/+/p/+/r/+/t/u/v/w/x/+/z":   true,
		"a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/++/##/z": false,
	}

	res := false
	for i := 0; i < b.N; i++ {
		for topic, _ := range topics {
			res = res && ContainsWildcard(topic)
		}
	}
}

func Test_ValidateSubscribeTopicName(t *testing.T) {
	topics := map[string]bool{
		"":         false,
		"a":        true,
		"#":        true,
		"+":        true,
		"/":        false,
		"/a":       true,
		"/#":       true,
		"a/":       false,
		"+/":       false,
		"#/":       false,
		"a/b":      true,
		"a/#":      true,
		"a/+":      true,
		"b/a":      true,
		"#/b":      false,
		"+/b":      true,
		"a/b/c/":   false,
		"/a/b/c":   true,
		"/a/++":    true,
		"a/++/b":   true,
		"/a/##":    true,
		"a/##/b":   true,
		"a/#/b":    false,
		"a/+/b":    true,
		"a/+/b/+":  true,
		"/a/+/b/#": true,
		"/a/+#/b":  true,
		"/a/#+/b":  true,
	}

	for topic, expectError := range topics {
		_, e := ValidateSubscribeTopicName(topic)
		if expectError == (e != nil) {
			t.Errorf("error validating topic \"%s\"", topic)
		}
	}
}

// 5833 ns/op
func Benchmark_ValidateSubscribeTopicName(b *testing.B) {
	topics := map[string]bool{
		"":         false,
		"a":        true,
		"#":        true,
		"+":        true,
		"/":        false,
		"/a":       true,
		"/#":       true,
		"a/":       false,
		"+/":       false,
		"#/":       false,
		"a/b":      true,
		"a/#":      true,
		"a/+":      true,
		"b/a":      true,
		"#/b":      false,
		"+/b":      true,
		"a/b/c/":   false,
		"/a/b/c":   true,
		"/a/++":    true,
		"a/++/b":   true,
		"/a/##":    true,
		"a/##/b":   true,
		"a/#/b":    false,
		"a/+/b":    true,
		"a/+/b/+":  true,
		"/a/+/b/#": true,
		"/a/+#/b":  true,
		"/a/#+/b":  true,
	}

	for i := 0; i < b.N; i++ {
		sum := 0
		for topic, _ := range topics {
			toks, _ := ValidateSubscribeTopicName(topic)
			sum += len(toks)
		}
	}
}

func Test_ValidatePublishTopicName(t *testing.T) {
	topics := map[string]bool{
		"":         false,
		"a":        true,
		"#":        false,
		"+":        false,
		"/":        false,
		"/a":       true,
		"/#":       false,
		"a/":       false,
		"+/":       false,
		"#/":       false,
		"a/b":      true,
		"a/#":      false,
		"a/+":      false,
		"b/a":      true,
		"#/b":      false,
		"+/b":      false,
		"a/b/c/":   false,
		"/a/b/c":   true,
		"/a/++":    true,
		"a/++/b":   true,
		"/a/##":    true,
		"a/##/b":   true,
		"a/#/b":    false,
		"a/+/b":    false,
		"a/+/b/+":  false,
		"/a/+/b/#": false,
		"/a/+#/b":  true,
		"/a/#+/b":  true,
	}

	for topic, expectError := range topics {
		_, e := ValidatePublishTopicName(topic)
		if expectError == (e != nil) {
			t.Errorf("error validating topic \"%s\"", topic)
		}
	}
}

// 8299 ns/op
func Benchmark_ValidatePublishTopicName(b *testing.B) {
	topics := map[string]bool{
		"":         false,
		"a":        true,
		"#":        false,
		"+":        false,
		"/":        false,
		"/a":       true,
		"/#":       false,
		"a/":       false,
		"+/":       false,
		"#/":       false,
		"a/b":      true,
		"a/#":      false,
		"a/+":      false,
		"b/a":      true,
		"#/b":      false,
		"+/b":      false,
		"a/b/c/":   false,
		"/a/b/c":   true,
		"/a/++":    true,
		"a/++/b":   true,
		"/a/##":    true,
		"a/##/b":   true,
		"a/#/b":    false,
		"a/+/b":    false,
		"a/+/b/+":  false,
		"/a/+/b/#": false,
		"/a/+#/b":  true,
		"/a/#+/b":  true,
	}

	for i := 0; i < b.N; i++ {
		sum := 0
		for topic, _ := range topics {
			toks, _ := ValidatePublishTopicName(topic)
			sum += len(toks)
		}
	}
}
