package astjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Scan(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected token
	}{
		"eof": {input: ``, expected: token{
			tp:       EOF,
			leftPos:  0,
			rightPos: 0,
		}},
		"string": {input: `"123"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 5,
		}},
		"positive integer": {input: "999", expected: token{
			tp:       Number,
			leftPos:  0,
			rightPos: 3,
		}},
		"negative integer": {input: "-999", expected: token{
			tp:       Number,
			leftPos:  0,
			rightPos: 4,
		}},
		"positive float": {input: "0.99", expected: token{
			tp:       Number,
			leftPos:  0,
			rightPos: 4,
		}},
		"negative float": {input: "-0.99", expected: token{
			tp:       Number,
			leftPos:  0,
			rightPos: 5,
		}},
		"zero": {input: "0", expected: token{
			tp:       Number,
			leftPos:  0,
			rightPos: 1,
		}},
		"null": {input: `null`, expected: token{
			tp:       Null,
			leftPos:  0,
			rightPos: 4,
		}},
		"true": {input: `true`, expected: token{
			tp:       Bool,
			leftPos:  0,
			rightPos: 4,
		}},
		"false": {input: `false`, expected: token{
			tp:       Bool,
			leftPos:  0,
			rightPos: 5,
		}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			l := newLexer([]byte(tc.input))
			tk := l.Scan()
			assert.Equal(t, tc.expected, tk)
		})
	}
}
