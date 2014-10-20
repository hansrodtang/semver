package semver

import "strings"

func containsOnly(s string, set string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(set, r)
	}) == -1
}

func hasLeadingZero(number string) bool {
	if len(number) > 1 {
		if strings.HasPrefix(number, "0") {
			return true
		}
	}
	return false
}

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

func rng(main, first, second *Version) bool {
	return gte(main, first) && lte(main, second)
}
