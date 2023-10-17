package astjson

import (
	"fmt"
	"strconv"
)

// Type represents the token type
type Type uint

const (
	WhiteSpace Type = iota
	String
	Number
	Bool
	Null
	EOF
	ObjectStart
	ObjectEnd
	ArrayStart
	ArrayEnd
	Comma
	Colon
)

var roster = map[Type]string{
	WhiteSpace:  "WhiteSpace",
	String:      "String",
	Number:      "Number",
	Bool:        "Bool",
	Null:        "Null",
	EOF:         "EOF",
	ObjectStart: "ObjectStart",
	ObjectEnd:   "ObjectEnd",
	Comma:       "Comma",
	Colon:       "Colon",
}

func (t Type) String() string {
	return roster[t]
}

// var lexicalRule = map[int]string{
// 	WhiteSpace: "WhiteSpace",
// }

// token represents the json token.
// currently, token only supports to the limited json value and limits primitive
// types only.
type token struct {
	tp Type

	// the token value is [ leftPos, rightPos)
	// index starts at 0
	leftPos, rightPos int
}

type lexer struct {
	bs []byte

	// todo: try to use uint
	curPos  int
	lastPos int
}

func newLexer(bs []byte) *lexer {
	return &lexer{
		bs:      bs,
		curPos:  0,
		lastPos: 0,
	}
}

// Scan returns one token or panic
// todo: return error instead of panic
func (l *lexer) Scan() token {
	// align sentries
	l.lastPos = l.curPos

	if l.curPos == len(l.bs) {
		return token{
			tp:       EOF,
			leftPos:  l.curPos,
			rightPos: l.curPos,
		}
	}

	c := l.bs[l.curPos]

	switch c {
	case '{':
		l.curPos += 1
		return token{
			tp:       ObjectStart,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	case '}':
		l.curPos += 1
		return token{
			tp:       ObjectEnd,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	case '[':
		l.curPos += 1
		return token{
			tp:       ArrayStart,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	case ']':
		l.curPos += 1
		return token{
			tp:       ArrayEnd,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	case '"':
		// string case
		return l.stringType()
	case 'f', 't':
		// bool case
		return l.boolType()
	case 'n':
		// null case
		return l.nullType()
	case ' ', '\t', '\n', '\r':
		l.curPos += 1
		return token{
			tp:       WhiteSpace,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	case ':':
		l.curPos += 1
		return token{
			tp:       Colon,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	case ',':
		l.curPos += 1
		return token{
			tp:       Comma,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	default:
		// number case
		return l.numberType()
	}
}

func (l *lexer) stringType() token {
	// move next to the starting "
	l.curPos++

	for l.curPos < len(l.bs) {
		if l.bs[l.curPos] == '\\' {
			l.curPos++
			switch l.bs[l.curPos] {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				l.curPos++
				continue
			case 'u':
				// u1234: check whether it's a hex digital
				s := l.bs[l.curPos+1 : l.curPos+5]
				_, err := strconv.ParseUint(string(s), 16, 64)
				if err != nil {
					panic(fmt.Errorf("invalid hex string at %d", l.curPos))
				}
				l.curPos += 5
			default:
				panic(fmt.Sprintf("invalid string \\ near %d", l.curPos))
			}
		}
		if l.bs[l.curPos] != '"' {
			l.curPos++
			continue
		}

		// move curPos right because we need to conclude " as wel
		l.curPos++
		return token{
			tp:      String,
			leftPos: l.lastPos,
			// the curPos ends at where the second " occurs
			rightPos: l.curPos,
		}
	}
	panic(fmt.Sprintf("invalid string from %d to %d", l.lastPos, l.curPos))
}

// todo: return error instead of panic
func (l *lexer) boolType() token {
	if string(l.bs[l.lastPos:l.curPos+len("true")]) == "true" {
		l.curPos += len("true")
		return token{
			tp:       Bool,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	}

	if string(l.bs[l.lastPos:l.curPos+len("false")]) == "false" {
		l.curPos += len("false")
		return token{
			tp:       Bool,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	}

	panic("not a valid json bool type")
}

// todo: return error instead of panic
func (l *lexer) nullType() token {
	l.curPos += 4
	str := string(l.bs[l.lastPos:l.curPos])

	if str == "null" {
		return token{
			tp:       Null,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	}

	panic("not a valid null value")
}

func (l *lexer) numberType() token {
	switch l.bs[l.curPos] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		l.curPos++
	Loop:
		for ; l.curPos < len(l.bs); l.curPos++ {
			switch l.bs[l.curPos] {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
				'.', 'e', 'E', '+', '-':
			default:
				break Loop
			}
		}
	}

	return token{
		tp:       Number,
		leftPos:  l.lastPos,
		rightPos: l.curPos,
	}
}
