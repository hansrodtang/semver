package semver

import "strings"

type containsFunc func(rune) bool

func containsOnly(s string, c containsFunc) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !c(r)
	}) == -1
}

func alphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || (r == '-')
}

func hasLeadingZero(number string) bool {
	if len(number) > 1 {
		if strings.HasPrefix(number, "0") {
			return true
		}
	}
	return false
}

type comparatorFunc func(*Version, *Version) bool
type satisfactionMap map[*Version]comparatorFunc

func gt(main, other *Version) bool {
	return main.Compare(other) > 0
}

func gte(main, other *Version) bool {
	return main.Compare(other) >= 0
}

func lt(main, other *Version) bool {
	return main.Compare(other) < 0
}

func lte(main, other *Version) bool {
	return main.Compare(other) <= 0
}

func eq(main, other *Version) bool {
	return main.Compare(other) == 0
}

func rng2cpm(main, other *Version) satisfactionMap {
	return satisfactionMap{main: gt, other: lt}
}
