package astjson

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Corner_Cases(t *testing.T) {

	// destination is not a pointer
	var i int
	err := NewDecoder().Unmarshal(NewParser([]byte("")).Parse(), i)
	assert.NotNil(t, err)
	assert.Equal(t, "dest must be a pointer", err.Error())

	// input value is nils
	err = NewDecoder().Unmarshal(nil, &i)
	assert.Equal(t, "value is a nil pointer", err.Error())

	wrongValue := &Value{
		NodeType: 7,
		AstValue: nil,
	}
	err = NewDecoder().Unmarshal(wrongValue, &i)
	assert.Equal(t, "invalid value", err.Error())
}

func Test_Unmarshal_Number(t *testing.T) {

	testCases := map[string]struct {
		input       string
		destination interface{}
		expected    interface{}
	}{
		"int": {
			input:       "-999",
			destination: new(int),
			expected:    -999,
		},
		"int8": {
			input:       "-128",
			destination: new(int8),
			expected:    int8(-128),
		},
		"int16": {
			input:       "-32768",
			destination: new(int16),
			expected:    int16(-32768),
		},
		"int32": {
			input:       "-2147483648",
			destination: new(int32),
			expected:    int32(-2147483648),
		},
		"int64": {
			input:       "-9223372036854775808",
			destination: new(int64),
			expected:    int64(-9223372036854775808),
		},
		"uint": {
			input:       "999",
			destination: new(uint),
			expected:    uint(999),
		},
		"uint8": {
			input:       "255",
			destination: new(uint8),
			expected:    uint8(255),
		},
		"uint16": {
			input:       "65535",
			destination: new(uint16),
			expected:    uint16(65535),
		},
		"uint32": {
			input:       "4294967295",
			destination: new(uint32),
			expected:    uint32(4294967295),
		},
		"uint64": {
			input:       "18446744073709551615",
			destination: new(uint64),
			expected:    uint64(18446744073709551615),
		},
		"float32": {
			input:       "2e3",
			destination: new(float32),
			expected:    float32(2e3),
		},
		"float64": {
			input:       "2e3",
			destination: new(float64),
			expected:    2e3,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v := NewParser([]byte(tc.input)).Parse()
			assert.NoError(t, NewDecoder().Unmarshal(v, tc.destination))

			// destVal stands the value of pointer underlying tc.destination interface.
			destVal := reflect.ValueOf(tc.destination).Elem().Interface()
			assert.Equal(t, tc.expected, destVal)
		})
	}
}

func Test_Unmarshal_String_and_Slice(t *testing.T) {
	makeslice := func(len int) *[]byte {
		sl := make([]byte, len)
		return &sl
	}
	testCases := map[string]struct {
		input       string
		destination interface{}
		expected    interface{}
	}{
		"String": {
			input:       `"helloworld"`,
			destination: new(string),
			expected:    "helloworld",
		},
		"Bytes Slice With Enough Space": {
			input:       `"hello-world"`,
			destination: makeslice(12),
			expected:    []byte("hello-world"),
		},
		"Bytes Slice Without Enough Space": {
			input:       `"hello-world"`,
			destination: makeslice(2),
			expected:    []byte("hello-world"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v := NewParser([]byte(tc.input)).Parse()
			assert.NoError(t, NewDecoder().Unmarshal(v, tc.destination))

			// destVal stands the value of pointer underlying tc.destination interface.
			destVal := reflect.ValueOf(tc.destination).Elem().Interface()
			assert.Equal(t, tc.expected, destVal)
		})
	}
}

func Test_Unmarshal_Array(t *testing.T) {
	var ar [3]byte
	v := NewParser([]byte(`"hello-world"`)).Parse()
	assert.NoError(t, NewDecoder().Unmarshal(v, &ar))

	expected := [3]byte{'h', 'e', 'l'}
	assert.Equal(t, expected, ar)
}

func Test_Unmarshal_Null(t *testing.T) {
	v := NewParser([]byte("null")).Parse()

	i := 1
	var ip = &i
	var ei *int = nil
	assert.NoError(t, NewDecoder().Unmarshal(v, &ip))
	assert.Equal(t, ei, ip)

	tr := true
	var bp = &tr
	assert.NoError(t, NewDecoder().Unmarshal(v, &bp))
	var bi *bool = nil
	assert.Equal(t, bi, bp)

	type tmp struct{}
	var tmpp = &tmp{}
	assert.NoError(t, NewDecoder().Unmarshal(v, &tmpp))
	var etmp *tmp = nil
	assert.Equal(t, etmp, tmpp)
}

func Test_Unmarshal_Bool(t *testing.T) {
	var b bool
	assert.NoError(t, NewDecoder().Unmarshal(NewParser([]byte("true")).Parse(), &b))
	assert.True(t, b)
	assert.NoError(t, NewDecoder().Unmarshal(NewParser([]byte("false")).Parse(), &b))
	assert.False(t, b)
}

func Test_Simple_JsonArray_To_Slice(t *testing.T) {
	i1 := make([]int, 1)
	i3 := make([]int, 3)
	i4 := make([]int, 4)
	testCases := map[string]struct {
		input    string
		expected interface{}
		dest     interface{}
	}{
		"int slice": {
			input:    "[1,2,3]",
			expected: []int{1, 2, 3},
			dest:     new([]int),
		},
		"int slice with allocation 1": {
			input:    "[1,2,3]",
			expected: []int{1, 2, 3},
			dest:     &i1,
		},
		"int slice with allocation 3": {
			input:    "[1,2,3]",
			expected: []int{1, 2, 3},
			dest:     &i3,
		},
		"int slice with allocation 4": {
			input:    "[1,2,3]",
			expected: []int{1, 2, 3, 0},
			dest:     &i4,
		},
		"float slice": {
			input:    "[-1.2,2.2,3.2]",
			expected: []float64{-1.2, 2.2, 3.2},
			dest:     new([]float64),
		},
		"bool slice": {
			input:    "[true,false,true]",
			expected: []bool{true, false, true},
			dest:     new([]bool),
		},
		"null slice": {
			input:    "[null,null]",
			expected: []*int{nil, nil},
			dest:     new([]*int),
		},
		"slice int slice": {
			input:    "[[1],[2],[3]]",
			expected: [][]int{{1}, {2}, {3}},
			dest:     new([][]int),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, NewDecoder().Unmarshal(NewParser([]byte(tc.input)).Parse(), tc.dest))
			assert.Equal(t, tc.expected, reflect.ValueOf(tc.dest).Elem().Interface())
		})
	}
}

func Test_Simple_JsonArray_To_Array(t *testing.T) {
	var i1 [1]int
	var i3 [3]int
	var i4 [4]int

	var f1 [1]float64
	var f3 [3]float64
	var f4 [4]float64

	var b1 [1]bool
	var b3 [3]bool
	var b4 [4]bool

	var null1 [1]*int
	var aa [2][1]int

	testCases := map[string]struct {
		input    string
		expected interface{}
		dest     interface{}
	}{
		"int array": {
			input:    "[1,2,3]",
			expected: [1]int{1},
			dest:     &i1,
		},
		"int array with allocation 3": {
			input:    "[1,2,3]",
			expected: [3]int{1, 2, 3},
			dest:     &i3,
		},
		"int array with allocation 4": {
			input:    "[1,2,3]",
			expected: [4]int{1, 2, 3, 0},
			dest:     &i4,
		},
		"float array with allocation 1": {
			input:    "[-1.2,2.2,3.2]",
			expected: [1]float64{-1.2},
			dest:     &f1,
		},
		"float array with allocation 3": {
			input:    "[-1.2,2.2,3.2]",
			expected: [3]float64{-1.2, 2.2, 3.2},
			dest:     &f3,
		},
		"float array with allocation 4": {
			input:    "[-1.2,2.2,3.2]",
			expected: [4]float64{-1.2, 2.2, 3.2, 0},
			dest:     &f4,
		},
		"bool array with allocation 1": {
			input:    "[true,false,true]",
			expected: [1]bool{true},
			dest:     &b1,
		},
		"bool array with allocation 3": {
			input:    "[true,false,true]",
			expected: [3]bool{true, false, true},
			dest:     &b3,
		},
		"bool array with allocation 4": {
			input:    "[true,false,true]",
			expected: [4]bool{true, false, true, false},
			dest:     &b4,
		},
		"null array with allocation 1": {
			input:    "[null,null]",
			expected: [1]*int{nil},
			dest:     &null1,
		},
		"array int array": {
			input:    "[[1],[1,2],[3]]",
			expected: [2][1]int{{1}, {1}},
			dest:     &aa,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, NewDecoder().Unmarshal(NewParser([]byte(tc.input)).Parse(), tc.dest))
			assert.Equal(t, tc.expected, reflect.ValueOf(tc.dest).Elem().Interface())
		})
	}
}

func Test_Simple_Object(t *testing.T) {
	type (
		Sub struct {
			ArrayInt []int  `json:"array_int"`
			Str      string `json:"str"`
		}
		Nest struct {
			Hello string `json:"hello"`
		}
		demo struct {
			Str       string  `json:"str"`
			Int       int     `json:"int"`
			Float64   float64 `json:"float64"`
			Bool      bool    `json:"bool"`
			Null      *int    `json:"null"`
			ArrayInt  []int   `json:"array_int"`
			ArrayInt1 [1]int  `json:"array_int1"`
			Sub       Sub     `json:"sub"`
			Useless   int
			Nest
		}
	)

	var p *int
	expected := demo{
		Str:       "str",
		Int:       999,
		Float64:   0.99,
		Bool:      true,
		Null:      p,
		ArrayInt:  []int{-1, 0, 1},
		ArrayInt1: [1]int{-1},
		Sub: Sub{
			ArrayInt: []int{-1, 0, 1},
			Str:      "str",
		},
		Nest: Nest{Hello: "hello"},
	}
	jsonStr := `
{
  "str": "str",
  "int": 999,
  "float64": 0.99,
  "bool": true,
  "null": null,
  "array_int": [-1,0,1],
  "array_int1": [-1,0,1],
  "sub": {
    "array_int": [-1,0,1],
    "str": "str"
  },
  "hello": "hello"
}
`

	var d demo
	assert.NoError(t, NewDecoder().Unmarshal(NewParser([]byte(jsonStr)).Parse(), &d))
	assert.Equal(t, expected, d)
}
