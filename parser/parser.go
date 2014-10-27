package parser

import (
	"errors"

	"github.com/hansrodtang/semver"
)

type parser struct {
	l      *lexer
	result node
	ibuf   []item
	pos    int
}

func (p *parser) run() (node, error) {

	n := handleRange(p)
	if n.Type() != errorNode {
		return n, nil
	}
	return nil, errors.New(n.String())

}

func (p *parser) next() item {
	if p.pos >= len(p.ibuf) {
		i := p.l.nextItem()
		p.ibuf = append(p.ibuf, i)
		p.pos++
		return i
	}
	i := p.ibuf[p.pos]
	p.pos++
	return i
}

func (p *parser) backup() {
	p.pos--
}

func Parse(input string) (node, error) {
	l := lex(input)
	p := &parser{l, nil, []item{}, 0}
	return p.run()

}

func handleOperator(p *parser) node {
	var set nodeSet
	for {
		i := p.next()

		switch i.typ {
		case itemVersion:
			ver1, _ := semver.New(i.val)
			if i = p.next(); i.typ == itemAdvanced {
				if i.val == string(operatorHY) {
					i = p.next()
					ver2, _ := semver.New(i.val)
					nc := hy2op(ver1, ver2)
					set = append(set, nc)
					return set
				}
			}
			p.backup()
			nc := nodeSet{nodeComparison{eq, ver1}}
			set = append(set, nc)
			return set
		case itemAdvanced:
			if i.val == string(operatorTR) {
				i := p.next()
				if i.typ != itemError {
					nc := tld2op(i)
					set = append(set, nc)
					return set
				}

			}
		case itemXRange:
			nc := xr2op(i)
			set = append(set, nc)
			return set
		default:
			v := p.next()
			ver, _ := semver.New(v.val)
			nc := nodeSet{nodeComparison{comparators[i.val], ver}}
			set = append(set, nc...)
			return set
		}

	}
}

func handleSet(p *parser) nodeSet {
	var set nodeSet

	for {
		i := p.next()

		switch i.typ {
		case itemSet:
			break
		case itemEOF:
			p.backup()
			return set
		case itemRange:
			return set
		default:
			p.backup()
			nc := handleOperator(p)
			set = append(set, nc)

		}
	}
}

func handleRange(p *parser) node {
	var ns nodeSet
	var rng nodeRange

	for {
		i := p.next()
		switch i.typ {
		case itemError:
			return nodeError{i}
		case itemEOF:
			return rng
		default:
			p.backup()
			ns = handleSet(p)
			rng.sets = append(rng.sets, ns)

		}
	}
}
