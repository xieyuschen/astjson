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

// Manipulator allows to change the Value and its children KvMap.
// Manipulator will be executed after Validator.
type Manipulator func(value *Value)

type Walker struct {
	// head marks the head of the walker
	head *Walker
	next *Walker

	value *Value
	err   error

	// field stores the current field for possible validator
	field              string
	compulsoryFields   []string
	validators         map[string]Validator
	optionalValidators map[string]Validator

	// literalValidator is registered by Validate when compulsoryFields is emtpy
	literalValidator Validator
}

func NewWalker(value *Value) *Walker {
	var w Walker
	w.value = value
	w.head = &w
	w.validators = make(map[string]Validator)
	w.optionalValidators = make(map[string]Validator)
	w.compulsoryFields = make([]string, 0, 5)
	return &w
}

// Field requires field is compulsory during walking on current layer.
func (w *Walker) Field(field string) *Walker {
	w.compulsoryFields = append(w.compulsoryFields, field)
	w.field = field
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
	w.validators[w.field] = validator
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

// Walk executes all handlers submitted to it and the sub-walker created by Path and EndPath.
// The returned value always be the value in original walker. See Clone to walk from current walker
// as the very beginning.
// The error is returned when validator fails, and the value is nil.
func (w *Walker) Walk() (*Value, error) {
	current := w.head
	for current != nil {
		if _, err := current.walk(); err != nil {
			return nil, err
		}
		current = current.next
	}
	return w.head.value, nil
}

// Clone return a new path which removes the relationships with its ancestor
func (w *Walker) Clone() *Walker {
	// todo: implement me!
	return nil
}

// walk executes all handlers submitted to it and return a final Value after walking.
// The error is returned when validator fails, and the value is nil.
func (w *Walker) walk() (*Value, error) {
	if w.err != nil {
		return nil, w.err
	}
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
	objectMap := w.value.AstValue.(*ObjectAst).KvMap
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

// Path jumps the walker inside the given path of ast.
// If the given path doesn't exist, error is raised when Walk.
// The scope of Path ends when EndPath is called.
func (w *Walker) Path(path string) *Walker {
	switch w.value.NodeType {
	case Object:
		obj := w.value.AstValue.(*ObjectAst)
		val, ok := obj.KvMap[path]
		if ok {
			n := Walker{
				head:  w.head,
				value: &val,
			}
			w.next = &n
			return &n
		}
		fallthrough
	case Number, Null, Array, String, Bool:
		w.err = fmt.Errorf("path %s doesn't exist in nodeype %s", path, w.value.NodeType)
		return w
	}
	return w
}

// EndPath returns to the root walker after Path enters a path
func (w *Walker) EndPath() *Walker {
	if w.err != nil {
		w.head.err = w.err
	}
	return w.head
}
