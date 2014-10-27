package parser

import (
	"testing"

	"github.com/fatih/color"
)

type results []item

type lexerTestables struct {
	expected bool
	value    string
	result   results
}

var constraints = []*lexerTestables{
	{true, "1.0.0 || >=2.5.0 || 5.0.0 - 7.2.3",
		results{
			{itemVersion, "1.0.0"},
			{itemRange, "||"},
			{itemOperator, ">="},
			{itemVersion, "2.5.0"},
			{itemRange, "||"},
			{itemVersion, "5.0.0"},
			{itemAdvanced, "-"},
			{itemVersion, "7.2.3"},
		},
	},
	{true, "4.5.2-alpha-1",
		results{{itemVersion, "4.5.2-alpha-1"}},
	},
	// Operators
	{true, "=2.3.2",
		results{{itemOperator, "="}, {itemVersion, "2.3.2"}},
	},
	{true, "<=1.2.3",
		results{{itemOperator, "<="}, {itemVersion, "1.2.3"}},
	},
	{true, ">=1.2.3",
		results{{itemOperator, ">="}, {itemVersion, "1.2.3"}},
	},
	// Sets
	{true, "5.3.5 4.3.5",
		results{{itemVersion, "5.3.5"}, {itemSet, " "}, {itemVersion, "4.3.5"}},
	},
	//Ranges
	{true, "5.3.5||4.3.5",
		results{{itemVersion, "5.3.5"}, {itemRange, "||"}, {itemVersion, "4.3.5"}},
	},
	{true, "5.3.5 ||4.3.5",
		results{{itemVersion, "5.3.5"}, {itemRange, "||"}, {itemVersion, "4.3.5"}},
	},
	{true, "5.3.5|| 4.3.5",
		results{{itemVersion, "5.3.5"}, {itemRange, "||"}, {itemVersion, "4.3.5"}},
	},
	{false, "5.3.5||  4.3.5",
		results{{itemVersion, "5.3.5"}, {itemRange, "||"}},
	},
	// Tilde and Caret Ranges
	{false, "~ 1.2.3",
		results{{itemAdvanced, "~"}},
	},
	{true, "~1.2.3",
		results{{itemAdvanced, "~"}, {itemVersion, "1.2.3"}},
	},
	{true, "^4.5.2-alpha.1",
		results{{itemAdvanced, "^"}, {itemVersion, "4.5.2-alpha.1"}},
	},
	{false, ">= 1.2.3",
		results{},
	},
	// Hyphen Range
	{false, "1.2.3 -3.2.5",
		results{{itemVersion, "1.2.3"}},
	},
	// X-Ranges
	{true, "*",
		results{{itemAdvanced, "*"}},
	},
	{false, "**",
		results{},
	},
	{true, "1.0",
		results{{itemAdvanced, "1.0"}},
	},
	{true, "1.x",
		results{{itemAdvanced, "1.x"}},
	},
	{true, "*.x",
		results{{itemAdvanced, "*"}},
	},
	{true, "*.1",
		results{{itemAdvanced, "*"}},
	},
	{true, "*.*.*",
		results{{itemAdvanced, "*"}},
	},
	{false, "1.x+98uihuhyg",
		results{},
	},
	{true, "1.*.2",
		results{{itemAdvanced, "1.*"}},
	},
	{true, "1.*.x || 1.x.4",
		results{{itemAdvanced, "1.*"}, {itemRange, "||"}, {itemAdvanced, "1.x"}},
	},
	{false, "1.*.2-beta",
		results{},
	},
	{false, "1.*.2-",
		results{},
	},
	{true, "*.1.2",
		results{{itemAdvanced, "*"}},
	},
	{false, "1x.2.*",
		results{},
	},
	{false, "x1.2.*",
		results{},
	},
	{false, "1.x2.*",
		results{},
	},
	{false, "1...1",
		results{},
	},
	{false, "1...",
		results{},
	},
	{false, "1.x.",
		results{},
	},
	{false, "1.x..",
		results{},
	},

	// Assorted syntax errors
	{false, "1.2.3-",
		results{},
	},
	{false, "1.2.3 >=",
		results{{itemVersion, "1.2.3"}, {itemSet, " "}},
	},
	{false, "5.3.5 |1| 4.3.5",
		results{{itemVersion, "5.3.5"}},
	},
	{false, "1.2.3.1",
		results{},
	},
	{false, "5. 4.4",
		results{},
	},
	{false, "<1<1",
		results{{itemOperator, "<"}},
	},
	{false, "<1||",
		results{{itemOperator, "<"}, {itemAdvanced, "1"}, {itemRange, "||"}},
	},
	{false, "M",
		results{},
	},
}

func init() {
	// Appends appropriate end token based on expected result.
	for _, c := range constraints {
		if c.expected {
			c.result = append(c.result, item{itemEOF, ""})
		} else {
			c.result = append(c.result, item{itemError, ""})
		}
	}
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

			if len(c.result) > x {
				if i.typ != c.result[x].typ {
					t.Errorf("lex(%v) => %v, want %v \n", cyan(c.value), yellow(i), yellow(c.result[x]))
				} else if i.val != c.result[x].val {
					if !(i.typ == itemError || i.typ == itemEOF) {
						t.Errorf("lex(%v) => %v, want %v \n", cyan(c.value), yellow(i), yellow(c.result[x]))
					}
				}
			} else {
				t.Errorf("lex(%v) => %v, want <nil>\n", cyan(c.value), yellow(i))
			}
			x++
		}
		if result != c.expected {
			t.Errorf("lex(%v) => %t, want %t \n", cyan(c.value), result, c.expected)
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
