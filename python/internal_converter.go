package python

import "fmt"

// InternalConverter is used to convert a type that is
// defined within the package that is being transpiled
type InternalConverter struct {
    Name string
}

// NewInternalConverter creates a new internal converter
// object for the internal type defined in the given TypeMap
func NewInternalConverter(name string) *InternalConverter {
    return &InternalConverter{name}
}

// GoType returns internal type's name
func (ic *InternalConverter) GoType() string {
    return ic.Name
}

// GoTransitionType returns int64 because it is the key for
// this item in the cache
func (ic *InternalConverter) GoTransitionType() string {
    return "int64"
}

// CTransitionType returns long long because it is the key for
// this item in the cache
func (ic *InternalConverter) CTransitionType() string {
    return "long long"
}

// ConvertGoFromC uses cache lookup to get object
func (ic *InternalConverter) ConvertGoFromC(varName string) string {
    return fmt.Sprintf("cache%s[%s]", ic.Name, varName)
}

// ConvertGoToC uses cache to set object
func (ic *InternalConverter) ConvertGoToC(varName string) string {
    return fmt.Sprintf("getCached%s(%s)", ic.Name, varName)
}

// ConvertCFromGo is a simple assigment of the long long from go
func (ic *InternalConverter) ConvertCFromGo(varName string) string {
    return fmt.Sprintf("%s", varName)
}

// ConvertCToGo just passes the vanilla long long value
func (ic *InternalConverter) ConvertCToGo(varName string) string {
    return fmt.Sprintf("%s", varName)
}

// ConvertPyFromC ...TODO... his should create a valid item
func (ic *InternalConverter) ConvertPyFromC(varName string) string {
    return fmt.Sprintf("NULL")
}

// ConvertPyToC accesses the cache key for this object
func (ic *InternalConverter) ConvertPyToC(varName string) string {
    return fmt.Sprintf("((%s*)%s)->go%s", ic.Name, varName, ic.Name)
}

// ValidatePyValue checks that the given var is of the right type
func (ic *InternalConverter) ValidatePyValue(varName string) string {
    return fmt.Sprintf("PyObject_TypeCheck(%s, &%s_type)", varName, ic.Name)
}

// PyTupleTarget is just a PyObject*
func (ic *InternalConverter) PyTupleTarget(ident int) string {
    return fmt.Sprintf("PyObject *%s%d", ic.Name, ident)
}

// PyParseTupleArgs returns multiple args because we can leverage
// python to type-check this one for us
func (ic *InternalConverter) PyParseTupleArgs(ident int) string {
    return fmt.Sprintf("&%s_type, &%s%d", ic.Name, ic.Name, ident)
}

// PyTupleResult returns the name of the checked var generated for PyTupleTarget
func (ic *InternalConverter) PyTupleResult(ident int) string {
    return fmt.Sprintf("%s%d", ic.Name, ident)
}

// PyTupleFormat returns the type for type-asserted object
func (ic *InternalConverter) PyTupleFormat() string {
    return "O!"
}

// CDeclarations returns nothing
func (ic *InternalConverter) CDeclarations() string {
    return ""
}

// CDefinitions returns nothing
func (ic *InternalConverter) CDefinitions() string {
    return ""
}

// GoDefinitions returns nothing
func (ic *InternalConverter) GoDefinitions() string {
    return ""
}
