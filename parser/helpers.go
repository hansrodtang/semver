package parser

import (
	"reflect"
	"runtime"
	"strconv"
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

func hy2op(v1, v2 *semver.Version) node {
	return nodeSet{
		nodeComparison{gte, v1},
		nodeComparison{lte, v2},
	}
}

func cr2op(i item) node {
	return nil
}

func tld2op(i item) node {
	if i.typ == itemXRange {
		return xr2op(i)
	}
	v1, _ := semver.New(i.val)
	v2 := *v1
	v2.IncrementMinor()
	v2.SetPatch(0)
	return nodeSet{
		nodeComparison{gte, v1},
		nodeComparison{lt, &v2},
	}
}

func xr2op(i item) node {

	version := i.val
	s := strings.Split(version, dot)

	for len(s) < 3 {
		s = append(s, "*")
	}
	major, err1 := strconv.ParseUint(s[0], 10, 0)
	minor, err2 := strconv.ParseUint(s[1], 10, 0)
	patch, _ := strconv.ParseUint(s[2], 10, 0)

	v1 := semver.Build(major, minor, patch)

	v2 := *v1

	if err1 != nil {
		v2.SetMinor(0)
		v2.SetPatch(0)
		return nodeSet{
			nodeComparison{gte, &v2},
		}
	}
	if err2 != nil {
		v2.IncrementMajor()
		v2.SetMinor(0)
		v2.SetPatch(0)
		return nodeSet{
			nodeComparison{gte, v1},
			nodeComparison{lt, &v2},
		}
	}

	v2.IncrementMinor()
	return nodeSet{
		nodeComparison{gte, v1},
		nodeComparison{lt, &v2},
	}

}

func getFunctionName(i interface{}) string {
	fname := strings.Split(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name(), ".")
	return fname[len(fname)-1]
}
