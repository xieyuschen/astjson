package astjson

import (
	"strconv"
)

// NodeType represents an AST node type
type NodeType uint

const (
	// NtNumber denotes an AST node is a number node
	// todo: Nt prefix sounds messy. try to find a better name
	NtNumber NodeType = iota
	NtNull
	NtString
	NtBool
	NtObject
	NtArray
)

// Value is the concrete AST representation
type Value struct {
	tp NodeType

	// use assert with the help of tp
	// val should be a pointer
	val interface{}
}

type NumberAst float64
type NullAst struct{} // only type make snese for NullAst
type BoolAst bool
type StringAst string

type ObjectAst struct {
	m map[Value]Value
}

type ArrayAst struct {
	values []Value
}

func Parse(bs []byte) *Value {
	return NewParser(bs).Parse()
}

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
	case Number, String, Bool, Null:
		return literal(p.bs, tk)
	case EOF:
		return nil
	case ArrayStart:
		return p.arrayParser()
	case ObjectStart:
		return p.objectParser()
	default:
		panic("invalid json syntax")
	}
}

func (a *ArrayAst) verifyNextType(ntp NodeType) bool {
	if len(a.values) == 0 {
		return true
	}
	if a.values[len(a.values)-1].tp == ntp {
		return true
	}
	return false
}

func (p *Parser) arrayParser() *Value {
	var ar ArrayAst

	for {
		tk := p.nextExceptWhitespace()
		if tk.tp == ArrayEnd {
			return &Value{
				tp:  NtArray,
				val: &ArrayAst{},
			}
		}
		val := p.parse(tk)

		if ar.verifyNextType(val.tp) {
			ar.values = append(ar.values, *val)
		} else {
			panic("inconsistent array value type")
		}

		// check whether an array ends
		then := p.nextExceptWhitespace()
		if then.tp == ArrayEnd {
			break
		} else if then.tp == Comma {
			continue
		} else {
			panic("invalid token after colon")
		}
	}

	return &Value{
		tp:  NtArray,
		val: &ar,
	}
}

func (p *Parser) objectParser() *Value {
	var v ObjectAst
	v.m = map[Value]Value{}

	for {
		start := p.nextExceptWhitespace()
		// an object is empty {}
		if start.tp == ObjectEnd {
			return &Value{
				tp:  NtObject,
				val: &v,
			}
		}

		if start.tp != String {
			panic("Invalid json schema for key")
		}

		key := literal(p.bs, start)

		if Colon != p.nextExceptWhitespace().tp {
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
		if then.tp == ObjectEnd {
			break
		} else if then.tp == Comma {
			continue
		} else {
			panic("invalid token after colon")
		}
	}

	return &Value{
		tp:  NtObject,
		val: &v,
	}
}

func NewParser(bs []byte) *Parser {
	return &Parser{
		bs: bs,
		l:  newLexer(bs),
	}
}

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

func (p *Parser) nextExceptWhitespace() token {
	return p.next(WhiteSpace)
}

func literal(bs []byte, tk token) *Value {
	var v Value
	switch tk.tp {
	case String:
		v.tp = NtString
		// remove left and right "
		v.val = StringAst(bs[tk.leftPos+1 : tk.rightPos-1])
	case Bool:
		v.tp = NtBool
		b, _ := strconv.ParseBool(string(bs[tk.leftPos:tk.rightPos]))
		v.val = BoolAst(b)
	case Number:
		v.tp = NtNumber
		f, _ := strconv.ParseFloat(string(bs[tk.leftPos:tk.rightPos]), 64)
		v.val = NumberAst(f)
	case Null:
		// the val of those types are useless
		v.tp = NtNull
		v.val = &NullAst{}
	}
	return &v
}
