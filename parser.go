package astjson

import (
	"strconv"
)

// Value is the concrete AST representation
type Value struct {
	// todo: which name is better? tp or op
	tp Type

	raw []byte
	// use assert with the help of tp
	val interface{}
}

type NumberAst float64
type NullAst struct{} // only type make snese for NullAst
type BoolAst bool
type StringAst string

func Parse(bs []byte) Value {
	l := newLexer(bs)
	tk := l.Scan()

	v := Value{
		tp: tk.tp,
	}

	// same leftPos and rightPos will get an array hold a nil
	// instead of a nil array
	if l.lastPos != l.curPos {
		v.raw = bs[tk.leftPos:tk.rightPos]
	}
	switch tk.tp {
	case String:
		v.val = StringAst(v.raw)
		return v
	case Bool:
		b, _ := strconv.ParseBool(string(v.raw))
		v.val = BoolAst(b)
		return v

	case Number:
		f, _ := strconv.ParseFloat(string(v.raw), 64)
		v.val = NumberAst(f)
		return v
	case EOF, Null:
		// the val of those types are useless
	}
	return v
}
