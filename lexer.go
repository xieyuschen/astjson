package astjson

import (
	"fmt"
	"strconv"
)

// Type represents the token type
//
//go:generate stringer -type=Type
type Type uint

const (
	tkWhiteSpace Type = iota
	tkString
	tkNumber
	tkBool
	tkNull
	tkEOF
	tkObjectStart
	tkObjectEnd
	tkArrayStart
	tkArrayEnd
	tkComma
	tkColon
)

// token represents the json token.
// currently, token only supports to the limited json value and limits primitive
// types only.
type token struct {
	tp Type

	// the token value is [ leftPos, rightPos)
	// index starts at 0
	leftPos, rightPos int

	// hasDash and isFloat only make sense for tkNumber because we don't
	// want to lose precise
	hasDash, isFloat bool
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

func (l *lexer) Reset() {
	l.curPos, l.lastPos = 0, 0
}

// Scan returns one token or panic
// todo: return error instead of panic
func (l *lexer) Scan() token {
	// align sentries
	l.lastPos = l.curPos

	if l.curPos == len(l.bs) {
		return token{
			tp:       tkEOF,
			leftPos:  l.curPos,
			rightPos: l.curPos,
		}
	}

	c := l.bs[l.curPos]

	switch c {
	case '{':
		l.curPos += 1
		return token{
			tp:       tkObjectStart,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	case '}':
		l.curPos += 1
		return token{
			tp:       tkObjectEnd,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	case '[':
		l.curPos += 1
		return token{
			tp:       tkArrayStart,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	case ']':
		l.curPos += 1
		return token{
			tp:       tkArrayEnd,
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
			tp:       tkWhiteSpace,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	case ':':
		l.curPos += 1
		return token{
			tp:       tkColon,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	case ',':
		l.curPos += 1
		return token{
			tp:       tkComma,
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
			tp:      tkString,
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
			tp:       tkBool,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	}

	if string(l.bs[l.lastPos:l.curPos+len("false")]) == "false" {
		l.curPos += len("false")
		return token{
			tp:       tkBool,
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
			tp:       tkNull,
			leftPos:  l.lastPos,
			rightPos: l.curPos,
		}
	}

	panic("not a valid null value")
}

func (l *lexer) numberType() token {
	t := token{
		tp:      tkNumber,
		leftPos: l.lastPos,
	}
	switch l.bs[l.curPos] {
	case '-':
		t.hasDash = true
		fallthrough
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.curPos++
	Loop:
		for ; l.curPos < len(l.bs); l.curPos++ {
			switch l.bs[l.curPos] {
			case '.', 'e', 'E':
				t.isFloat = true
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-':
			default:
				break Loop
			}
		}
	}

	t.rightPos = l.curPos
	return t
}
