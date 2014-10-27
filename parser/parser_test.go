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
	"1.2 <1.2.9 || >2.0.0": {
		{false, semver.Build(1, 2, 10)},
		{false, semver.Build(1, 5, 1)},
		{true, semver.Build(1, 2, 8)},
		{true, semver.Build(1, 2, 7)},
	},
	"* || >2.0.0": {
		{true, semver.Build(1, 0, 10)},
		{true, semver.Build(100, 5, 1)},
		{true, semver.Build(1, 100, 8)},
		{true, semver.Build(1, 2, 100)},
	},
	"* >2.0.0": {
		{false, semver.Build(1, 0, 10)},
		{true, semver.Build(100, 5, 1)},
		{false, semver.Build(1, 100, 8)},
		{false, semver.Build(1, 2, 100)},
	},
	"1.0.0 - 2.0.0": {
		{false, semver.Build(0, 0, 10)},
		{false, semver.Build(3, 5, 1)},
		{true, semver.Build(1, 1, 5)},
		{true, semver.Build(1, 9, 7)},
	},
	"~1.2.3": {
		{false, semver.Build(1, 3, 2)},
		{false, semver.Build(1, 2, 2)},
		{true, semver.Build(1, 2, 5)},
		{true, semver.Build(1, 2, 9)},
	},
	"~1.2": {
		{false, semver.Build(1, 3, 2)},
		{false, semver.Build(1, 1, 9)},
		{true, semver.Build(1, 2, 3)},
		{true, semver.Build(1, 2, 9)},
	},
}

func TestParser(t *testing.T) {

	for k, v := range parsables {
		n, err := Parse(k)
		if err != nil {
			t.Error(err)
		} else {
			for _, x := range v {
				if response := n.Run(x.version); response != x.expected {
					t.Errorf("%q.Run(%q) => %t, want %t", k, x.version, response, x.expected)
				}
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
