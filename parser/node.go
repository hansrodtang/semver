package parser

import "github.com/hansrodtang/semver"

type node interface {
	Run() bool
}

type nodeComparison struct {
	action comparatorFunc
	arg    *semver.Version
}

func (n nodeComparison) Run(main *semver.Version) bool {
	return n.action(main, n.arg)
}

type nodeRange struct {
	comparisons []nodeSet
}

func (n nodeRange) Run(main *semver.Version) bool {
	for _, c := range n.comparisons {
		if c.Run(main) != false {
			return true
		}
	}
	return false
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

var comparators = map[string]comparatorFunc{
	string(operatorGT): gt,
	string(operatorGE): gte,
	string(operatorLT): lt,
	string(operatorLE): lte,
	string(operatorEQ): eq,
}
