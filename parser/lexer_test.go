package parser

import (
	"testing"

	"github.com/fatih/color"
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

type lexerTestables struct {
	expected bool
	value    string
	result   []itemType
}

var constraints = []lexerTestables{
	{true, "1.0.0 || >=2.5.0 || 5.0.0 - 7.2.3",
		[]itemType{
			itemVersion,
			itemRange,
			itemOperator,
			itemVersion,
			itemRange,
			itemVersion,
			itemAdvanced,
			itemVersion,
			itemEOF,
		},
	},
	{true, "~1.2.3",
		[]itemType{
			itemAdvanced,
			itemVersion,
			itemEOF,
		},
	},
	{true, "^4.5.2-alpha.1",
		[]itemType{
			itemAdvanced,
			itemVersion,
			itemEOF,
		},
	},
	{true, "=2.3.2",
		[]itemType{
			itemOperator,
			itemVersion,
			itemEOF,
		},
	},
	{true, "<=1.2.3",
		[]itemType{
			itemOperator,
			itemVersion,
			itemEOF,
		},
	},
	{true, "5.3.5||4.3.5",
		[]itemType{
			itemVersion,
			itemRange,
			itemVersion,
			itemEOF,
		},
	},
	{true, "5.3.5 ||4.3.5",
		[]itemType{
			itemVersion,
			itemRange,
			itemVersion,
			itemEOF,
		},
	},
	{true, "5.3.5|| 4.3.5",
		[]itemType{
			itemVersion,
			itemRange,
			itemVersion,
			itemEOF,
		},
	},
	{true, "5.3.5 4.3.5",
		[]itemType{
			itemVersion,
			itemSet,
			itemVersion,
			itemEOF,
		},
	},
	{true, ">=1.2.3",
		[]itemType{
			itemOperator,
			itemVersion,
			itemEOF,
		},
	},
	//
	{false, "~ 1.2.3",
		[]itemType{
			itemAdvanced,
			itemError,
		},
	},
	{false, ">= 1.2.3",
		[]itemType{
			itemError,
		},
	},
	{false, "1.2.3 >=",
		[]itemType{
			itemVersion,
			itemSet,
			itemError,
		},
	},
	{false, "5.3.5 |1| 4.3.5",
		[]itemType{
			itemVersion,
			itemError,
		},
	},
	{false, "5. 4.4",
		[]itemType{
			itemError,
		},
	},
	{false, "<1<1",
		[]itemType{
			itemOperator,
			itemError,
		},
	},
	{false, "<1||",
		[]itemType{
			itemOperator,
			itemVersion,
			itemError,
		},
	},
	{false, "M",
		[]itemType{
			itemError,
		},
	},
	{true, "1.0",
		[]itemType{
			itemAdvanced,
			itemEOF,
		},
	},
	{true, "1.x",
		[]itemType{
			itemAdvanced,
			itemEOF,
		},
	},
	{false, "1.x+98uihuhyg",
		[]itemType{
			itemError,
		},
	},
	{true, "1.*.2",
		[]itemType{
			itemAdvanced,
			itemEOF,
		},
	},
	{true, "1.*.2-beta",
		[]itemType{
			itemAdvanced,
			itemEOF,
		},
	},
	{true, "*.1.2",
		[]itemType{
			itemAdvanced,
			itemEOF,
		},
	},
	{false, "1x.2.*",
		[]itemType{
			itemError,
		},
	},
	{false, "1.x2.*",
		[]itemType{
			itemError,
		},
	},
	{false, "1...1",
		[]itemType{
			itemError,
		},
	},
	{false, "1.x.",
		[]itemType{
			itemError,
		},
	},
}

var cyan = color.New(color.FgCyan).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()

func TestLexer(t *testing.T) {
	for _, c := range constraints {
		_, ch := lex(c.value)
		result := true
		x := 0
		for i := range ch {

			result = (i.typ != itemError)

			if i.typ != c.result[x] {
				t.Logf("lex(%v) => %v, want %v \n", cyan(c.value), items[i.typ], items[c.result[x]])
				t.Logf("lex(%v) => %v: %v \n", cyan(c.value), items[i.typ], yellow(i))
			}
			x++
		}
		if result != c.expected {
			t.Logf("lex(%v) => %t, want %t \n", cyan(c.value), result, c.expected)
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
