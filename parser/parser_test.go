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
