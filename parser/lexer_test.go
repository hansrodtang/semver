package parser

import (
	"fmt"
	"testing"
)

var items = map[itemType]string{
	itemVersion:  "itemVersion",
	itemOperator: "itemOperator",
	itemSet:      "itemSet",
	itemRange:    "itemRange",
	itemAdvanced: "itemAdvanced",
	itemError:    "itemError",
	itemEOF:      "itemEOF",
}

var constraints = []string{
	"1.0 || >=2.5.0 || 5.0.0 - 7.2.3",
	"~1.2.3",
	"^4.5.2-alpha.1",
	"=2.3.2",
	"<=1.2.3",
	"5.3.5||4.3.5",
	"5.3.5 ||4.3.5",
	"5.3.5|| 4.3.5",
	"5.3.5 4.3.5",
	">=1.2.3",
	">= 1.2.3",
	"M",
}

// Just for debugging, not a real test. REMOVE THIS.
func TestLexer(t *testing.T) {
	for _, c := range constraints {
		_, ch := lex(c)
		for {
			s, ok := <-ch
			if ok != false {
				fmt.Printf("%v: '%v' \n", items[s.typ], s)
			} else {
				break
			}
		}
	}
}

// Poor implementation, just for initial testing.
func BenchmarkLexerComplex(b *testing.B) {
	const VERSION = "1.0.0 || >=2.5.0 || 5.0.0 - 7.2.3 || ~4.3.1 ^2.1.1"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, ch := lex(VERSION)
		for {
			_, ok := <-ch
			if ok == false {
				//fmt.Printf("%v: '%v' \n", items[s.typ], s)
				//} else {
				break
			}
		}
	}
}

// Poor implementation, just for initial testing.
func BenchmarkLexerSimple(b *testing.B) {
	const VERSION = "1.0.0"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, ch := lex(VERSION)
		for {
			_, ok := <-ch
			if ok == false {
				//fmt.Printf("%v: '%v' \n", items[s.typ], s)
				//} else {
				break
			}
		}
	}
}
