package astjson

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	mixedNode = &Value{
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
	val, err := NewWalker(mixedNode).
		Optional("key1", func(value *Value) error {
			panic("")
		}).Walk()
	assert.NoError(t, err)
	assert.Equal(t, mixedNode, val)

	val, err = NewWalker(mixedNode).
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
	assert.Equal(t, mixedNode, val)

	val, err = NewWalker(mixedNode).
		Field("str").
		Optional("num", shouldBe123).
		Field("bool").Validate(shouldBeTrue).
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
	assert.Equal(t, mixedNode, val)

	val, err = NewWalker(mixedNode).
		Field("str").
		Optional("non-exist", shouldBe234).Walk()

	assert.NoError(t, err)
	assert.Equal(t, mixedNode, val)

	val, err = NewWalker(mixedNode).
		Field("str").
		Optional("num", shouldBe234).Walk()

	assert.Equal(t, "num should be 234", err.Error())
	assert.Nil(t, val)

	val, err = NewWalker(mixedNode).
		Field("str").
		Field("num").Validate(shouldBe234).
		Walk()

	assert.Equal(t, "num should be 234", err.Error())
	assert.Nil(t, val)

	val, err = NewWalker(mixedNode).
		Field("str").
		ValidateKey("num", shouldBe234).Walk()

	assert.Equal(t, "num should be 234", err.Error())
	assert.Nil(t, val)

}

func Test_WalkPath_and_EndPath(t *testing.T) {
	input := &Value{
		NodeType: Object,
		AstValue: &ObjectAst{map[string]Value{}},
	}
	w := NewWalker(input)

	assert.NotNil(t, w.Path("").err)
	val, err := w.Path("").Walk()
	assert.NotNil(t, err)
	assert.Nil(t, val)

	assert.NotNil(t, w.Path("123").err)
	assert.NotNil(t, w.Path("123").EndPath().err)

	w = NewWalker(mixedNode)
	nw := w.Path("str")
	assert.NoError(t, w.err)
	strNode := &Value{NodeType: String, AstValue: StringAst(`123\b\t\r\n`)}
	assert.Equal(t, strNode, nw.value)
	assert.Equal(t, w, nw.EndPath())

	w = NewWalker(mixedNode)
	nw = w.Path("embed-object").Path("hello")
	assert.NoError(t, w.err)
	strNode = &Value{NodeType: String, AstValue: StringAst("world")}
	assert.Equal(t, strNode, nw.value)
	assert.Equal(t, w, nw.EndPath())
	assert.Equal(t, mixedNode, nw.EndPath().value)
}

func Test_Walk_Object_Mixed(t *testing.T) {
	w := NewWalker(mixedNode).
		Field("str").
		Optional("num", shouldBe123).
		Field("bool").Validate(shouldBeTrue).
		Path("embed-object").Validate(shouldBeObject).
		Path("hello").Validate(shouldBeWorld)

	val, err := w.EndPath().Walk()
	assert.Equal(t, mixedNode, val)
	assert.NoError(t, err)

	val, err = w.Walk()
	assert.Equal(t, mixedNode, val)
	assert.NoError(t, err)
}

func Test_WalkerClone(t *testing.T) {
	assert.Nil(t, NewWalker(nil).Clone())
}
func shouldBeWorld(value *Value) error {
	if value.AstValue.(StringAst) != "world" {
		return errors.New("value should be world")
	}
	return nil
}

func shouldBeObject(value *Value) error {
	if value.AstValue != Object {
		return errors.New("value should be an object")
	}
	return nil
}

func shouldBe234(value *Value) error {
	actual := value.AstValue.(NumberAst).GetInt64()
	if actual == 234 {
		return nil
	}
	return errors.New("num should be 234")
}

func shouldBe123(value *Value) error {
	actual := value.AstValue.(NumberAst).GetInt64()
	if actual == 123 {
		return nil
	}
	return errors.New("num should be 123")
}

func shouldBeTrue(value *Value) error {
	if value.AstValue.(BoolAst) {
		return nil
	}
	return errors.New("bool should be true")
}
