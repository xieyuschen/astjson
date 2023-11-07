# astjson
[![codecov](https://codecov.io/gh/xieyuschen/astjson/graph/badge.svg?token=FFXVZQYUWF)](https://codecov.io/gh/xieyuschen/astjson)

## Requirement
Astjson requires at least `go1.18` as it uses some features from reflect library.

```shell
go get github.com/xieyuschen/astjson
```

## Example

```go
package main

import (
	"errors"
	"fmt"

	"github.com/xieyuschen/astjson"
)

const (
	jsonStr = `
{
  "num": 999,
  "str": "helloworld",
  "enabled": true ,
  "sub1": {
    "sub2": {
      "key": 999
    }
  },
  "name": "astjson"
}`
)

type (
	Sub1 struct {
		Sub2 Sub2 `json:"sub2"`
	}
	Sub2 struct {
		Key int `json:"key"`
	}
	demo struct {
		Num     int    `json:"num"`
		Str     string `json:"str"`
		Enabled bool   `json:"enabled"`
		Sub     Sub1   `json:"sub1"`
	}
)

func main() {

	astValue := astjson.NewParser([]byte(jsonStr)).Parse()

	val, err := astjson.NewWalker(astValue).
		// field "str" is required
		Field("str").
		// field "num" is optional, but if it exists the validator will be triggered
		Optional("num", equal999).
		Field("enabled").Validate(isTrue).
		Path("sub1").Path("sub2").Validate(equal999).EndPath().
		ValidateKey("name", hasAstjsonName).
		Walk()

	if err != nil {
		panic(err)
	}
	var d demo
	_ = astjson.NewDecoder().Unmarshal(val, &d)
	fmt.Println(d)
}

func equal999(value *astjson.Value) error {
	actual := value.AstValue.(astjson.NumberAst).GetInt64()
	if actual == 999 {
		return nil
	}
	return errors.New("num should be 999")
}

func isTrue(value *astjson.Value) error {
	if bool(value.AstValue.(astjson.BoolAst)) {
		return nil
	}
	return errors.New("bool should be true")
}
func hasAstjsonName(value *astjson.Value) error {
	if value.AstValue.(astjson.StringAst) == "astjson" {
		return nil
	}
	return errors.New("name should be astjson")
}
```

See more [examples here](astjson_example_test.go).

## Motivation
Haskell json library [aeson](https://github.com/haskell/aeson) parsed AST first in its early version as a default way.
However, it skips to convert AST first to speed up parsing as a new default way.

Parsing AST first is slower than directly parsing json, however, it provides more flexibilities.
Hence, inspired by aeson, I wrote a ast json to help parsing json to AST.

- [json tokens](https://www.json.org/json-en.html)
- [aeson](https://hackage.haskell.org/package/aeson-2.2.1.0/docs/Data-Aeson.html)
