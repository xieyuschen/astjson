// Code generated by "stringer -type=numberType"; DO NOT EDIT.

package astjson

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[floatNumber-0]
	_ = x[unsignedInteger-1]
	_ = x[integer-2]
}

const _numberType_name = "floatNumberunsignedIntegerinteger"

var _numberType_index = [...]uint8{0, 11, 26, 33}

func (i numberType) String() string {
	if i >= numberType(len(_numberType_index)-1) {
		return "numberType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _numberType_name[_numberType_index[i]:_numberType_index[i+1]]
}
