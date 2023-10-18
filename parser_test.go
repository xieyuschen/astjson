package astjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parse_Literal(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected *Value
	}{
		"eof": {input: "", expected: nil},
		"string": {input: `"123"`, expected: &Value{
			tp:  NtString,
			val: StringAst("123"),
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
			tp:  NtNull,
			val: &NullAst{},
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

func Test_Parse_Object(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *Value
	}{
		{
			name:  "empty object",
			input: "{}",
			expected: &Value{
				tp:  NtObject,
				val: &ObjectAst{map[Value]Value{}},
			},
		},
		{
			name:  "string: string",
			input: `{"123": "123"}`,
			expected: &Value{
				tp: NtObject,
				val: &ObjectAst{map[Value]Value{
					Value{tp: NtString, val: StringAst("123")}: {tp: NtString, val: StringAst("123")}},
				},
			},
		},
		{
			name:  "string: number",
			input: `{"123": 123}`,
			expected: &Value{
				tp: NtObject,
				val: &ObjectAst{map[Value]Value{
					Value{tp: NtString, val: StringAst("123")}: {tp: NtNumber, val: NumberAst(123)}},
				},
			},
		},
		{
			name:  "string: bool true",
			input: `{"123": true}`,
			expected: &Value{
				tp: NtObject,
				val: &ObjectAst{map[Value]Value{
					Value{tp: NtString, val: StringAst("123")}: {tp: NtBool, val: BoolAst(true)}},
				},
			},
		},
		{
			name:  "string: bool false",
			input: `{"123": false}`,
			expected: &Value{
				tp: NtObject,
				val: &ObjectAst{map[Value]Value{
					Value{tp: NtString, val: StringAst("123")}: {tp: NtBool, val: BoolAst(false)}},
				},
			},
		},
		{
			name:  "string: null",
			input: `{"123": null}`,
			expected: &Value{
				tp: NtObject,
				val: &ObjectAst{map[Value]Value{
					Value{tp: NtString, val: StringAst("123")}: {tp: NtNull, val: &NullAst{}}},
				},
			},
		},
		{
			name:  "string: null and string: null",
			input: `{"123": null, "12": null}`,
			expected: &Value{
				tp: NtObject,
				val: &ObjectAst{map[Value]Value{
					Value{tp: NtString, val: StringAst("123")}: {tp: NtNull, val: &NullAst{}},
					Value{tp: NtString, val: StringAst("12")}:  {tp: NtNull, val: &NullAst{}},
				},
				},
			},
		},
		{
			name: "all possible cases",
			input: `
			{
				"str": "123\b\t\r\n",
				"num": 123,
				"bool": true,
				"null": null,
				"empty": {},
				"integrate": { "hello": "world" }
			}`,
			expected: &Value{
				tp: NtObject,
				val: &ObjectAst{map[Value]Value{
					Value{tp: NtString, val: StringAst("str")}:   {tp: NtString, val: StringAst(`123\b\t\r\n`)},
					Value{tp: NtString, val: StringAst("num")}:   {tp: NtNumber, val: NumberAst(123)},
					Value{tp: NtString, val: StringAst("bool")}:  {tp: NtBool, val: BoolAst(true)},
					Value{tp: NtString, val: StringAst("null")}:  {tp: NtNull, val: &NullAst{}},
					Value{tp: NtString, val: StringAst("empty")}: {tp: NtObject, val: &ObjectAst{map[Value]Value{}}},
					Value{tp: NtString, val: StringAst("integrate")}: {tp: NtObject, val: &ObjectAst{map[Value]Value{
						Value{tp: NtString, val: StringAst("hello")}: {tp: NtString, val: StringAst("world")},
					}}},
				},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := Parse([]byte(tc.input))
			assert.Equal(t, tc.expected, value)
		})
	}
}
