package semver

import (
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

func TestContainsOnly(t *testing.T) {
	// Check matches
	matchesExpected := true
	for _, x := range containsMatches {
		if response := containsOnly(x, set); response != matchesExpected {
			t.Errorf("containsOnly(%q, %q) => %t, want %t", x, set, response, matchesExpected)
		}
	}
	// Check mismatches
	mismatchesExpected := false
	for _, x := range containsMismatches {
		if response := containsOnly(x, set); response != mismatchesExpected {
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
