package parser

import (
	"testing"

	"github.com/hansrodtang/semver"
)

type test struct {
	expected bool
	version  *semver.Version
}

var parsables = map[string][]test{
	"1.2.7 || >=1.2.9 <2.0.0": {
		{true, semver.Build(1, 2, 7)},
		{true, semver.Build(1, 2, 9)},
		{true, semver.Build(1, 4, 6)},
		{false, semver.Build(1, 2, 8)},
		{false, semver.Build(2, 0, 0)},
	},
	//"1.2 <1.2.9 || >2.0.0": {
	//	{false, semver.Build(1,2,10)},
	//  {true, semver.Build(1,2,8)}
	//}
}

func TestParser(t *testing.T) {

	for k, v := range parsables {
		n, _ := Parse(k)
		for _, x := range v {
			if response := n.Run(x.version); response != x.expected {
				t.Errorf("%q.Run(%q) => %t, want %t", k, x.version, response, x.expected)
			}
		}
	}
}

func BenchmarkParser(b *testing.B) {
	const VERSION = "1.2.7 || >=1.2.9 <2.0.0"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Parse(VERSION)
	}
}

func BenchmarkRunner(b *testing.B) {
	const VERSION = "1.2.7 || >=1.2.9 <2.0.0"
	p, _ := Parse(VERSION)
	v := semver.Build(2, 0, 0)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		p.Run(v)
	}
}
