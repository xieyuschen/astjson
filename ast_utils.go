package astjson

func IsBool(value *Value) bool {
	return value.NodeType == Bool
}

func IsNull(value *Value) bool {
	return value.NodeType == Null
}

func IsNumber(value *Value) bool {
	return value.NodeType == Number
}

func IsString(value *Value) bool {
	return value.NodeType == String
}

func IsArray(value *Value) bool {
	return value.NodeType == Array
}

func IsObject(value *Value) bool {
	return value.NodeType == Object
}

func GetArrayValues(value *Value) []Value {
	if !IsArray(value) {
		panic("value is not an object")
	}
	return value.AstValue.(*ArrayAst).Values
}

func GetObjectKvMap(value *Value) map[string]Value {
	if !IsObject(value) {
		panic("value is not an object")
	}
	return value.AstValue.(*ObjectAst).KvMap
}

func GetString(value *Value) string {
	if !IsString(value) {
		panic("value is not a string")
	}
	return string(value.AstValue.(StringAst))
}

func GetBool(value *Value) bool {
	if !IsBool(value) {
		panic("value is not a bool")
	}
	return bool(value.AstValue.(BoolAst))
}

func GetNumber(value *Value) NumberAst {
	if !IsNumber(value) {
		panic("value is not a number")
	}
	return value.AstValue.(NumberAst)
}
