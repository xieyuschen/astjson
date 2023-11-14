package astjson

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNumberAst_GetInt64(t *testing.T) {
	testCases := map[string]struct {
		input    *Value
		overflow bool
		expected int64
	}{
		"positive unsigned integer": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: unsignedInteger,
					u:  999,
				},
			},
			expected: 999,
		},
		"negative integer": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: integer,
					i:  -999,
				},
			},
			expected: -999,
		},
		"max positive integer64": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: integer,
					i:  9223372036854775807,
				},
			},
			expected: 9223372036854775807,
		},
		"max negative integer64": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: integer,
					i:  -9223372036854775808,
				},
			},
			expected: -9223372036854775808,
		},
		"max unsigned integer64": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: unsignedInteger,
					u:  18446744073709551615,
				},
			},
			overflow: true,
		},
		"zero": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{Nt: unsignedInteger, u: 0},
			},
			expected: 0,
		},

		"positive float which overflow int": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  float64(18446744073709551615),
				},
			},
			overflow: true,
		},
		"positive float": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  0.99,
				},
			},
			expected: 0,
		},
		"positive float1": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  1.49,
				},
			},
			expected: 1,
		},
		"positive float2": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  1.99,
				},
			},
			expected: 1,
		},
		"negative float": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  -0.99,
				},
			},
			expected: 0,
		},
		"negative float2": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  -1.99,
				},
			},
			expected: -1,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.input.AstValue.(NumberAst).GetInt64()
			if tc.overflow {
				assert.NotEqual(t, tc.expected, result)
				return
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNumberAst_GetUint64(t *testing.T) {
	testCases := map[string]struct {
		input    *Value
		overflow bool
		expected uint64
	}{
		"positive unsigned integer": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: unsignedInteger,
					u:  999,
				},
			},
			expected: 999,
		},
		"negative integer": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: integer,
					i:  -999,
				},
			},
			// overflow
			expected: math.MaxUint - 999 + 1,
		},
		"max positive integer64": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: integer,
					i:  9223372036854775807,
				},
			},
			expected: 9223372036854775807,
		},
		"max negative integer64": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: integer,
					i:  -9223372036854775808,
				},
			},
			expected: 1 << 63,
		},
		"max unsigned integer64": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: unsignedInteger,
					u:  18446744073709551615,
				},
			},
			overflow: true,
		},
		"zero": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{Nt: unsignedInteger, u: 0},
			},
			expected: 0,
		},

		"positive float which overflow int": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  float64(18446744073709551615),
				},
			},
			overflow: true,
		},
		"positive float": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  0.99,
				},
			},
			expected: 0,
		},
		"positive float1": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  1.49,
				},
			},
			expected: 1,
		},
		"positive float2": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  1.99,
				},
			},
			expected: 1,
		},
		"negative float": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  -0.99,
				},
			},
			expected: 0,
		},
		"negative float2": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  -1.99,
				},
			},
			// the behavior differs in different platform
			// arm result is 0, amd is 0xffffffffffffffff
			overflow: true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.input.AstValue.(NumberAst).GetUint64()
			if tc.overflow {
				return
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNumberAst_GetFloat64(t *testing.T) {
	testCases := map[string]struct {
		input    *Value
		overflow bool
		expected float64
	}{
		"positive unsigned integer": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: unsignedInteger,
					u:  999,
				},
			},
			expected: 999,
		},
		"negative integer": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: integer,
					i:  -999,
				},
			},
			// overflow
			expected: -999,
		},
		"max positive integer64": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: integer,
					i:  9223372036854775807,
				},
			},
			expected: 9223372036854775807,
		},
		"max negative integer64": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: integer,
					i:  -9223372036854775808,
				},
			},
			expected: -9223372036854775808,
		},
		"max unsigned integer64": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: unsignedInteger,
					u:  18446744073709551615,
				},
			},
			overflow: true,
		},
		"zero": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{Nt: unsignedInteger, u: 0},
			},
			expected: 0,
		},

		"positive float which overflow int": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  float64(18446744073709551615),
				},
			},
			overflow: true,
		},
		"positive float": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  0.99,
				},
			},
			expected: 0.99,
		},
		"positive float1": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  1.49,
				},
			},
			expected: 1.49,
		},
		"positive float2": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  1.99,
				},
			},
			expected: 1.99,
		},
		"negative float": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  -0.99,
				},
			},
			expected: -0.99,
		},
		"negative float2": {
			input: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					Nt: floatNumber,
					f:  -1.99,
				},
			},
			expected: -1.99,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.input.AstValue.(NumberAst).GetFloat64()
			if tc.overflow {
				assert.NotEqual(t, tc.expected, result)
				return
			}
			assert.Equal(t, tc.expected, result)
		})
	}

}
