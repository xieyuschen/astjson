package astjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parse(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected *Value
	}{
		"eof": {input: "", expected: nil},
		"string": {input: `"123"`, expected: &Value{
			tp:  NtString,
			val: StringAst(`"123"`),
		}},
		"positive integer": {input: "999", expected: &Value{
			tp:  NtNumber,
			val: NumberAst(999),
		}},
		"negative integer": {input: "-999", expected: &Value{
			tp:  NtNumber,
			val: NumberAst(-999),
		}},
		"zero": {input: "0", expected: &Value{
			tp:  NtNumber,
			val: NumberAst(0),
		}},
		"positive float": {input: "0.99", expected: &Value{
			tp:  NtNumber,
			val: NumberAst(0.99),
		}},
		"negative float": {input: "-0.99", expected: &Value{
			tp:  NtNumber,
			val: NumberAst(-0.99),
		}},
		"null": {input: "null", expected: &Value{
			tp: NtNull,
		}},
		"true": {input: "true", expected: &Value{
			tp:  NtBool,
			val: BoolAst(true),
		}},
		"false": {input: "false", expected: &Value{
			tp:  NtBool,
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
