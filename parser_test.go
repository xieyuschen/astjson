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
			NodeType: String,
			AstValue: StringAst("123"),
		}},
		"positive integer": {input: "999", expected: &Value{
			NodeType: Number,
			AstValue: NumberAst{
				Nt: unsignedInteger,
				u:  999,
			},
		}},
		"negative integer": {input: "-999", expected: &Value{
			NodeType: Number,
			AstValue: NumberAst{
				Nt: integer,
				i:  -999,
			},
		}},
		"zero": {input: "0", expected: &Value{
			NodeType: Number,
			AstValue: NumberAst{Nt: unsignedInteger, u: 0},
		}},
		"positive float": {input: "0.99", expected: &Value{
			NodeType: Number,
			AstValue: NumberAst{
				Nt: floatNumber,
				f:  0.99,
			},
		}},
		"negative float": {input: "-0.99", expected: &Value{
			NodeType: Number,
			AstValue: NumberAst{
				Nt: floatNumber,
				f:  -0.99,
			},
		}},
		"null": {input: "null", expected: &Value{
			NodeType: Null,
			AstValue: &NullAst{},
		}},
		"true": {input: "true", expected: &Value{
			NodeType: Bool,
			AstValue: BoolAst(true),
		}},
		"false": {input: "false", expected: &Value{
			NodeType: Bool,
			AstValue: BoolAst(false),
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
				NodeType: Object,
				AstValue: &ObjectAst{map[string]Value{}},
			},
		},
		{
			name:  "string: string",
			input: `{"123": "123"}`,
			expected: &Value{
				NodeType: Object,
				AstValue: &ObjectAst{map[string]Value{
					"123": {NodeType: String, AstValue: StringAst("123")}},
				},
			},
		},
		{
			name:  "string: number",
			input: `{"123": 123}`,
			expected: &Value{
				NodeType: Object,
				AstValue: &ObjectAst{map[string]Value{
					"123": {NodeType: Number, AstValue: NumberAst{
						Nt: unsignedInteger,
						u:  123,
					}}},
				},
			},
		},
		{
			name:  "string: bool true",
			input: `{"123": true}`,
			expected: &Value{
				NodeType: Object,
				AstValue: &ObjectAst{map[string]Value{
					"123": {NodeType: Bool, AstValue: BoolAst(true)}},
				},
			},
		},
		{
			name:  "string: bool false",
			input: `{"123": false}`,
			expected: &Value{
				NodeType: Object,
				AstValue: &ObjectAst{map[string]Value{
					"123": {NodeType: Bool, AstValue: BoolAst(false)}},
				},
			},
		},
		{
			name:  "string: null",
			input: `{"123": null}`,
			expected: &Value{
				NodeType: Object,
				AstValue: &ObjectAst{map[string]Value{
					"123": {NodeType: Null, AstValue: &NullAst{}}},
				},
			},
		},
		{
			name:  "string: null and string: null",
			input: `{"123": null, "12": null}`,
			expected: &Value{
				NodeType: Object,
				AstValue: &ObjectAst{map[string]Value{
					"123": {NodeType: Null, AstValue: &NullAst{}},
					"12":  {NodeType: Null, AstValue: &NullAst{}},
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
				NodeType: Array,
				AstValue: &ArrayAst{},
			},
		},
		{
			name:  "single string array",
			input: `[ "123"]`,
			expected: &Value{
				NodeType: Array,
				AstValue: &ArrayAst{[]Value{
					{NodeType: String, AstValue: StringAst("123")},
				}},
			},
		},
		{
			name:  "double string array",
			input: `[ "123", "456"]`,
			expected: &Value{
				NodeType: Array,
				AstValue: &ArrayAst{[]Value{
					{NodeType: String, AstValue: StringAst("123")},
					{NodeType: String, AstValue: StringAst("456")},
				}},
			},
		},
		{
			name:  "int array",
			input: `[ -1,0,1]`,
			expected: &Value{
				NodeType: Array,
				AstValue: &ArrayAst{[]Value{
					{NodeType: Number, AstValue: NumberAst{Nt: integer, i: -1}},
					{NodeType: Number, AstValue: NumberAst{Nt: unsignedInteger, u: 0}},
					{NodeType: Number, AstValue: NumberAst{Nt: unsignedInteger, u: 1}},
				}},
			},
		},
		{
			name:  "float array",
			input: `[ -0.99, 0, 9.99 ]`,
			expected: &Value{
				NodeType: Array,
				AstValue: &ArrayAst{[]Value{
					{NodeType: Number, AstValue: NumberAst{Nt: floatNumber, f: -0.99}},
					{NodeType: Number, AstValue: NumberAst{Nt: unsignedInteger, u: 0}},
					{NodeType: Number, AstValue: NumberAst{Nt: floatNumber, f: 9.99}},
				}},
			},
		},
		{
			name:  "null array",
			input: `[ null, null ]`,
			expected: &Value{
				NodeType: Array,
				AstValue: &ArrayAst{[]Value{
					{NodeType: Null, AstValue: &NullAst{}},
					{NodeType: Null, AstValue: &NullAst{}},
				}},
			},
		},
		{
			name:  "bool array",
			input: `[ true, false ]`,
			expected: &Value{
				NodeType: Array,
				AstValue: &ArrayAst{[]Value{
					{NodeType: Bool, AstValue: BoolAst(true)},
					{NodeType: Bool, AstValue: BoolAst(false)},
				}},
			},
		},
		{
			name:  "empty array of array",
			input: `[ [] ]`,
			expected: &Value{
				NodeType: Array,
				AstValue: &ArrayAst{[]Value{
					{NodeType: Array, AstValue: &ArrayAst{}},
				}},
			},
		},
		{
			name:  "two empty array of array",
			input: `[ [], [] ]`,
			expected: &Value{
				NodeType: Array,
				AstValue: &ArrayAst{[]Value{
					{NodeType: Array, AstValue: &ArrayAst{}},
					{NodeType: Array, AstValue: &ArrayAst{}},
				}},
			},
		},
		{
			name:  "embed string array of array",
			input: `[ ["123"], ["123"] ]`,
			expected: &Value{
				NodeType: Array,
				AstValue: &ArrayAst{[]Value{
					{NodeType: Array, AstValue: &ArrayAst{[]Value{
						{NodeType: String, AstValue: StringAst("123")},
					}}},
					{NodeType: Array, AstValue: &ArrayAst{[]Value{
						{NodeType: String, AstValue: StringAst("123")},
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
				NodeType: Object,
				AstValue: &ObjectAst{map[string]Value{
					"str":   {NodeType: String, AstValue: StringAst(`123\b\t\r\n`)},
					"num":   {NodeType: Number, AstValue: NumberAst{Nt: unsignedInteger, u: 123}},
					"bool":  {NodeType: Bool, AstValue: BoolAst(true)},
					"null":  {NodeType: Null, AstValue: &NullAst{}},
					"empty": {NodeType: Object, AstValue: &ObjectAst{map[string]Value{}}},
					"embed-object": {
						NodeType: Object,
						AstValue: &ObjectAst{map[string]Value{
							"hello": {NodeType: String, AstValue: StringAst("world")},
						}}},
					"array-in-object": {
						NodeType: Object,
						AstValue: &ObjectAst{map[string]Value{
							"hello": {
								NodeType: Array,
								AstValue: &ArrayAst{Values: []Value{{NodeType: String, AstValue: StringAst("world")}}}}},
						}},
					"array": {
						NodeType: Array,
						AstValue: &ArrayAst{[]Value{
							{NodeType: String, AstValue: StringAst("world")},
						}},
					},
					"empty-array": {
						NodeType: Array,
						AstValue: &ArrayAst{},
					},
					"embed-empty-array": {
						NodeType: Array,
						AstValue: &ArrayAst{[]Value{
							{NodeType: Array, AstValue: &ArrayAst{}},
							{NodeType: Array, AstValue: &ArrayAst{}},
						}},
					},
					"array-empty-obj": {
						NodeType: Array,
						AstValue: &ArrayAst{[]Value{
							{NodeType: Object, AstValue: &ObjectAst{KvMap: map[string]Value{}}},
							{NodeType: Object, AstValue: &ObjectAst{KvMap: map[string]Value{}}},
						}},
					},
					"array-obj": {
						NodeType: Array,
						AstValue: &ArrayAst{[]Value{
							{NodeType: Object, AstValue: &ObjectAst{map[string]Value{"hello": {NodeType: String, AstValue: StringAst("world")}}}},
							{NodeType: Object, AstValue: &ObjectAst{map[string]Value{"hello": {NodeType: String, AstValue: StringAst("world")}}}},
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
