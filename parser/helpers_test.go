package parser

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/hansrodtang/semver"
)

func getFunctionName(i interface{}) string {
	fname := strings.Split(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name(), ".")
	return fname[len(fname)-1]
}

func TestComparators(t *testing.T) {
	ver1 := Build(1, 2, 4)
	ver2 := Build(3, 2, 1)

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
