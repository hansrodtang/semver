package parser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type stateFn func(*lexer) stateFn

const (
	itemVersion  itemType = iota // Version string
	itemXRange                   // Version partials
	itemOperator                 // <, <=, >, >= =
	itemSet                      // Set seperated by whitespace
	itemRange                    // || ,
	itemAdvanced                 // ~, ^, -, x-ranges
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

	numbers string = "0123456789"
	letters        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-"

	dot        = "."
	hyphen     = "-"
	plus       = "+"
	delimiters = dot + hyphen + plus

	allchars  = alphanum + delimiters
	alphanum  = letters + numbers
	wildcards = "Xx*"
)

var items = map[itemType]string{
	itemVersion:  "itemVersion",
	itemXRange:   "itemXRange",
	itemOperator: "itemOperator",
	itemSet:      "itemSet",
	itemRange:    "itemRange",
	itemAdvanced: "itemAdvanced",
	itemError:    "itemError",
	itemEOF:      "itemEOF",
}

type itemType int

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	return fmt.Sprintf("%v(%v)", items[i.typ], i.val)
}

type lexer struct {
	input string    // the string being scanned.
	start int       // start position of this item.
	pos   int       // current position in the input.
	width int       // width of last rune read from input.
	items chan item // channel of scanned items.
}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run() // Concurrently run state machine.
	return l
}

func (l *lexer) run() {
	for state := lexMain; state != nil; {
		state = state(l)
	}
}

func (l *lexer) nextItem() item {
	return <-l.items
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

func (l *lexer) rewind() {
	l.pos = l.start
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

func (l *lexer) check(valid string) bool {
	if strings.IndexRune(valid, l.peek()) >= 0 {
		return true
	}
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
	switch r := l.peek(); {

	case r == eof || r == '\n':
		l.emit(itemEOF) // Useful to make EOF a token.
		return nil      // Stop the run loop.

	case '0' <= r && r <= '9':
		return lexVersion
	case r == operatorLT:
		return lexOperator
	case r == operatorGT:
		return lexOperator
	case r == operatorEQ:
		return lexOperator
	case r == operatorTR:
		return lexAdvancedRange
	case r == operatorCR:
		return lexAdvancedRange
	case r == operatorRG:
		return lexRange
	case r == operatorST:
		return lexSet
	case l.check(wildcards):
		return lexAdvancedVersion
	default:
		return l.errorf("invalid character:%v: %q", l.pos, string(r))
	}
}

func lexVersion(l *lexer) stateFn {

	l.acceptRun(numbers)
	if l.accept(dot) {
		if l.accept(numbers) {
			l.acceptRun(numbers)

			if l.accept(dot) {
				if l.accept(numbers) {
					l.acceptRun(numbers)

					if l.accept("+-") {
						if !l.accept(allchars) {
							return l.errorf("invalid character:%v: %q", l.pos, string(l.next()))
						}
						l.acceptRun(allchars)
					}

					if !isEnd(l.peek()) {
						return l.errorf("invalid character:%v: %q", l.pos, string(l.next()))
					}

					l.emit(itemVersion)
					return lexMain
				}
			}
		}
	}
	l.rewind()
	return lexAdvancedVersion
}

func lexOperator(l *lexer) stateFn {
	l.accept(string(operatorGT) + string(operatorLT))
	l.accept(string(operatorEQ))
	if !l.check(numbers) {
		return l.errorf("invalid character:%v: %q", l.pos, string(l.next()))
	}
	l.emit(itemOperator)
	return lexMain
}

func lexSet(l *lexer) stateFn {
	if l.accept(string(operatorST)) {
		if l.peek() == operatorRG {
			l.ignore()
			return lexRange
		}
		if l.peek() == operatorHY {
			l.ignore()
			return lexAdvancedRange
		}
		l.emit(itemSet)
	}
	return lexMain
}

func lexRange(l *lexer) stateFn {
	l.accept(string(operatorRG))
	if l.accept(string(operatorRG)) {
		l.emit(itemRange)
		if l.peek() == operatorST {
			l.next()
			l.ignore()
		}
		if isEnd(l.peek()) {
			return l.errorf("invalid character:%v: %q", l.pos, string(l.next()))
		}
		return lexMain
	}
	return l.errorf("invalid character:%v: %q", l.pos, string(l.next()))

}

func lexAdvancedRange(l *lexer) stateFn {
	if l.accept(string(operatorHY)) {
		if l.peek() == operatorST {
			l.emit(itemAdvanced)
			l.next()
			l.ignore()
		} else {
			return l.errorf("invalid character:%v: %q", l.pos, string(l.next()))
		}
		return lexMain
	}
	if l.accept(string(operatorCR) + string(operatorTR)) {
		l.emit(itemAdvanced)

		if !l.check(numbers) {
			return l.errorf("invalid character:%v: %q", l.pos, string(l.next()))
		}
	}

	return lexMain
}

func lexAdvancedVersion(l *lexer) stateFn {
	// Syntax check
	for i := 0; i <= 2; i++ {
		if !l.accept(wildcards) {
			if !l.accept(numbers) {
				return l.errorf("invalid character:%v: %q", l.pos, string(l.next()))
			}
			l.acceptRun(numbers)
		}

		if !l.accept(dot) {
			if !isEnd(l.peek()) {
				return l.errorf("invalid character:%v: %q", l.pos, string(l.next()))
			}
			l.rewind()
			break
		}
	}

	// Generate item
	for i := 0; i <= 2; i++ {
		if !l.accept(wildcards) {
			l.acceptRun(numbers)
			l.accept(dot)
		} else {
			l.emit(itemXRange)
			break
		}
		if isEnd(l.peek()) {
			l.emit(itemXRange)
			break
		}
	}

	for !isEnd(l.next()) {
	}
	l.backup()
	l.ignore()

	return lexMain

}

func isEnd(r rune) bool {
	return (r == operatorST || r == eof || r == operatorRG)
}
