package semver

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var set = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789"

var containsMatches = []string{
	"abce345656",
	"controller",
	"matches",
	"12345678",
}

var containsMismatches = []string{
	"this is not right",
	"digitalmagasinet.no",
	"@twitter",
}

var zeroesMismatches = []string{
	"12434",
	"5000",
	"56032",
	"12345678",
	"0",
}

var zeroesMatches = []string{
	"0003",
	"02",
	"034",
}

func getFunctionName(i interface{}) string {
	fname := strings.Split(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name(), ".")
	return fname[len(fname)-1]
}

func TestContainsOnly(t *testing.T) {
	// Check matches
	matchesExpected := true
	for _, x := range containsMatches {
		if response := containsOnly(x, alphanumeric); response != matchesExpected {
			t.Errorf("containsOnly(%q, %q) => %t, want %t", x, set, response, matchesExpected)
		}
	}
	// Check mismatches
	mismatchesExpected := false
	for _, x := range containsMismatches {
		if response := containsOnly(x, alphanumeric); response != mismatchesExpected {
			t.Errorf("containsOnly(%q, %q) => %t, want %t", x, set, response, mismatchesExpected)
		}
	}
}

func TestLeadingZeroes(t *testing.T) {
	matchesExpected := true
	for _, x := range zeroesMatches {
		if response := hasLeadingZero(x); response != matchesExpected {
			t.Errorf("hasLeadingZero(%q): => %t, want %t", x, response, matchesExpected)
		}
	}
	// Check mismatches
	mismatchesExpected := false
	for _, x := range zeroesMismatches {
		if response := hasLeadingZero(x); response != mismatchesExpected {
			t.Errorf("hasLeadingZero(%q): => %t, want %t", x, response, mismatchesExpected)
		}
	}
}

func TestComparators(t *testing.T) {
	ver1 := Build(1, 2, 4)
	ver2 := Build(3, 2, 1)
	ver3 := Build(2, 1, 1)

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
		Build(1, 2, 0): gte,
		Build(3, 3, 1): lte,
		Build(3, 2, 1): eq,
	}
	ver := Build(3, 2, 1)

	expected := true
	for v, f := range v {
		response := f(ver, v)
		if !response {
			t.Errorf("%v(%q,%q): => %t, want %t", getFunctionName(f), ver, v, response, expected)
		}
	}
}
