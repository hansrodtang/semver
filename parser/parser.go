package parser

import "github.com/hansrodtang/semver"

type parser struct {
	items  chan item // channel of scanned items.
	result node
	ibuf   []item
	pos    int
}

func (p *parser) run() (node, error) {

	return handleRange(p), nil

}

func (p *parser) next() item {
	if p.pos >= len(p.ibuf) {
		i := <-p.items
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
	_, ch := lex(input)
	p := &parser{ch, nil, []item{}, 0}
	return p.run()

}

func handleOperator(p *parser) nodeComparison {
	var nc nodeComparison
	for {
		i := p.next()

		switch i.typ {
		case itemEOF:
			p.backup()
			return nc
		case itemSet:
			return nc
		case itemRange:
			return nc
		case itemVersion:
			ver, _ := semver.New(i.val)
			nc = nodeComparison{eq, ver}
			return nc
		default:
			v := p.next()
			ver, _ := semver.New(v.val)
			nc = nodeComparison{comparators[i.val], ver}
			return nc
		}
	}
}

func handleSet(p *parser) nodeSet {
	var nc nodeComparison
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
			nc = handleOperator(p)
			set.comparisons = append(set.comparisons, nc)

		}
	}
}

func handleRange(p *parser) node {
	var ns nodeSet
	var rng nodeRange

	for {
		i := p.next()
		switch i.typ {
		case itemEOF:
			return rng
		default:
			p.backup()
			ns = handleSet(p)
			rng.sets = append(rng.sets, ns)

		}
	}
}
