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

var correctVersions = []string{
	"0.0.1-alpha.preview+123.456",
	"1.2.3-alpha.1+123.456",
	"1.2.3-alpha.1",
	"1.2.3+123.456",
	"1.2.3-alpha.b-eta+123.b-uild",
	"1.2.3+123.b-uild",
	"1.2.3-alpha.b-eta",
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
	{semver.Build(1, 0, 0, []string{"alpha", "1"}), semver.Build(1, 0, 0, []string{"alpha"}), 1},
	{semver.Build(1, 0, 0, []string{"alpha", "1"}), semver.Build(1, 0, 0, []string{"alpha", "beta"}), -1},
	{semver.Build(1, 0, 0, []string{"alpha", "beta"}), semver.Build(1, 0, 0, []string{"beta"}), -1},
	{semver.Build(1, 0, 0, []string{"beta"}), semver.Build(1, 0, 0, []string{"beta", "2"}), -1},
	{semver.Build(1, 0, 0, []string{"beta", "2"}), semver.Build(1, 0, 0, []string{"beta", "11"}), -1},
	{semver.Build(1, 0, 0, []string{"beta", "11"}), semver.Build(1, 0, 0, []string{"beta", "2"}), 1},
	{semver.Build(1, 0, 0, []string{"beta", "11"}), semver.Build(1, 0, 0, []string{"rc", "1"}), -1},
	{semver.Build(1, 0, 0, []string{"beta", "beta"}), semver.Build(1, 0, 0, []string{"beta", "alpha"}), 1},
	{semver.Build(1, 0, 0, []string{"rc", "1"}), semver.Build(1, 0, 0), -1},

	{semver.Build(1, 0, 0, []string{"beta", "alpha", "1"}), semver.Build(1, 0, 0, []string{"beta", "alpha"}), 1},

	{semver.Build(1, 0, 0, []string{"rc", "1"}), semver.Build(1, 0, 0, []string{"rc", "1"}, []string{"435345345"}), 0},
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
	expectedMajor := uint64(1)
	expectedMinor := uint64(2)
	expectedPatch := uint64(3)
	expectedPrerelease := "alpha.1"
	expectedMetadata := "35.45"

	ver := semver.Build(expectedMajor, expectedMinor, expectedPatch)
	ver.SetPrerelease("alpha", "1")
	ver.SetMetadata("35", "45")

	if result := ver.Major(); result != expectedMajor {
		t.Errorf("%q.Major() => %q, wanted %q", ver, result, expectedMajor)
	}
	if result := ver.Minor(); result != expectedMinor {
		t.Errorf("%q.Minor() => %q, wanted %q", ver, result, expectedMinor)
	}
	if result := ver.Patch(); result != expectedPatch {
		t.Errorf("%q.Patch() => %q, wanted %q", ver, result, expectedPatch)
	}
	if result := ver.Prerelease(); result != expectedPrerelease {
		t.Errorf("%q.Prerelease() => %q, wanted %q", ver, result, expectedPrerelease)
	}
	if result := ver.Metadata(); result != expectedMetadata {
		t.Errorf("%q.Metadata() => %q, wanted %q", ver, result, expectedMetadata)
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

func TestIncrementers(t *testing.T) {
	ver := semver.Build(1, 2, 3)
	ver.IncrementMajor()
	ver.IncrementMinor()
	ver.IncrementPatch()

	expected := "2.3.4"
	result := fmt.Sprint(ver)
	if result != expected {
		t.Errorf("Stringer() => %q, wanted %q", result, expected)
	}
}

func TestDecrementers(t *testing.T) {
	ver := semver.Build(1, 2, 3)
	ver.DecrementMajor()
	ver.DecrementMinor()
	ver.DecrementPatch()

	expected := "0.1.2"
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

func TestCorrectFormat(t *testing.T) {
	for _, version := range correctVersions {
		_, err := semver.New(version)
		if err != nil {
			t.Errorf("New(%q) => %v, want <nil>", version, err)
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

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		semver.New(VERSION)
	}
}

func BenchmarkParseComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		semver.New(VERSION)
	}
}

func BenchmarkParseAverage(b *testing.B) {
	l := len(correctVersions)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		semver.New(correctVersions[n%l])
	}
}

func BenchmarkCompareSimple(b *testing.B) {
	const VERSION = "0.0.1"
	v, _ := semver.New(VERSION)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Compare(v)
	}
}

func BenchmarkCompareComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	const VERSION2 = "0.0.1-alpha.preview+123.456"
	v, _ := semver.New(VERSION)
	v2, _ := semver.New(VERSION2)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Compare(v2)
	}
}

func BenchmarkCompareAverage(b *testing.B) {
	l := len(comparisons)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		comparisons[n%l].main.Compare(comparisons[n%l].other)
	}
}

func ExampleCompare() {
	v1, _ := semver.New("1.6.0")
	v2, _ := semver.New("1.5.0")
	// do something with error
	if v1.Compare(v2) > 0 {
		fmt.Println("v1 is larger")
	}
	// Output: v1 is larger
}
