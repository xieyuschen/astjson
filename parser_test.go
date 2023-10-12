package astjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parse(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected Value
	}{
		"eof": {input: "", expected: Value{
			tp: EOF,
		}},
		"string": {input: `"123"`, expected: Value{
			tp:  String,
			raw: []byte(`"123"`),
			val: StringAst(`"123"`),
		}},
		"positive integer": {input: "999", expected: Value{
			tp:  Number,
			raw: []byte("999"),
			val: NumberAst(999),
		}},
		"negative integer": {input: "-999", expected: Value{
			tp:  Number,
			raw: []byte("-999"),
			val: NumberAst(-999),
		}},
		"zero": {input: "0", expected: Value{
			tp:  Number,
			raw: []byte("0"),
			val: NumberAst(0),
		}},
		"positive float": {input: "0.99", expected: Value{
			tp:  Number,
			raw: []byte("0.99"),
			val: NumberAst(0.99),
		}},
		"negative float": {input: "-0.99", expected: Value{
			tp:  Number,
			raw: []byte("-0.99"),
			val: NumberAst(-0.99),
		}},
		"null": {input: "null", expected: Value{
			tp:  Null,
			raw: []byte("null"),
		}},
		"true": {input: "true", expected: Value{
			tp:  Bool,
			raw: []byte("true"),
			val: BoolAst(true),
		}},
		"false": {input: "false", expected: Value{
			tp:  Bool,
			raw: []byte("false"),
			val: BoolAst(false),
		}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			value := Parse([]byte(tc.input))
			assert.Equal(t, tc.expected, value)
		})
	}
}
