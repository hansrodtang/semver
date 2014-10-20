package semver_test

import (
	"fmt"
	"testing"

	"github.com/hansrodtang/semver"
)

type comparison struct {
	main     *semver.Version
	other    *semver.Version
	expected int
}

var badVersions = []string{
	"",
	".",
	"1.",
	".1",
	"a.b.c",
	"1.a.b",
	"1.1.a",
	"1.a.1",
	"a.1.1",
	"..",
	"1..",
	"1.1.",
	"1..1",
	"1.1.+123",
	"1.1.-beta",
	"-1.1.1",
	"1.-1.1",
	"1.1.-1",
	// Leading zeroes
	"01.1.1",
	"001.1.1",
	"1.01.1",
	"1.001.1",
	"1.1.01",
	"1.1.001",
	"1.1.1-01",
	"1.1.1-001",
	"1.1.1-beta.01",
	"1.1.1-beta.001",
	"0.0.0-!",
	"0.0.0+!",
	// empty prerelease
	"0.0.0-.alpha",
	// empty build metadata
	"0.0.0-alpha+",
	"0.0.0-alpha+test.",
}

var comparisons = []comparison{
	{semver.Build(1, 0, 0), semver.Build(1, 0, 0), 0},
	{semver.Build(2, 0, 0), semver.Build(1, 0, 0), 1},
	{semver.Build(0, 1, 0), semver.Build(0, 1, 0), 0},
	{semver.Build(0, 2, 0), semver.Build(0, 1, 0), 1},
	{semver.Build(0, 0, 1), semver.Build(0, 0, 1), 0},
	{semver.Build(0, 0, 2), semver.Build(0, 0, 1), 1},
	{semver.Build(1, 2, 3), semver.Build(1, 2, 3), 0},
	{semver.Build(2, 2, 4), semver.Build(1, 2, 4), 1},
	{semver.Build(1, 3, 3), semver.Build(1, 2, 3), 1},
	{semver.Build(1, 2, 4), semver.Build(1, 2, 3), 1},

	// Spec Examples #11
	{semver.Build(1, 0, 0), semver.Build(2, 0, 0), -1},
	{semver.Build(2, 0, 0), semver.Build(2, 1, 0), -1},
	{semver.Build(2, 1, 0), semver.Build(2, 1, 1), -1},

	// Spec Examples #9
	{semver.Build(1, 0, 0), semver.Build(1, 0, 0, []string{"alpha"}), 1},
	{semver.Build(1, 0, 0, []string{"alpha"}), semver.Build(1, 0, 0, []string{"alpha", "1"}), -1},
	{semver.Build(1, 0, 0, []string{"alpha", "1"}), semver.Build(1, 0, 0, []string{"alpha", "beta"}), -1},
	{semver.Build(1, 0, 0, []string{"alpha", "beta"}), semver.Build(1, 0, 0, []string{"beta"}), -1},
	{semver.Build(1, 0, 0, []string{"beta"}), semver.Build(1, 0, 0, []string{"beta", "2"}), -1},
	{semver.Build(1, 0, 0, []string{"beta", "2"}), semver.Build(1, 0, 0, []string{"beta", "11"}), -1},
	{semver.Build(1, 0, 0, []string{"beta", "11"}), semver.Build(1, 0, 0, []string{"beta", "2"}), 1},
	{semver.Build(1, 0, 0, []string{"beta", "11"}), semver.Build(1, 0, 0, []string{"rc", "1"}), -1},
	{semver.Build(1, 0, 0, []string{"beta", "beta"}), semver.Build(1, 0, 0, []string{"beta", "alpha"}), 1},
	{semver.Build(1, 0, 0, []string{"rc", "1"}), semver.Build(1, 0, 0), -1},
}

var badPreRelease = [][]string{
	{"alpha", "-3"},
}

func TestStringer(t *testing.T) {
	ver := semver.Build(1, 2, 3)
	ver.SetPrerelease("alpha", "1")
	ver.SetMetadata("35", "45")

	expected := "1.2.3-alpha.1+35.45"
	result := fmt.Sprint(ver)

	if result != expected {
		t.Errorf("String() => %q, wanted %q", result, expected)
	}
}

func TestGetters(t *testing.T) {
	expected_major := uint64(1)
	expected_minor := uint64(2)
	expected_patch := uint64(3)
	expected_prerelease := "alpha.1"
	expected_metadata := "35.45"

	ver := semver.Build(expected_major, expected_minor, expected_patch)
	ver.SetPrerelease("alpha", "1")
	ver.SetMetadata("35", "45")

	if result := ver.Major(); result != expected_major {
		t.Errorf("%q.Major() => %q, wanted %q", ver, result, expected_major)
	}
	if result := ver.Minor(); result != expected_minor {
		t.Errorf("%q.Minor() => %q, wanted %q", ver, result, expected_minor)
	}
	if result := ver.Patch(); result != expected_patch {
		t.Errorf("%q.Patch() => %q, wanted %q", ver, result, expected_patch)
	}
	if result := ver.Prerelease(); result != expected_prerelease {
		t.Errorf("%q.Prerelease() => %q, wanted %q", ver, result, expected_prerelease)
	}
	if result := ver.Metadata(); result != expected_metadata {
		t.Errorf("%q.Metadata() => %q, wanted %q", ver, result, expected_metadata)
	}
}

func TestSetters(t *testing.T) {
	ver := semver.Build(1, 1, 1)
	ver.SetMajor(2)
	ver.SetMinor(3)
	ver.SetPatch(4)
	if err := ver.SetPrerelease("beta", "1"); err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	if err := ver.SetMetadata("22", "43"); err != nil {
		t.Errorf(fmt.Sprint(err))
	}

	expected := "2.3.4-beta.1+22.43"
	result := fmt.Sprint(ver)
	if result != expected {
		t.Errorf("Stringer() => %q, wanted %q", result, expected)
	}
}

func TestNew(t *testing.T) {
	expected := "1.0.3-alpha.1+35.45"

	ver, err := semver.New(expected)
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}

	result := fmt.Sprint(ver)
	if result != expected {
		t.Errorf("Stringer() => %q, want %q", result, expected)
	}

}

func TestBadFormat(t *testing.T) {
	for _, version := range badVersions {
		_, err := semver.New(version)
		if err == nil {
			// TODO: Set up error types
			t.Errorf("New(%q) => %v, want Error", version, err)
		}
	}
}

func TestComparison(t *testing.T) {
	for _, c := range comparisons {
		result := c.main.Compare(c.other)
		if result != c.expected {
			t.Errorf("%q.Compare(%q) => %#v, want %#v", c.main, c.other, result, c.expected)
		}
	}
}

func TestSetPrerelease(t *testing.T) {
	ver := semver.Build(1, 2, 3)
	err := ver.SetPrerelease("12", "1", "0", "43")
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
}

func BenchmarkParseSimple(b *testing.B) {
	const VERSION = "0.0.1"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		semver.New(VERSION)
	}
}

func BenchmarkParseComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		semver.New(VERSION)
	}
}

func BenchmarkCompareSimple(b *testing.B) {
	const VERSION = "0.0.1"
	v, _ := semver.New(VERSION)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Compare(v)
	}
}

func BenchmarkCompareComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	v, _ := semver.New(VERSION)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Compare(v)
	}
}
