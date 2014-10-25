package parser

import "github.com/hansrodtang/semver"

type Node interface {
	Run() bool
}

type NodeComparison struct {
	action comparatorFunc
	arg    *semver.Version
}

func (n NodeComparison) Run(main *semver.Version) bool {
	return n.action(main, n.arg)
}

type NodeRange struct {
	comparisons []NodeSet
}

func (n NodeRange) Run(main *semver.Version) bool {
	for _, c := range n.comparisons {
		if c.Run(main) != false {
			return true
		}
	}
	return false
}

type NodeSet struct {
	comparisons []NodeComparison
}

func (n NodeSet) Run(main *semver.Version) bool {
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
