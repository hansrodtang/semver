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

type results []itemType

type lexerTestables struct {
	expected bool
	value    string
	result   results
}

var constraints = []*lexerTestables{
	{true, "1.0.0 || >=2.5.0 || 5.0.0 - 7.2.3",
		results{
			itemVersion,
			itemRange,
			itemOperator,
			itemVersion,
			itemRange,
			itemVersion,
			itemAdvanced,
			itemVersion,
		},
	},
	// Operators
	{true, "=2.3.2",
		results{itemOperator, itemVersion},
	},
	{true, "<=1.2.3",
		results{itemOperator, itemVersion},
	},
	{true, ">=1.2.3",
		results{itemOperator, itemVersion},
	},
	// Sets
	{true, "5.3.5 4.3.5",
		results{itemVersion, itemSet, itemVersion},
	},
	//Ranges
	{true, "5.3.5||4.3.5",
		results{itemVersion, itemRange, itemVersion},
	},
	{true, "5.3.5 ||4.3.5",
		results{itemVersion, itemRange, itemVersion},
	},
	{true, "5.3.5|| 4.3.5",
		results{itemVersion, itemRange, itemVersion},
	},
	// Tilde and Caret Ranges
	{false, "~ 1.2.3",
		results{itemAdvanced},
	},
	{true, "~1.2.3",
		results{itemAdvanced, itemVersion},
	},
	{true, "^4.5.2-alpha.1",
		results{itemAdvanced, itemVersion},
	},
	{false, ">= 1.2.3",
		results{},
	},
	// X-Ranges
	{true, "*",
		results{itemAdvanced},
	},
	{true, "1.0",
		results{itemAdvanced},
	},
	{true, "1.x",
		results{itemAdvanced},
	},
	{false, "1.x+98uihuhyg",
		results{},
	},
	{true, "1.*.2",
		results{itemAdvanced},
	},
	{true, "1.*.2 || 1.x.4",
		results{itemAdvanced, itemSet, itemAdvanced},
	},
	{true, "1.*.2-beta",
		results{itemAdvanced},
	},
	{true, "*.1.2",
		results{itemAdvanced},
	},
	{false, "1x.2.*",
		results{},
	},
	{false, "1.x2.*",
		results{},
	},
	{false, "1...1",
		results{},
	},
	{false, "1.x.",
		results{},
	},

	// Assorted syntax errors
	{false, "1.2.3 >=",
		results{itemVersion, itemSet},
	},
	{false, "5.3.5 |1| 4.3.5",
		results{itemVersion},
	},
	{false, "5. 4.4",
		results{},
	},
	{false, "<1<1",
		results{itemOperator},
	},
	{false, "<1||",
		results{itemOperator, itemVersion},
	},
	{false, "M",
		results{},
	},
}

func init() {
	// Appends appropriate end token based on expected result.
	for _, c := range constraints {
		if c.expected {
			c.result = append(c.result, itemEOF)
		} else {
			c.result = append(c.result, itemError)
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

			if i.typ != c.result[x] {
				t.Logf("lex(%v) => %v(%v), want %v \n", cyan(c.value), items[i.typ], yellow(i), items[c.result[x]])
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
