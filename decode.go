package astjson

import (
	"errors"
	"reflect"
)

const (
	jsonTAG = "json"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	d := &Decoder{}
	return d
}

// Unmarshal decodes the AST to a structure.
// This API is an EXPERIENTIAL one and might be removed in the future.
func (d *Decoder) Unmarshal(val *Value, dest interface{}) error {
	if !isPointer(dest) {
		return errors.New("dest must be a pointer")
	}

	return d.unmarshal(val, dest)
}

// todo: think about using reflect.Value instead of interface{} as the type of dest
func (d *Decoder) unmarshal(val *Value, dest interface{}) error {
	if val == nil {
		return errors.New("value is a nil pointer")
	}
	// todo: currently I only consider the valid conversion, need ensure behaviors if fail
	switch val.NodeType {
	case Number:
		return setNumber(val, dest)
	case String:
		return setString(val, dest)
	case Null:
		return setNull(dest)
	case Bool:
		return setBool(val, dest)
	case Array:
		return setArray(val, dest)
	case Object:
		return setObject(val, dest)
	}
	return errors.New("invalid value")
}

func setObject(val *Value, dest interface{}) error {
	obj := val.AstValue.(*ObjectAst)
	// *T --> T
	rval := reflect.ValueOf(dest).Elem()
	rtyp := reflect.TypeOf(dest).Elem()

	var err error
	for index := 0; index < rtyp.NumField(); index++ {

		fieldTyp := rtyp.Field(index)
		tag := fieldTyp.Tag.Get(jsonTAG)
		if tag == "" && !fieldTyp.Anonymous {
			continue
		}

		// handle embedded fields
		if fieldTyp.Anonymous {
			fieldVal := reflect.New(fieldTyp.Type)
			err = NewDecoder().unmarshal(val, fieldVal.Interface())
			if err != nil {
				return err
			}

			rval.Field(index).Set(fieldVal.Elem())
			continue
		}

		astVal := obj.KvMap[tag]

		// construct a pointer type of field type to store the data
		// we cannot pass the field directly because it's not a reference but a value, however,
		// we want to change the value itself. So pass it in side unmarshal as an interface doesn't work.
		fieldVal := reflect.New(rtyp.Field(index).Type)
		err = NewDecoder().unmarshal(&astVal, fieldVal.Interface())
		if err != nil {
			return err
		}

		rval.Field(index).Set(fieldVal.Elem())
	}
	return nil
}

// setArray sets the json array into golang a slice or an array.
func setArray(val *Value, dest interface{}) error {
	ars := val.AstValue.(*ArrayAst).Values

	kind := reflect.TypeOf(dest).Elem().Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		// todo: ignore error or report it?
		return nil
	}

	boundary := reflect.ValueOf(dest).Elem().Len()

	for i, value := range ars {
		// all available fields in an array are filled, we needn't to continue
		if i >= boundary && kind == reflect.Array {
			return nil
		}

		// double elem get from a pointer of array to the array element type
		// *[]T --Elem()--> []T --Elem()--> T
		elemType := reflect.TypeOf(dest).Elem().Elem()
		newVal := reflect.New(elemType)
		err := NewDecoder().unmarshal(&value, newVal.Interface())
		if err != nil {
			return err
		}
		elem := newVal.Elem()

		// this logic only applies to slice because array has a fixed length.
		if i >= boundary {
			// append the element into the array,
			// instead of the pointer to the array
			na := reflect.Append(reflect.ValueOf(dest).Elem(), elem)
			reflect.ValueOf(dest).Elem().Set(na)
			continue
		}

		reflect.ValueOf(dest).Elem().Index(i).Set(elem)
	}

	return nil
}

func setNull(dest interface{}) error {
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
		reflect.ValueOf(dest).Elem().SetInt(numberAst.GetInt64())
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		reflect.ValueOf(dest).Elem().SetUint(numberAst.GetUint64())
		return nil

	case reflect.Float32,
		reflect.Float64:
		reflect.ValueOf(dest).Elem().SetFloat(numberAst.GetFloat64())
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
