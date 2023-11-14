package astjson

// NodeType represents an AST node type
//
//go:generate stringer -type=NodeType
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
type NumberType uint

const (
	floatNumber NumberType = iota
	// unsignedInteger refers we could store the value inside an uint64
	unsignedInteger
	// integer refers we could store the value inside an int64
	integer
)

type NumberAst struct {
	Nt NumberType
	f  float64
	u  uint64
	i  int64
}

// GetInt64 returns the KvMap inside a NumberAst, it's possible to
// lose precise for float64 or overflow for uint64
func (n NumberAst) GetInt64() int64 {
	switch n.Nt {
	case integer:
		return n.i
	case unsignedInteger:
		return int64(n.u)
	case floatNumber:
		// precise is acceptable because users need us to cast it.
		return int64(n.f)
	}
	panic("")
}

func (n NumberAst) GetUint64() uint64 {
	switch n.Nt {
	case integer:
		return uint64(n.i)
	case unsignedInteger:
		return n.u
	case floatNumber:
		// precise is acceptable because users need us to cast it.
		return uint64(n.f)
	}
	panic("")
}

func (n NumberAst) GetFloat64() float64 {
	switch n.Nt {
	case integer:
		return float64(n.i)
	case unsignedInteger:
		return float64(n.u)
	case floatNumber:
		// todo: check further whether this logic is correct
		return n.f
	}
	panic("")
}

type NullAst struct{}
type BoolAst bool
type StringAst string

type ObjectAst struct {
	KvMap map[string]Value
}

type ArrayAst struct {
	Values []Value
}
