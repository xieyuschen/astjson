package astjson

import (
	"errors"
	"fmt"
)

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
		Nest
	}
)

const jsonStr = `
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

func ExampleDecoder_Unmarshal() {
	var d demo
	err := NewDecoder().Unmarshal(NewParser([]byte(jsonStr)).Parse(), &d)
	dieIf(err)
	fmt.Println(d)
}

func ExampleWalker_Walk() {

	astValue := NewParser([]byte(jsonStr)).Parse()

	val, err := NewWalker(astValue).
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
		ValidateKey("null", func(value *Value) error {
			return nil
		}).
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
	dieIf(err)

	var d demo
	err = NewDecoder().Unmarshal(val, &d)
	dieIf(err)
	fmt.Println(d)
}

func dieIf(err error) {
	if err != nil {
		panic(err)
	}
}
