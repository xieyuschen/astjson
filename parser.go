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
	val interface{}
}

type NumberAst float64
type NullAst struct{} // only type make snese for NullAst
type BoolAst bool
type StringAst string

func literal(bs []byte, tk token) *Value {
	var v Value
	rawStr := string(bs[tk.leftPos:tk.rightPos])
	switch tk.tp {
	case String:
		v.tp = NtString
		v.val = StringAst(rawStr)
	case Bool:
		v.tp = NtBool
		b, _ := strconv.ParseBool(rawStr)
		v.val = BoolAst(b)
	case Number:
		v.tp = NtNumber
		f, _ := strconv.ParseFloat(rawStr, 64)
		v.val = NumberAst(f)
	case Null:
		// the val of those types are useless
		v.tp = NtNull
	}
	return &v
}

func Parse(bs []byte) *Value {
	l := newLexer(bs)

	tk := l.Scan()
	switch tk.tp {
	case Number, String, Bool, Null:
		return literal(bs, tk)
	case EOF:
		return nil
	default:
		// todo: further impl
		return nil
	}
}
