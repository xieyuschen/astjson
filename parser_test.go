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
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := Parse([]byte(tc.input))
			assert.Equal(t, tc.expected, value)
		})
	}
}

func Test_Parse_Array(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *Value
	}{
		{
			name:  "empty array",
			input: "[]",
			expected: &Value{
				tp:  NtArray,
				val: &ArrayAst{},
			},
		},
		{
			name:  "single string array",
			input: `[ "123"]`,
			expected: &Value{
				tp: NtArray,
				val: &ArrayAst{[]Value{
					{tp: NtString, val: StringAst("123")},
				}},
			},
		},
		{
			name:  "double string array",
			input: `[ "123", "456"]`,
			expected: &Value{
				tp: NtArray,
				val: &ArrayAst{[]Value{
					{tp: NtString, val: StringAst("123")},
					{tp: NtString, val: StringAst("456")},
				}},
			},
		},
		{
			name:  "int array",
			input: `[ -1,0,1]`,
			expected: &Value{
				tp: NtArray,
				val: &ArrayAst{[]Value{
					{tp: NtNumber, val: NumberAst(-1)},
					{tp: NtNumber, val: NumberAst(0)},
					{tp: NtNumber, val: NumberAst(1)},
				}},
			},
		},
		{
			name:  "float array",
			input: `[ -0.99, 0, 9.99 ]`,
			expected: &Value{
				tp: NtArray,
				val: &ArrayAst{[]Value{
					{tp: NtNumber, val: NumberAst(-0.99)},
					{tp: NtNumber, val: NumberAst(0)},
					{tp: NtNumber, val: NumberAst(9.99)},
				}},
			},
		},
		{
			name:  "null array",
			input: `[ null, null ]`,
			expected: &Value{
				tp: NtArray,
				val: &ArrayAst{[]Value{
					{tp: NtNull, val: &NullAst{}},
					{tp: NtNull, val: &NullAst{}},
				}},
			},
		},
		{
			name:  "bool array",
			input: `[ true, false ]`,
			expected: &Value{
				tp: NtArray,
				val: &ArrayAst{[]Value{
					{tp: NtBool, val: BoolAst(true)},
					{tp: NtBool, val: BoolAst(false)},
				}},
			},
		},
		{
			name:  "empty array of array",
			input: `[ [] ]`,
			expected: &Value{
				tp: NtArray,
				val: &ArrayAst{[]Value{
					{tp: NtArray, val: &ArrayAst{}},
				}},
			},
		},
		{
			name:  "two empty array of array",
			input: `[ [], [] ]`,
			expected: &Value{
				tp: NtArray,
				val: &ArrayAst{[]Value{
					{tp: NtArray, val: &ArrayAst{}},
					{tp: NtArray, val: &ArrayAst{}},
				}},
			},
		},
		{
			name:  "embed string array of array",
			input: `[ ["123"], ["123"] ]`,
			expected: &Value{
				tp: NtArray,
				val: &ArrayAst{[]Value{
					{tp: NtArray, val: &ArrayAst{[]Value{
						{tp: NtString, val: StringAst("123")},
					}}},
					{tp: NtArray, val: &ArrayAst{[]Value{
						{tp: NtString, val: StringAst("123")},
					}}},
				}},
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

func Test_Parse_Mixture(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *Value
	}{
		{
			name: "all possible cases",
			input: `
			{
				"str": "123\b\t\r\n",
				"num": 123,
				"bool": true,
				"null": null,
				"empty": {},
				"embed-object": { "hello": "world" },
				"array-in-object": { "hello": ["world"] },
				"array": [ "world"],
				"empty-array": [],
			    "embed-empty-array": [[],[]],
				"array-empty-obj": [ {}, {} ],
				"array-obj": [ { "hello": "world" }, { "hello": "world" } ]

			}`,
			expected: &Value{
				tp: NtObject,
				val: &ObjectAst{map[Value]Value{
					Value{tp: NtString, val: StringAst("str")}:   {tp: NtString, val: StringAst(`123\b\t\r\n`)},
					Value{tp: NtString, val: StringAst("num")}:   {tp: NtNumber, val: NumberAst(123)},
					Value{tp: NtString, val: StringAst("bool")}:  {tp: NtBool, val: BoolAst(true)},
					Value{tp: NtString, val: StringAst("null")}:  {tp: NtNull, val: &NullAst{}},
					Value{tp: NtString, val: StringAst("empty")}: {tp: NtObject, val: &ObjectAst{map[Value]Value{}}},
					Value{tp: NtString, val: StringAst("embed-object")}: {
						tp: NtObject,
						val: &ObjectAst{map[Value]Value{
							Value{tp: NtString, val: StringAst("hello")}: {tp: NtString, val: StringAst("world")},
						}}},
					Value{tp: NtString, val: StringAst("array-in-object")}: {
						tp: NtObject,
						val: &ObjectAst{map[Value]Value{
							Value{tp: NtString, val: StringAst("hello")}: {
								tp:  NtArray,
								val: &ArrayAst{values: []Value{{tp: NtString, val: StringAst("world")}}}}},
						}},
					Value{tp: NtString, val: StringAst("array")}: {
						tp: NtArray,
						val: &ArrayAst{[]Value{
							{tp: NtString, val: StringAst("world")},
						}},
					},
					Value{tp: NtString, val: StringAst("empty-array")}: {
						tp:  NtArray,
						val: &ArrayAst{},
					},
					Value{tp: NtString, val: StringAst("embed-empty-array")}: {
						tp: NtArray,
						val: &ArrayAst{[]Value{
							{tp: NtArray, val: &ArrayAst{}},
							{tp: NtArray, val: &ArrayAst{}},
						}},
					},
					Value{tp: NtString, val: StringAst("array-empty-obj")}: {
						tp: NtArray,
						val: &ArrayAst{[]Value{
							{tp: NtObject, val: &ObjectAst{m: map[Value]Value{}}},
							{tp: NtObject, val: &ObjectAst{m: map[Value]Value{}}},
						}},
					},
					Value{tp: NtString, val: StringAst("array-obj")}: {
						tp: NtArray,
						val: &ArrayAst{[]Value{
							{tp: NtObject, val: &ObjectAst{map[Value]Value{{tp: NtString, val: StringAst("hello")}: {tp: NtString, val: StringAst("world")}}}},
							{tp: NtObject, val: &ObjectAst{map[Value]Value{{tp: NtString, val: StringAst("hello")}: {tp: NtString, val: StringAst("world")}}}},
						}},
					},
				}},
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
