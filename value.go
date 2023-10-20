package astjson

import (
	"fmt"
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

func (v *Value) String() string {
	switch v.NodeType {
	case Number:
		return fmt.Sprintf("%f", v.AstValue)
	case String:
		return v.AstValue.(string)
	case Bool:
		return strconv.FormatBool(v.AstValue.(bool))
	case Null:
		return "null"
	}
	return ""
}
