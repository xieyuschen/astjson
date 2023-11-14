package astjson

import (
	"errors"
	"fmt"
)

func ShouldEqualTrue() Validator {
	return func(value *Value) error {
		if !IsBool(value) {
			return fmt.Errorf("value should be an bool type: %s", value)
		}
		if !GetBool(value) {
			return errors.New("value should be true")
		}
		return nil
	}
}

func ShouldEqualFalse() Validator {
	return func(value *Value) error {
		if !IsBool(value) {
			return fmt.Errorf("value should be an bool type: %s", value)
		}
		if GetBool(value) {
			return errors.New("value should be false")
		}
		return nil
	}
}

func ShouldEqualString(str string) Validator {
	return func(value *Value) error {
		if !IsString(value) {
			return fmt.Errorf("value should be a stirng type: %s", value)
		}
		actual := GetString(value)
		if actual != str {
			return fmt.Errorf("value is %s, not equal with expected value is %s", actual, str)
		}
		return nil
	}
}

func ShouldNotEqualString(str string) Validator {
	return func(value *Value) error {
		if !IsString(value) {
			return fmt.Errorf("value should be a stirng type: %s", value)
		}
		actual := GetString(value)
		if actual == str {
			return fmt.Errorf("value is %s, equal with expected value", str)

		}
		return nil
	}
}

func ShouldEqualNull() Validator {
	return func(value *Value) error {
		if !IsNull(value) {
			return fmt.Errorf("value should be a null type: %s", value)
		}

		return nil
	}
}
