package astjson

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func Test_Non_Pointer(t *testing.T) {
	var i int
	err := NewDecoder().Unmarshal(NewParser([]byte("")).Parse(), i)
	assert.NotNil(t, err)
	assert.Equal(t, "dest must be a pointer", err.Error())
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
		"int array": {
			input:    "[1,2,3]",
			expected: []int{1, 2, 3},
			dest:     new([]int),
		},
		"int array with allocation 1": {
			input:    "[1,2,3]",
			expected: []int{1, 2, 3},
			dest:     &i1,
		},
		"int array with allocation 3": {
			input:    "[1,2,3]",
			expected: []int{1, 2, 3},
			dest:     &i3,
		},
		"int array with allocation 4": {
			input:    "[1,2,3]",
			expected: []int{1, 2, 3, 0},
			dest:     &i4,
		},
		"float array": {
			input:    "[-1.2,2.2,3.2]",
			expected: []float64{-1.2, 2.2, 3.2},
			dest:     new([]float64),
		},
		"bool array": {
			input:    "[true,false,true]",
			expected: []bool{true, false, true},
			dest:     new([]bool),
		},
		"null array": {
			input:    "[null,null]",
			expected: []*int{nil, nil},
			dest:     new([]*int),
		},
		"array int array": {
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
