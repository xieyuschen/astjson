package astjson

import (
	"errors"
	"reflect"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

// Unmarshal decodes the AST to a structure.
// This API is an EXPERIENTIAL one and might be removed in the future.
func (d *Decoder) Unmarshal(val *Value, dest interface{}) error {
	if !isPointer(dest) {
		return errors.New("dest must be a pointer")
	}
	switch val.NodeType {
	case Number:
		return setNumber(val, dest)
	case String:
		return setString(val, dest)
	case Null:
		return setNull(val, dest)
	case Bool:
		return setBool(val, dest)
	}
	return nil
}

func setNull(val *Value, dest interface{}) error {
	// pointer owns type, we cannot assign it a nil directly
	v := reflect.ValueOf(dest).Elem()
	v.Set(reflect.Zero(v.Type()))
	return nil
}

func setBool(val *Value, dest interface{}) error {
	reflect.ValueOf(dest).Elem().SetBool(bool(val.AstValue.(BoolAst)))
	return nil
}

func setNumber(val *Value, dest interface{}) error {
	kind := reflect.ValueOf(dest).Elem().Kind()
	numberAst := val.AstValue.(NumberAst)
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		reflect.ValueOf(dest).Elem().SetInt(numberAst.getInt64())
		return nil
	
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		reflect.ValueOf(dest).Elem().SetUint(numberAst.getUint64())
		return nil
	
	case reflect.Float32,
		reflect.Float64:
		reflect.ValueOf(dest).Elem().SetFloat(numberAst.getFloat64())
		return nil
	}
	panic("fail to set number")
}

// setString set the string value into dest
// todo: support []byte and []int8
// todo: think about whether support the implicitly cast from byte to the other types
func setString(val *Value, dest interface{}) error {
	v := reflect.ValueOf(dest).Elem()
	kind := v.Kind()
	
	strAst := val.AstValue.(StringAst)
	switch kind {
	case reflect.String:
		v.SetString(string(strAst))
		return nil
	case reflect.Slice:
		v.SetBytes([]byte(strAst))
		return nil
	case reflect.Array:
		l := v.Len()
		bs := []byte(strAst)
		for i := 0; i < l; i++ {
			v.Index(i).Set(reflect.ValueOf(bs[i]))
		}
		
		return nil
	}
	panic("fail to set string")
}

func isPointer(dest interface{}) bool {
	if reflect.ValueOf(dest).Kind() != reflect.Pointer {
		return false
	}
	return true
}
