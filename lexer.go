package astjson

import (
	"fmt"
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
)

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
	if l.curPos == len(l.bs) {
		return token{tp: EOF}
	}
	
	c := l.bs[l.curPos]
	// align sentries
	l.lastPos = l.curPos
	switch c {
	case '"':
		// string case
		// todo: support backward slash in the future
		return l.stringType()
	case 'f', 't':
		// bool case
		return l.boolType()
	case 'n':
		// null case
		return l.nullType()
	default:
		// number case
		return l.numberType()
	}
}

// todo: if a string is longer than 20 chars after ", we panic now.
const maxStringLength = 10

func (l *lexer) stringType() token {
	var counter int
	for l.curPos < len(l.bs) {
		counter++
		if counter > maxStringLength {
			panic(fmt.Sprintf("excessive string from %d to %d", l.lastPos, l.curPos))
		}
		
		l.curPos++
		if l.bs[l.curPos] != '"' {
			continue
		}
		
		return token{
			tp:      String,
			leftPos: l.lastPos,
			// the curPos ends at where the second " occurs
			rightPos: l.curPos + 1,
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
	
	panic("not a valid json string")
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
	
	panic("not a valid json string")
}

func (l *lexer) numberType() token {
	switch l.bs[l.curPos] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		l.curPos++
		for ; l.curPos < len(l.bs); l.curPos++ {
			switch l.bs[l.curPos] {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
				'.', 'e', 'E', '+', '-':
			default:
				break
			}
		}
	}
	
	return token{
		tp:       Number,
		leftPos:  l.lastPos,
		rightPos: l.curPos,
	}
}
