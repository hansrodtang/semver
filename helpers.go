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

func numbers(r rune) bool {
	return (r >= '0' && r <= '9')
}

func hasLeadingZero(number string) bool {
	if len(number) > 1 {
		if strings.HasPrefix(number, "0") {
			return true
		}
	}
	return false
}
