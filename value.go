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

//go:generate stringer -type=numberType
type numberType uint

const (
	floatNumber numberType = iota
	// unsignedInteger refers we could store the value inside an uint64
	unsignedInteger
	// integer refers we could store the value inside an int64
	integer
)

type NumberAst struct {
	nt numberType
	f  float64
	u  uint64
	i  int64
}

type NullAst struct{}
type BoolAst bool
type StringAst string

type ObjectAst struct {
	m map[Value]Value
}

type ArrayAst struct {
	values []Value
}
