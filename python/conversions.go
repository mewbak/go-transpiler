package python

import "fmt"

// converter represents a single go-type and provides the
// necessary functionality to convert variables of this
// type between go and c code
type converter interface {

    // GoType should return the name of the go type
    // that this converter is representing
    GoType() string

    // CMemberType should return the c type name
    // for this converters go type. This is the type of variable
    // that should be stored in a c struct to define the underlying
    // data of a python object
    CMemberType() string

    // GoIncomingArgType returns the argument type that should be used
    // on go functions exported to c (usually C.* types)
    GoIncomingArgType() string

    // COutgoingArgType returns the argument type that c should use to
    // define extern functions exported from go
    COutgoingArgType() string

    // ConvertFromCValue should return valid go code that
    // takes a c representation of this type and converts it
    // back to go code. The result of this function should be assignable
    ConvertFromCValue(varName string) string

    // ConvertToCValue should return valid go code that
    // takes a go representation of this type and converts it
    // to c code. The result of this function should be assignable
    ConvertToCValue(varName string) string

    // ConvertFromGoValue should return valid C code that
    // takes a go representation of this type and converts it
    // back to c code. The result of this function should be assignable
    //
    // Because Go is considered the owner of all data, and manages
    // conversions through caches, this should
    // usually be a straight assignment or simple lookup, not
    // creating new data
    ConvertFromGoValue(varName string) string

    // ConvertToGoValue should return valid c code that
    // takes a c representation of this type and converts it
    // to the type expected by go code. The result of this function
    // should be assignable
    ConvertToGoValue(varName string) string

    // PyMemberDefTypeEnum returns the python enum that properly
    // describes this type in a member definition (eg PY_OBJECT_EX)
    PyMemberDefTypeEnum() string

    // PyTupleTarget returns a string that defines target variables
    // to be used when parsing this varaible type out of a tuple set.
    // The given identifier string should be appended to variable names
    // in order to ensure uniqueness
    PyTupleTarget(ident int) string

    // PyParseTupleArgs returns a string that are arguments
    // to be passed to tuple argument parding functions in python api.
    // The given identifier string should be appended to variable names
    // in order to match what would be returned from PythonTupleTarget
    PyParseTupleArgs(ident int) string

    // PyTupleResult returns the name of the var generated for PyTupleTarget
    PyTupleResult(ident int) string

    // PyTupleFormat returns the set of format character(s) that
    // define this variable type as represented in python tuples
    // ex (int vars would return "i")
    PyTupleFormat() string
}

// SimpleConverter is a simple set of
// strings to define a type conversion
type SimpleConverter struct {
    Name         string
    CType        string
    GoCType      string
    CGoType      string
    FromC        string
    ToC          string
    FromGo       string
    ToGo         string
    GoConversion string
    PyMemberType string
    PyTupleFmt   string
}

// GoType returns internal type's name
func (sc *SimpleConverter) GoType() string {
    return sc.Name
}

// CMemberType returns the CType string
func (sc *SimpleConverter) CMemberType() string {
    return sc.CType
}

// GoIncomingArgType returns the GoCType string
func (sc *SimpleConverter) GoIncomingArgType() string {
    return sc.GoCType
}

// COutgoingArgType returns the CGoType string
func (sc *SimpleConverter) COutgoingArgType() string {
    return sc.CGoType
}

// ConvertFromCValue formats the FromC string
func (sc *SimpleConverter) ConvertFromCValue(varName string) string {
    return fmt.Sprintf(sc.FromC, varName)
}

// ConvertToCValue formats the ToC string
func (sc *SimpleConverter) ConvertToCValue(varName string) string {
    return fmt.Sprintf(sc.ToC, varName)
}

// ConvertFromGoValue formats the FromGo string
func (sc *SimpleConverter) ConvertFromGoValue(varName string) string {
    return fmt.Sprintf(sc.FromGo, varName)
}

// ConvertToGoValue formats the ToGo string
func (sc *SimpleConverter) ConvertToGoValue(varName string) string {
    return fmt.Sprintf(sc.ToGo, varName)
}

// PyMemberDefTypeEnum returns the PyMemberType string
func (sc *SimpleConverter) PyMemberDefTypeEnum() string {
    return sc.PyMemberType
}

// PyTupleTarget uses the CType to create the tuple target
func (sc *SimpleConverter) PyTupleTarget(ident int) string {
    return fmt.Sprintf("%s %s%d", sc.CType, sc.Name, ident)
}

// PyParseTupleArgs just returns the same var name as PyTupleTarget
func (sc *SimpleConverter) PyParseTupleArgs(ident int) string {
    return fmt.Sprintf("&%s%d", sc.Name, ident)
}

// PyTupleResult returns the name of the var generated for PyTupleTarget
func (sc *SimpleConverter) PyTupleResult(ident int) string {
    return fmt.Sprintf("%s%d", sc.Name, ident)
}

// PyTupleFormat returns the type defined in this struct
func (sc *SimpleConverter) PyTupleFormat() string {
    return sc.PyTupleFmt
}

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
