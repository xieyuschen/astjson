package astjson

// ParseE is an EXPERIENTIAL function which might be removed in the long run development
func (p *Parser) ParseE() ValueE {
	return ValueE{
		Value: p.Parse(),
		e:     nil,
	}
}

// ValueE is an EXPERIENTIAL structure which might be removed in the long run development
// we won't implement error interface for ValueE because it is error-prone
type ValueE struct {
	*Value
	e error
}

func (v ValueE) Decompose() (*Value, error) {
	return v.Value, v.e
}

func (v ValueE) Error() error {
	return v.e
}
