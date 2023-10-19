package astjson

import (
	"strconv"
)

// NodeType represents an AST node type
type NodeType uint

const (
	Number NodeType = iota
	Null
	String
	Bool
	Object
	Array
)

// Value is the concrete AST representation
type Value struct {
	NodeType
	
	// AstValue stores the real value of an AST Value
	// the literal type(Number, String, Bool and Null) stores them by value
	// the Object and Array stores them by pointer
	AstValue interface{}
}

type NumberAst float64
type NullAst struct{}
type BoolAst bool
type StringAst string

type ObjectAst struct {
	m map[Value]Value
}

type ArrayAst struct {
	values []Value
}

// Parse transforms the json bytes to AST value, it will return nil or panic
// when the input is invalid.
// todo: fix me: return error instead of panic
func Parse(bs []byte) *Value {
	return NewParser(bs).Parse()
}

// Parser helps to parse the bytes to AST value
type Parser struct {
	bs []byte
	l  *lexer
}

func (p *Parser) Parse() *Value {
	p.l.Reset()
	tk := p.nextExceptWhitespace()
	return p.parse(tk)
}

// parse helps to get a whole object, array or a literal type.
func (p *Parser) parse(tk token) *Value {
	switch tk.tp {
	case tkNumber, tkString, tkBool, tkNull:
		return literal(p.bs, tk)
	case tkEOF:
		return nil
	case tkArrayStart:
		return p.arrayParser()
	case tkObjectStart:
		return p.objectParser()
	default:
		panic("invalid json syntax")
	}
}

// verifyNextType verifies whether the next ntp node type satisfies
// the array type. It returns true when the array is empty or the type is same
// with the last element.
// because we always verify the type before appending the array, it's safe to
// compare the tail element.
func (a *ArrayAst) verifyNextType(ntp NodeType) bool {
	if len(a.values) == 0 {
		return true
	}
	if a.values[len(a.values)-1].NodeType == ntp {
		return true
	}
	return false
}

// arrayParser parses the remained part of an array after tkArrayStart is found before.
func (p *Parser) arrayParser() *Value {
	var ar ArrayAst
	
	for {
		tk := p.nextExceptWhitespace()
		if tk.tp == tkArrayEnd {
			return &Value{
				NodeType: Array,
				AstValue: &ArrayAst{},
			}
		}
		val := p.parse(tk)
		
		if ar.verifyNextType(val.NodeType) {
			ar.values = append(ar.values, *val)
		} else {
			panic("inconsistent array value type")
		}
		
		// check whether an array ends
		then := p.nextExceptWhitespace()
		if then.tp == tkArrayEnd {
			break
		} else if then.tp == tkComma {
			continue
		} else {
			panic("invalid token after colon")
		}
	}
	
	return &Value{
		NodeType: Array,
		AstValue: &ar,
	}
}

// objectParser parses the remained part of an array after tkObjectStart is found before.
func (p *Parser) objectParser() *Value {
	var v ObjectAst
	v.m = map[Value]Value{}
	
	for {
		start := p.nextExceptWhitespace()
		// an object is empty {}
		if start.tp == tkObjectEnd {
			return &Value{
				NodeType: Object,
				AstValue: &v,
			}
		}
		
		if start.tp != tkString {
			panic("Invalid json schema for key")
		}
		
		key := literal(p.bs, start)
		
		if tkColon != p.nextExceptWhitespace().tp {
			panic("invalid json schema after key")
		}
		if _, ok := v.m[*key]; ok {
			panic("duplicated key")
		}
		
		val := p.parse(p.nextExceptWhitespace())
		v.m[*key] = *val
		
		// check whether an object ends
		// todo: refine me: the logic here is duplicated with the beginning of the for loop
		then := p.nextExceptWhitespace()
		if then.tp == tkObjectEnd {
			break
		} else if then.tp == tkComma {
			continue
		} else {
			panic("invalid token after colon")
		}
	}
	
	return &Value{
		NodeType: Object,
		AstValue: &v,
	}
}

// NewParser creates a new Parser to parse full json bytes to AST node.
func NewParser(bs []byte) *Parser {
	return &Parser{
		bs: bs,
		l:  newLexer(bs),
	}
}

// next keep retrieving tokens and return the token which type is not contained inside skips.
func (p *Parser) next(skips ...Type) token {
	shouldSkip := func(tk Type) bool {
		for _, skip := range skips {
			if tk == skip {
				return true
			}
		}
		return false
	}
	
	tk := p.l.Scan()
	for shouldSkip(tk.tp) {
		tk = p.l.Scan()
	}
	return tk
}

// nextExceptWhitespace returns the token which is not a tkWhiteSpace type.
func (p *Parser) nextExceptWhitespace() token {
	return p.next(tkWhiteSpace)
}

// literal constructs the AST value for Number, String, Bool and Null type.
// The AstValue inside Value is not a pointer.
func literal(bs []byte, tk token) *Value {
	var v Value
	switch tk.tp {
	case tkString:
		v.NodeType = String
		// remove left and right "
		v.AstValue = StringAst(bs[tk.leftPos+1 : tk.rightPos-1])
	case tkBool:
		v.NodeType = Bool
		b, _ := strconv.ParseBool(string(bs[tk.leftPos:tk.rightPos]))
		v.AstValue = BoolAst(b)
	case tkNumber:
		v.NodeType = Number
		f, _ := strconv.ParseFloat(string(bs[tk.leftPos:tk.rightPos]), 64)
		v.AstValue = NumberAst(f)
	case tkNull:
		// the AstValue of those types are useless
		v.NodeType = Null
		v.AstValue = &NullAst{}
	}
	return &v
}
