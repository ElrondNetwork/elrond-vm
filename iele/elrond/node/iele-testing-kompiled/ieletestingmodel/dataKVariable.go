// File provided by the K Framework Go backend. Timestamp: 2019-07-04 01:26:11.488

package ieletestingmodel

// KVariable is a KObject representing a KVariable item in K
type KVariable struct {
	Name string
}

func (*KVariable) referenceType() kreferenceType {
	return kvariableRef
}


// NewKVariable creates a new object and returns the reference.
func (ms *ModelState) NewKVariable(name string) KReference {
	return ms.addObject(&KVariable{Name: name})
}