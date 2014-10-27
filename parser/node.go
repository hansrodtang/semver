package parser

import (
	"bytes"
	"fmt"

	"github.com/hansrodtang/semver"
)

type nodeType int

const (
	errorNode nodeType = iota
	rangeNode
	comparisonNode
	setNode
)

type node interface {
	Run(*semver.Version) bool
	String() string
	Type() nodeType
}

type nodeContainer node

type nodeError struct {
	error item
}

func (n nodeError) Run(main *semver.Version) bool {
	return false
}

func (n nodeError) String() string {
	return n.error.String()
}

func (n nodeError) Type() nodeType {
	return errorNode
}

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

func (n nodeComparison) Type() nodeType {
	return comparisonNode
}

type nodeRange struct {
	sets []node
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

func (n nodeRange) Type() nodeType {
	return rangeNode
}

type nodeSet []node

func (n nodeSet) Run(main *semver.Version) bool {
	for _, c := range n {
		if c.Run(main) != true {
			return false
		}
	}
	return true
}

func (n nodeSet) String() string {
	var b bytes.Buffer
	for i, v := range n {
		b.WriteString(v.String())

		if len(n)-1 > i {
			b.WriteString(" ")
		}
	}
	return b.String()
}

func (n nodeSet) Type() nodeType {
	return setNode
}

var comparators = map[string]comparatorFunc{
	string(operatorGT): gt,
	string(operatorGE): gte,
	string(operatorLT): lt,
	string(operatorLE): lte,
	string(operatorEQ): eq,
}
