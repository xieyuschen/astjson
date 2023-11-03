# astjson
[![codecov](https://codecov.io/gh/xieyuschen/astjson/graph/badge.svg?token=FFXVZQYUWF)](https://codecov.io/gh/xieyuschen/astjson)

## Requirement
Astjson requires at least `go1.18` as it uses some features from reflect library.

```shell
go get github.com/xieyuschen/astjson
```

## Usage

```go
package main

import (
	"errors"
	"fmt"

	"github.com/xieyuschen/astjson"
)

func main() {
	jsonStr := `{ "num": 999, "str": "helloworld", "enabled": true , "name": "astjson" }`
	astValue := astjson.NewParser([]byte(jsonStr)).Parse()

	val, err := astjson.NewWalker(astValue).
		// field "str" is required
		Field("str").
		// field "num" is optional, but if it exists the validator will be triggered
		Optional("num", func(value *astjson.Value) error {
			actual := value.AstValue.(astjson.NumberAst).GetInt64()
			if actual == 999 {
				return nil
			}
			return errors.New("num should be 999")
		}).
		Field("enabled").
		Validate(func(value *astjson.Value) error {
            if bool(value.AstValue.(astjson.BoolAst)) {
                return nil
            }
            return errors.New("bool should be true")
	    }).
		ValidateKey("name", func(value *astjson.Value) error {
			if value.AstValue.(astjson.StringAst)=="astjson"{
				return nil
            }
			return errors.New("name should be astjson")
		}).
		Walk()

	fmt.Println(val, err)
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
