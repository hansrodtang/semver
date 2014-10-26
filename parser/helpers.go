package parser

import "github.com/hansrodtang/semver"

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
