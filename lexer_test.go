package astjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Scan_Simple_Tokens(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected token
	}{
		"eof": {input: ``, expected: token{
			tp:       EOF,
			leftPos:  0,
			rightPos: 0,
		}},
		"string": {input: `"123"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 5,
		}},
		"string with backward slash": {input: `"\"123"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 7,
		}},
		`string with \"`: {input: `"\"123"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 7,
		}},
		`string with \\"`: {input: `"\\123"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 7,
		}},
		`string with \/"`: {input: `"\/123"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 7,
		}},
		`string with \b"`: {input: `"\b123"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 7,
		}},
		`string with \f"`: {input: `"\f123"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 7,
		}},
		`string with \n"`: {input: `"\n123"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 7,
		}},
		`string with \r"`: {input: `"\r123"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 7,
		}},
		`string with \t"`: {input: `"\t123"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 7,
		}},
		`string with \u1234"`: {input: `"\u1234"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 8,
		}},
		`string with \uabcd"`: {input: `"\uabcd"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 8,
		}},
		`string with \uffff"`: {input: `"\uffff"`, expected: token{
			tp:       String,
			leftPos:  0,
			rightPos: 8,
		}},
		"positive integer": {input: "999", expected: token{
			tp:       Number,
			leftPos:  0,
			rightPos: 3,
		}},
		"negative integer": {input: "-999", expected: token{
			tp:       Number,
			leftPos:  0,
			rightPos: 4,
		}},
		"positive float": {input: "0.99", expected: token{
			tp:       Number,
			leftPos:  0,
			rightPos: 4,
		}},
		"negative float": {input: "-0.99", expected: token{
			tp:       Number,
			leftPos:  0,
			rightPos: 5,
		}},
		"zero": {input: "0", expected: token{
			tp:       Number,
			leftPos:  0,
			rightPos: 1,
		}},
		"null": {input: `null`, expected: token{
			tp:       Null,
			leftPos:  0,
			rightPos: 4,
		}},
		"true": {input: `true`, expected: token{
			tp:       Bool,
			leftPos:  0,
			rightPos: 4,
		}},
		"false": {input: `false`, expected: token{
			tp:       Bool,
			leftPos:  0,
			rightPos: 5,
		}},
		"left {": {input: `{`, expected: token{
			tp:       ObjectStart,
			leftPos:  0,
			rightPos: 1,
		}},
		"right }": {input: "}", expected: token{
			tp:       ObjectEnd,
			leftPos:  0,
			rightPos: 1,
		}},
		"left [": {input: `[`, expected: token{
			tp:       ArrayStart,
			leftPos:  0,
			rightPos: 1,
		}},
		"right ]": {input: "]", expected: token{
			tp:       ArrayEnd,
			leftPos:  0,
			rightPos: 1,
		}},
		"colon :": {input: ":", expected: token{
			tp:       Colon,
			leftPos:  0,
			rightPos: 1,
		}},
		"comma ,": {input: ",", expected: token{
			tp:       Comma,
			leftPos:  0,
			rightPos: 1,
		},
		},
		"whitespace space": {input: " ", expected: token{
			tp:       WhiteSpace,
			leftPos:  0,
			rightPos: 1,
		},
		},
		"whitespace linefeed": {input: "\r", expected: token{
			tp:       WhiteSpace,
			leftPos:  0,
			rightPos: 1,
		},
		},
		"whitespace carriage return": {input: "\n", expected: token{
			tp:       WhiteSpace,
			leftPos:  0,
			rightPos: 1,
		},
		},
		"whitespace tab": {input: "\t", expected: token{
			tp:       WhiteSpace,
			leftPos:  0,
			rightPos: 1,
		},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			l := newLexer([]byte(tc.input))
			tk := l.Scan()
			assert.Equal(t, tc.expected, tk)
		})
	}
}

func Test_Scan_Multiple_Tokens(t *testing.T) {
	testCases := map[string]struct {
		input string
		// check token type here to ease writing test cases
		// check curPos and lastPos inside code
		expected []Type
	}{
		"object start and end":                  {input: "{}", expected: []Type{ObjectStart, ObjectEnd, EOF, EOF}},
		"object end and start":                  {input: "}{", expected: []Type{ObjectEnd, ObjectStart, EOF}},
		"object end and end":                    {input: "}}", expected: []Type{ObjectEnd, ObjectEnd, EOF}},
		"object start, space and end":           {input: "{ }", expected: []Type{ObjectStart, WhiteSpace, ObjectEnd, EOF}},
		"object start, linefeed and end":        {input: "{\r}", expected: []Type{ObjectStart, WhiteSpace, ObjectEnd, EOF}},
		"object start, carriage return and end": {input: "{\n}", expected: []Type{ObjectStart, WhiteSpace, ObjectEnd, EOF}},
		"object start, horizontal tab and end":  {input: "{\t}", expected: []Type{ObjectStart, WhiteSpace, ObjectEnd, EOF}},
		`{"str": string}`:                       {input: `{"str": "string"}`, expected: []Type{ObjectStart, String, Colon, WhiteSpace, String, ObjectEnd, EOF}},
		`{"str": 123}`:                          {input: `{"str": 123}`, expected: []Type{ObjectStart, String, Colon, WhiteSpace, Number, ObjectEnd, EOF}},
		`{"str": -123}`:                         {input: `{"str": -123}`, expected: []Type{ObjectStart, String, Colon, WhiteSpace, Number, ObjectEnd, EOF}},
		`{"str": 0}`:                            {input: `{"str": 0}`, expected: []Type{ObjectStart, String, Colon, WhiteSpace, Number, ObjectEnd, EOF}},
		`{"str": 0.99}`:                         {input: `{"str": 0.99}`, expected: []Type{ObjectStart, String, Colon, WhiteSpace, Number, ObjectEnd, EOF}},
		`{"str": -0.99}`:                        {input: `{"str": -0.99}`, expected: []Type{ObjectStart, String, Colon, WhiteSpace, Number, ObjectEnd, EOF}},
		`{"str": 123e456}`:                      {input: `{"str": 123e456}`, expected: []Type{ObjectStart, String, Colon, WhiteSpace, Number, ObjectEnd, EOF}},
		`{"str": 123-e456}`:                     {input: `{"str": 123e-456}`, expected: []Type{ObjectStart, String, Colon, WhiteSpace, Number, ObjectEnd, EOF}},
		`{"str": true}`:                         {input: `{"str": true}`, expected: []Type{ObjectStart, String, Colon, WhiteSpace, Bool, ObjectEnd, EOF}},
		`{"str": false}`:                        {input: `{"str": false}`, expected: []Type{ObjectStart, String, Colon, WhiteSpace, Bool, ObjectEnd, EOF}},
		`{"str": null}`:                         {input: `{"str": null}`, expected: []Type{ObjectStart, String, Colon, WhiteSpace, Null, ObjectEnd, EOF}},
		`{,}`:                                   {input: `{,}`, expected: []Type{ObjectStart, Comma, ObjectEnd, EOF}},
		`{,,}`:                                  {input: `{,,}`, expected: []Type{ObjectStart, Comma, Comma, ObjectEnd, EOF}},
		`{123,}`:                                {input: `{123,}`, expected: []Type{ObjectStart, Number, Comma, ObjectEnd, EOF}},
		`{1.234,}`:                              {input: `{1.234,}`, expected: []Type{ObjectStart, Number, Comma, ObjectEnd, EOF}},
		`{"123",}`:                              {input: `{"123",}`, expected: []Type{ObjectStart, String, Comma, ObjectEnd, EOF}},
		`[]`:                                    {input: `[]`, expected: []Type{ArrayStart, ArrayEnd, EOF}},
		`["1"]`:                                 {input: `["1"]`, expected: []Type{ArrayStart, String, ArrayEnd, EOF}},
		`[1]`:                                   {input: `[1]`, expected: []Type{ArrayStart, Number, ArrayEnd, EOF}},
		`[1.23]`:                                {input: `[1.23]`, expected: []Type{ArrayStart, Number, ArrayEnd, EOF}},
		`[-1.23]`:                               {input: `[-1.23]`, expected: []Type{ArrayStart, Number, ArrayEnd, EOF}},
		`[1,2]`:                                 {input: `[1,2]`, expected: []Type{ArrayStart, Number, Comma, Number, ArrayEnd, EOF}},
		`"\ufffff"`:                             {input: `"\ufffff"`, expected: []Type{String, EOF}},
		`"\uffffg"`:                             {input: `"\uffffg"`, expected: []Type{String, EOF}},
		`"\uffff\uffff"`:                        {input: `"\uffff\uffff"`, expected: []Type{String, EOF}},
		`"\"\/\\\b\f\n\r\t\uabcd"`:              {input: `"\"\/\\\b\f\n\r\t\uabcd"`, expected: []Type{String, EOF}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			l := newLexer([]byte(tc.input))
			var counter, lastPos int
			var tk token

			for tk.tp != EOF {
				assert.LessOrEqual(t, counter, len(tc.expected))
				tk = l.Scan()

				assert.Equal(t, tc.expected[counter].String(), tk.tp.String(), "token types don't match")
				assert.Equal(t, lastPos, tk.leftPos, "positions are wrong")
				t.Logf("validation from [%d,%d) has passed for token type: %s",
					lastPos, tk.rightPos, tk.tp.String())
				lastPos = tk.rightPos
				counter++

			}
		})
	}
}

func Test_Scan_Panic(t *testing.T) {
	testCases := map[string]struct {
		input string
	}{
		`invalid string \d`:     {input: `"\d"`},
		`invalid string \uabcg`: {input: `"\uabcg"`},
		`invalid string "abc`:   {input: `"abc`},
		`invalid bool truu`:     {input: "truu"},
		`invalid bool falss `:   {input: "falss"},
		`invalid null nul `:     {input: "nul"},
		`invalid null nul1 `:    {input: "nul1"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Panics(t, func() {
				l := newLexer([]byte(tc.input))
				_ = l.Scan()
			})
		})
	}

}
