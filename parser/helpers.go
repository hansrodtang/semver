package parser

import (
	"fmt"
	"strings"

	"github.com/hansrodtang/semver"
)

type comparatorFunc func(*semver.Version, *semver.Version) bool
type satisfactionMap map[*semver.Version]comparatorFunc

func gt(main, other *semver.Version) bool {
	return main.Compare(other) > 0
}

func gte(main, other *semver.Version) bool {
	return main.Compare(other) >= 0
}

func lt(main, other *semver.Version) bool {
	return main.Compare(other) < 0
}

func lte(main, other *semver.Version) bool {
	return main.Compare(other) <= 0
}

func eq(main, other *semver.Version) bool {
	return main.Compare(other) == 0
}

func hy2op(v1, v2 *semver.Version) []nodeComparison {
	return []nodeComparison{
		{gte, v1},
		{lte, v2},
	}
}

func cr2op(i item) []nodeComparison {
	return nil
}

func tld2op(i item) []nodeComparison {
	return nil
}

func xr2op(i item) []nodeComparison {
	r := strings.NewReplacer("x", "*", "X", "*")

	var version string

	version = i.val

	if strings.Contains(version, plus) {
		output := strings.Split(version, plus)
		version = output[0]
	}

	if strings.Contains(version, hyphen) {
		output := strings.SplitN(version, hyphen, 2)
		version = output[0]
	}

	version = r.Replace(version)
	s := strings.SplitN(version, dot, 3)

	for len(s) < 3 {
		s = append(s, "*")
	}

	version = strings.Join(s, dot)
	version = strings.Replace(version, "*", "0", 3)

	v1, err := semver.New(version)
	if err != nil {
		fmt.Println(err)
	}
	v2 := *v1

	if s[0] == "*" {
		return []nodeComparison{
			{gte, v1},
		}
	}
	if s[1] == "*" {
		v2.IncrementMajor()
		return []nodeComparison{
			{gte, v1},
			{lt, &v2},
		}
	}

	v2.IncrementMinor()
	return []nodeComparison{
		{gte, v1},
		{lt, &v2},
	}

}
