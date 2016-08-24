// InternalConverter is used to convert a type that is
import "fmt"

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

// CMemberType returns PyObject* becuase it should be properly instantiated
func (ic *InternalConverter) CMemberType() string {
    return "PyObject*"
}

// GoIncomingArgType returns int64 because it is the key for
// this item in the cache
func (ic *InternalConverter) GoIncomingArgType() string {
    return "int64"
}

// COutgoingArgType returns long long because it is the key for
// this item in the cache
func (ic *InternalConverter) COutgoingArgType() string {
    return "long long"
}

// ConvertFromCValue uses cache lookup to get object
func (ic *InternalConverter) ConvertFromCValue(varName string) string {
    return fmt.Sprintf("cache%s[%s]", ic.Name, varName)
}

// ConvertToCValue uses cache to set object
func (ic *InternalConverter) ConvertToCValue(varName string) string {
    return fmt.Sprintf("getCached%s(%s)", ic.Name, varName)
}

// ConvertFromGoValue is a simple assigment of the long long from go
func (ic *InternalConverter) ConvertFromGoValue(varName string) string {
    return fmt.Sprintf("%s", varName)
}

// ConvertToGoValue accesses the cache key from the pyobject
func (ic *InternalConverter) ConvertToGoValue(varName string) string {
    return fmt.Sprintf("((%s*)%s)->go%s", ic.Name, varName, ic.Name)
}

// PyMemberDefTypeEnum returns T_OBJECT_EX because these should
// be instantiable python objects
func (ic *InternalConverter) PyMemberDefTypeEnum() string {
    return "T_OBJECT_EX"
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