package astjson

import (
	"errors"
	"fmt"
)

var (
	ErrFieldNotExist = errors.New("field not exist")
)

// todo: think about is it possible to embed those features into AST parser directly
// By this way, during parsing, we could parse it directly instead of hindsight.

// Validator is used in Walker and validates a Value according to the logic you need.
// Validator will be executed before Manipulator.
type Validator func(value *Value) error

// Manipulator allows to change the Value and its children values.
// Manipulator will be executed after Validator.
type Manipulator func(value *Value)

type Walker struct {
	value              *Value
	compulsoryFields   []string
	validators         map[string]Validator
	optionalValidators map[string]Validator

	// literalValidator is registered by Validate when compulsoryFields is emtpy
	literalValidator Validator
}

func NewWalker(value *Value) *Walker {
	var w Walker
	w.value = value
	w.validators = make(map[string]Validator)
	w.optionalValidators = make(map[string]Validator)
	w.compulsoryFields = make([]string, 0, 5)
	return &w
}

// Field requires field is compulsory during walking on current layer.
func (w *Walker) Field(field string) *Walker {
	w.compulsoryFields = append(w.compulsoryFields, field)
	return w
}

// Validate validates a key's value last registered be Field,
// or a literal or an array on current layer.
// It will end up walking ast Value when the validator returns a non-nil error.
// It panics when no field is submitted before.
// It overrides the old one if multiple validators are submitted
func (w *Walker) Validate(validator Validator) *Walker {
	l := len(w.compulsoryFields)
	if l == 0 {
		w.literalValidator = validator
		return w
	}
	key := w.compulsoryFields[l-1]
	w.validators[key] = validator
	return w
}

// ValidateKey is an explicit method of Validate as it requires a key on current layer.
func (w *Walker) ValidateKey(key string, validator Validator) *Walker {
	w.Field(key).Validate(validator)
	return w
}

// Optional marks a key is optional alongside a validator on current layer.
// The validator is respect only if the key exit.
func (w *Walker) Optional(key string, validator Validator) *Walker {
	w.optionalValidators[key] = validator
	return w
}

// Walk executes all handlers submitted to it and return a final Value after walking.
// The error is returned when validator fails, and the value is nil.
func (w *Walker) Walk() (*Value, error) {
	nodeTyp := w.value.NodeType
	switch nodeTyp {
	case String, Number, Bool, Null, Array:
		return w.checkArrayAndLiteral()

	case Object:
		return w.checkObject()
	}
	return w.value, nil
}

func (w *Walker) checkArrayAndLiteral() (*Value, error) {
	if w.literalValidator != nil {
		if err := w.literalValidator(w.value); err != nil {
			return nil, err
		}
	}
	return w.value, nil
}

func (w *Walker) checkObject() (*Value, error) {
	objectMap := w.value.AstValue.(*ObjectAst).m
	for _, field := range w.compulsoryFields {
		if _, ok := objectMap[field]; !ok {
			return nil, fmt.Errorf("%w: %s", ErrFieldNotExist, field)
		}
	}
	for key, validator := range w.optionalValidators {
		val, ok := objectMap[key]
		// optional fields are allowed
		if !ok {
			continue
		}
		if err := validator(&val); err != nil {
			return nil, err
		}
	}
	for key, validator := range w.validators {
		// the key for sure exist because of the registering way of Validator
		val, _ := objectMap[key]
		if err := validator(&val); err != nil {
			return nil, err
		}
	}
	return w.value, nil
}

// Path allows you to jump inside the Value with the given path,
// while the original walker status is still preserved in the new returned one.
// It panics if the path doesn't exist.
// todo: think about how to preserve the original sub-date while adding new information.
// todo: could consider how gin framework and context design.
// todo: don't panic, return error with a functional way.
func (w *Walker) Path(paths ...string) *Walker {
	// todo: implement me in the future.
	return w
}
