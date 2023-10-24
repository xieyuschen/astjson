package astjson

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
