package parser

import (
	"bytes"
	"fmt"

	"github.com/hansrodtang/semver"
)

type node interface {
	Run(*semver.Version) bool
	String() string
}

type nodeContainer node

type nodeComparison struct {
	action comparatorFunc
	arg    *semver.Version
}

func (n nodeComparison) Run(main *semver.Version) bool {
	return n.action(main, n.arg)
}

func (n nodeComparison) String() string {
	nm := getFunctionName(n.action)
	for k, v := range comparators {
		if getFunctionName(v) == nm {
			return fmt.Sprintf("%v%v", k, n.arg)
		}
	}
	return ""
}

type nodeRange struct {
	sets []nodeSet
}

func (n nodeRange) Run(main *semver.Version) bool {
	for _, c := range n.sets {
		if c.Run(main) != false {
			return true
		}
	}
	return false
}

func (n nodeRange) String() string {
	var b bytes.Buffer
	for i, v := range n.sets {
		b.WriteString(v.String())
		if len(n.sets)-1 > i {
			b.WriteString(" || ")
		}
	}

	return b.String()
}

type nodeSet struct {
	comparisons []nodeComparison
}

func (n nodeSet) Run(main *semver.Version) bool {
	for _, c := range n.comparisons {
		if c.Run(main) != true {
			return false
		}
	}
	return true
}

func (n nodeSet) String() string {
	var b bytes.Buffer
	for i, v := range n.comparisons {
		b.WriteString(v.String())

		if len(n.comparisons)-1 > i {
			b.WriteString(" ")
		}
	}
	return b.String()
}

var comparators = map[string]comparatorFunc{
	string(operatorGT): gt,
	string(operatorGE): gte,
	string(operatorLT): lt,
	string(operatorLE): lte,
	string(operatorEQ): eq,
}
