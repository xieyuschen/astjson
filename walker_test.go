package astjson

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WalkTopLevel_LiteralAndArray(t *testing.T) {
	testCases := map[string]struct {
		expected *Value
		validate Validator
		errStr   string
	}{
		"string: pass anyway": {
			expected: &Value{
				NodeType: String,
				AstValue: StringAst("123"),
			},
			validate: func(value *Value) error {
				return nil
			},
		},
		"string: check value and pass": {
			expected: &Value{
				NodeType: String,
				AstValue: StringAst("123"),
			},
			validate: func(value *Value) error {
				if value.AstValue.(StringAst) == "123" {
					return nil
				}
				return errors.New("the string should be 123")
			},
		},
		"string: check value and return error": {
			expected: &Value{
				NodeType: String,
				AstValue: StringAst("234"),
			},
			validate: func(value *Value) error {
				if value.AstValue.(StringAst) == "123" {
					return nil
				}
				return errors.New("the string should be 123")
			},
			errStr: "the string should be 123",
		},
		"integer: check value and pass": {
			expected: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					nt: unsignedInteger,
					u:  999,
				},
			},
			validate: func(value *Value) error {
				number := value.AstValue.(NumberAst)
				if number.GetInt64() == 999 {
					return nil
				}
				return errors.New("the number should be 999")
			},
		},
		"integer: check value and return error": {
			expected: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					nt: unsignedInteger,
					u:  888,
				},
			},
			validate: func(value *Value) error {
				number := value.AstValue.(NumberAst)
				if number.GetInt64() == 999 {
					return nil
				}
				return errors.New("the number should be 999")
			},
			errStr: "the number should be 999",
		},
		"unsigned integer: check and pass": {
			expected: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					nt: unsignedInteger,
					u:  999,
				},
			},
			validate: func(value *Value) error {
				number := value.AstValue.(NumberAst)
				if number.GetUint64() == 999 {
					return nil
				}
				return errors.New("the number should be 999")
			},
			errStr: "the number should be 999",
		},
		"unsigned integer:check and return error": {
			expected: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					nt: unsignedInteger,
					u:  888,
				},
			},
			validate: func(value *Value) error {
				number := value.AstValue.(NumberAst)
				if number.GetUint64() == 999 {
					return nil
				}
				return errors.New("the number should be 999")
			},
			errStr: "the number should be 999",
		},
		"negative integer: check and pass": {
			expected: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					nt: integer,
					i:  -1,
				},
			},
			validate: func(value *Value) error {
				number := value.AstValue.(NumberAst)
				if number.GetInt64() == -1 {
					return nil
				}
				return errors.New("the number should be -1")
			},
		},
		"negative integer: check and return error": {
			expected: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					nt: integer,
					i:  -1,
				},
			},
			validate: func(value *Value) error {
				number := value.AstValue.(NumberAst)
				if number.GetInt64() == -1 {
					return nil
				}
				return errors.New("the number should be -1")
			},
		},
		"negative float: check and pass": {
			expected: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					nt: floatNumber,
					f:  -0.99,
				},
			},
			validate: func(value *Value) error {
				number := value.AstValue.(NumberAst)
				if number.GetFloat64() == -0.99 {
					return nil
				}
				return errors.New("the number should be -0.99")
			},
		},
		"negative float: check and return error": {
			expected: &Value{
				NodeType: Number,
				AstValue: NumberAst{
					nt: floatNumber,
					f:  -1.99,
				},
			},
			validate: func(value *Value) error {
				number := value.AstValue.(NumberAst)
				if number.GetFloat64() == -0.99 {
					return nil
				}
				return errors.New("the number should be -0.99")
			},
			errStr: "the number should be -0.99",
		},
		"bool: check and pass": {
			expected: &Value{
				NodeType: Bool,
				AstValue: BoolAst(true),
			},
			validate: func(value *Value) error {
				if value.AstValue.(BoolAst) {
					return nil
				}
				return errors.New("bool value should be true")
			},
		},
		"true: check and return error": {
			expected: &Value{
				NodeType: Bool,
				AstValue: BoolAst(true),
			},
			validate: func(value *Value) error {
				if !value.AstValue.(BoolAst) {
					return nil
				}
				return errors.New("bool value should be false")
			},
			errStr: "bool value should be false",
		},
		"array: check and pass": {
			expected: &Value{
				NodeType: Array,
				AstValue: &ArrayAst{[]Value{
					{NodeType: Number, AstValue: NumberAst{nt: integer, i: -1}},
					{NodeType: Number, AstValue: NumberAst{nt: unsignedInteger, u: 0}},
					{NodeType: Number, AstValue: NumberAst{nt: unsignedInteger, u: 1}},
				}},
			},
			validate: func(value *Value) error {
				if len(value.AstValue.(*ArrayAst).values) == 3 {
					return nil
				}
				return errors.New("array should have three element")
			},
			errStr: "bool value should be false",
		},
		"array: check and return error": {
			expected: &Value{
				NodeType: Array,
				AstValue: &ArrayAst{[]Value{
					{NodeType: Number, AstValue: NumberAst{nt: integer, i: -1}},
					{NodeType: Number, AstValue: NumberAst{nt: unsignedInteger, u: 0}},
					{NodeType: Number, AstValue: NumberAst{nt: unsignedInteger, u: 1}},
					{NodeType: Number, AstValue: NumberAst{nt: unsignedInteger, u: 2}},
				}},
			},
			validate: func(value *Value) error {
				if len(value.AstValue.(*ArrayAst).values) == 3 {
					return nil
				}
				return errors.New("array should have three element")
			},
			errStr: "array should have three element",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			val, err := NewWalker(tc.expected).Validate(tc.validate).Walk()
			if err != nil {
				assert.Equal(t, tc.errStr, err.Error())
				return
			}
			assert.Equal(t, tc.expected, val)
		})
	}
}

func Test_WalkTopLevel_Object_Empty(t *testing.T) {
	input := &Value{
		NodeType: Object,
		AstValue: &ObjectAst{map[string]Value{}},
	}

	val, err := NewWalker(input).
		Optional("key1", func(value *Value) error {
			panic("")
		}).Walk()

	assert.NoError(t, err)
	assert.Equal(t, input, val)

	val, err = NewWalker(input).Field("key").Walk()
	assert.ErrorIs(t, err, ErrFieldNotExist)
	assert.Nil(t, val)

	val, err = NewWalker(input).Field("key").Validate(func(value *Value) error {
		return nil
	}).Walk()
	assert.ErrorIs(t, err, ErrFieldNotExist)
	assert.Nil(t, val)
}

func Test_WalkTopLevel_Object_Mixed(t *testing.T) {
	input := &Value{
		NodeType: Object,
		AstValue: &ObjectAst{map[string]Value{
			"str":   {NodeType: String, AstValue: StringAst(`123\b\t\r\n`)},
			"num":   {NodeType: Number, AstValue: NumberAst{nt: unsignedInteger, u: 123}},
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
						AstValue: &ArrayAst{values: []Value{{NodeType: String, AstValue: StringAst("world")}}}}},
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
					{NodeType: Object, AstValue: &ObjectAst{m: map[string]Value{}}},
					{NodeType: Object, AstValue: &ObjectAst{m: map[string]Value{}}},
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
	}

	val, err := NewWalker(input).
		Optional("key1", func(value *Value) error {
			panic("")
		}).Walk()
	assert.NoError(t, err)
	assert.Equal(t, input, val)

	val, err = NewWalker(input).
		Field("str").
		Field("num").
		Field("bool").
		Field("null").
		Field("empty").
		Field("embed-object").
		Field("array-in-object").
		Field("array").
		Field("empty-array").
		Field("embed-empty-array").
		Field("array-empty-obj").
		Field("embed-empty-array").
		Field("array-obj").
		Walk()

	assert.NoError(t, err)
	assert.Equal(t, input, val)

	val, err = NewWalker(input).
		Field("str").
		Optional("num", func(value *Value) error {
			actual := value.AstValue.(NumberAst).GetInt64()
			if actual == 123 {
				return nil
			}
			return errors.New("num should be 123")
		}).
		Field("bool").
		Validate(func(value *Value) error {
			if value.AstValue.(BoolAst) {
				return nil
			}
			return errors.New("bool should be true")
		}).
		Field("null").
		Field("empty").
		Field("embed-object").
		Field("array-in-object").
		Field("array").
		Field("empty-array").
		Field("embed-empty-array").
		Field("array-empty-obj").
		Field("embed-empty-array").
		Field("array-obj").
		Walk()

	assert.NoError(t, err)
	assert.Equal(t, input, val)

	val, err = NewWalker(input).
		Field("str").
		Optional("non-exist", func(value *Value) error {
			actual := value.AstValue.(NumberAst).GetInt64()
			if actual == 234 {
				return nil
			}
			return errors.New("num should be 234")
		}).Walk()

	assert.NoError(t, err)
	assert.Equal(t, input, val)

	val, err = NewWalker(input).
		Field("str").
		Optional("num", func(value *Value) error {
			actual := value.AstValue.(NumberAst).GetInt64()
			if actual == 234 {
				return nil
			}
			return errors.New("num should be 234")
		}).Walk()

	assert.Equal(t, "num should be 234", err.Error())
	assert.Nil(t, val)

	val, err = NewWalker(input).
		Field("str").
		Field("num").Validate(func(value *Value) error {
		actual := value.AstValue.(NumberAst).GetInt64()
		if actual == 234 {
			return nil
		}
		return errors.New("num should be 234")
	}).Walk()

	assert.Equal(t, "num should be 234", err.Error())
	assert.Nil(t, val)

	val, err = NewWalker(input).
		Field("str").
		ValidateKey("num", func(value *Value) error {
			actual := value.AstValue.(NumberAst).GetInt64()
			if actual == 234 {
				return nil
			}
			return errors.New("num should be 234")
		}).Walk()

	assert.Equal(t, "num should be 234", err.Error())
	assert.Nil(t, val)

}
