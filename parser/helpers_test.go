package parser

import (
	"fmt"
	"testing"

	"github.com/hansrodtang/semver"
)

func TestComparators(t *testing.T) {
	ver1 := semver.Build(1, 2, 4)
	ver2 := semver.Build(3, 2, 1)

	expected := false
	if response := gt(ver1, ver2); response != expected {
		t.Errorf("gt(%q, %q): => %t, want %t", ver1, ver2, response, expected)
	}

	expected = true
	if response := gte(ver1, ver1); response != expected {
		t.Errorf("gte(%q, %q): => %t, want %t", ver1, ver2, response, expected)
	}

	expected = true
	if response := lt(ver1, ver2); response != expected {
		t.Errorf("lt(%q, %q): => %t, want %t", ver1, ver2, response, expected)
	}

	expected = true
	if response := lte(ver2, ver2); response != expected {
		t.Errorf("lte(%q, %q): => %t, want %t", ver1, ver2, response, expected)
	}

	expected = false
	if response := eq(ver1, ver2); response != expected {
		t.Errorf("eq(%q, %q): => %t, want %t", ver1, ver2, response, expected)
	}
}

func TestComparatorFunc(t *testing.T) {

	v := satisfactionMap{
		semver.Build(1, 2, 0): gte,
		semver.Build(3, 3, 1): lte,
		semver.Build(3, 2, 1): eq,
	}
	ver := semver.Build(3, 2, 1)

	expected := true
	for v, f := range v {
		response := f(ver, v)
		if !response {
			t.Errorf("%v(%q,%q): => %t, want %t", getFunctionName(f), ver, v, response, expected)
		}
	}
}

func TestXRangesConverter(t *testing.T) {
	i := item{itemAdvanced, "1.x"}
	c := xr2op(i)
	for _, v := range c {

		fmt.Println(getFunctionName(v.action), v.arg)
	}
}
