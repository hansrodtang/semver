package semver

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type stateFn func(*lexer) stateFn

const (
	itemVersion  itemType = iota // Version string
	itemOperator                 // <, <=, >, >= =, ~, ^
	itemSet                      // Set seperated by whitespace
	itemRange                    // || , -
	itemError
	itemEOF // End of input

	versionDEL = '.'
	operatorGT = '>'
	operatorGE = ">="
	operatorLT = '<'
	operatorLE = "<="
	operatorEQ = '='
	operatorTR = '~'
	operatorCR = '^'

	operatorRG = '|'
	operatorST = ' '
	operatorHY = '-'

	eof = -1

	numbers   string = "0123456789"
	letters          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-"
	allchars         = alphanum + delimiters
	alphanum         = letters + numbers
	wildcards        = "Xx*"
)

type Node interface {
	Run() bool
}

type NodeComparison struct {
	action comparatorFunc
	arg    *Version
}

func (n NodeComparison) Run(main *Version) bool {
	return n.action(main, n.arg)
}

type NodeRange struct {
	comparisons []NodeSet
}

func (n NodeRange) Run(main *Version) bool {
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

func (n NodeSet) Run(main *Version) bool {
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

type itemType int

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	}
	return fmt.Sprintf("%v", i.val)
}

type lexer struct {
	name  string    // used only for error reports.
	input string    // the string being scanned.
	start int       // start position of this item.
	pos   int       // current position in the input.
	width int       // width of last rune read from input.
	items chan item // channel of scanned items.
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run() // Concurrently run state machine.
	return l, l.items
}

func (l *lexer) run() {
	for state := lexMain; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// next returns the next rune in the input.
func (l *lexer) next() (rn rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	rn, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rn
}

func (l *lexer) ignore() {
	l.start = l.pos
}

// peek returns but does not consume
// the next rune in the input.
func (l *lexer) peek() rune {
	rn := l.next()
	l.backup()
	return rn
}

func (l *lexer) backup() {
	l.pos -= l.width
}

// accept consumes the next rune
// if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func lexMain(l *lexer) stateFn {
	// Correctly reached EOF.
	switch r := l.next(); {

	case r == eof || r == '\n':
		l.emit(itemEOF) // Useful to make EOF a token.
		return nil      // Stop the run loop.

	case '0' <= r && r <= '9':
		l.backup()
		return lexVersion

	case r == operatorLT:
		l.backup()
		return lexOperator
	case r == operatorGT:
		l.backup()
		return lexOperator
	case r == operatorEQ:
		l.backup()
		return lexOperator
	case r == operatorTR:
		l.backup()
		return lexOperator
	case r == operatorCR:
		l.backup()
		return lexOperator
	case r == operatorRG:
		l.backup()
		return lexRange
	case r == operatorST:
		l.backup()
		return lexRange
	default:
		l.errorf("no version data found")
	}
	return nil
}

func lexVersion(l *lexer) stateFn {
	l.acceptRun(numbers + dot)
	if l.accept("+-") {
		l.acceptRun(allchars)
	}
	l.emit(itemVersion)
	return lexMain
}

func lexOperator(l *lexer) stateFn {

	if l.accept(string(operatorGT)) {
		l.accept(string(operatorEQ))
		l.emit(itemOperator)
		return lexMain
	}
	if l.accept(string(operatorLT)) {
		l.accept(string(operatorEQ))
		l.emit(itemOperator)
		return lexMain
	}
	if l.accept(string(operatorEQ)) {
		l.emit(itemOperator)
		return lexMain
	}
	if l.accept(string(operatorCR)) {
		l.emit(itemOperator)
		return lexMain
	}
	if l.accept(string(operatorTR)) {
		l.emit(itemOperator)
	}
	return lexMain
}

func lexRange(l *lexer) stateFn {
	if l.accept(string(operatorRG)) {
		if l.accept(string(operatorRG)) {
			l.emit(itemRange)
			if l.peek() == operatorST {
				l.next()
				l.ignore()
			}
			return lexMain
		}
	}
	if l.accept(string(operatorST)) {
		if l.peek() == operatorRG || l.peek() == operatorHY {
			l.ignore()
			if l.accept(string(operatorRG)) {
				if l.accept(string(operatorRG)) {
					l.emit(itemRange)
					if l.peek() == operatorST {
						l.next()
						l.ignore()
					}
					return lexMain
				}
			}
			if l.accept(string(operatorHY)) {
				l.emit(itemRange)
				if l.peek() == operatorST {
					l.next()
					l.ignore()
				}
				return lexMain
			}
		}
		l.emit(itemSet)
	}
	return lexMain
}
