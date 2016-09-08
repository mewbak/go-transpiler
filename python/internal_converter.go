package python

import (
    "fmt"
    "regexp"

    "github.com/KloudKtrl/go-transpiler/transpiler"
)

var intTypeRegex = regexp.MustCompile(`\*?\w+`)

// InternalConverter is used to convert a type that is
// defined within the package that is being transpiled
type InternalConverter struct {
    Name         string
    PackagedName string
    fm           *transpiler.FieldMap
}

// NewInternalConverter creates a new internal converter
// object for the internal type defined in the given TypeMap
func NewInternalConverter(fm *transpiler.FieldMap) (*InternalConverter, error) {

    if fm.TypeExpr != intTypeRegex.FindString(fm.TypeExpr) {
        return nil, fmt.Errorf(
            "cannot create converter, not valid internal type %s", fm.TypeExpr)
    }

    return &InternalConverter{
        Name:         fm.Type,
        PackagedName: fmt.Sprintf("%s.%s", fm.Package.Name, fm.Type),
        fm:           fm,
    }, nil
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

// ConvertGoParamForCFunc casts a go int64 to C.lonlong
func (ic *InternalConverter) ConvertGoParamForCFunc(varName string) string {
    return fmt.Sprintf("C.longlong(%s)", varName)
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
    return fmt.Sprintf("Py_None; Py_INCREF(Py_None)")
}

// ConvertPyToC accesses the cache key for this object
func (ic *InternalConverter) ConvertPyToC(varName string) string {
    return fmt.Sprintf("((%s*)%s)->go%s", ic.Name, varName, ic.Name)
}

// ValidatePyValue checks that the given var is of the right type
func (ic *InternalConverter) ValidatePyValue(varName string) string {
    return fmt.Sprintf("PyObject_TypeCheck(%s, &%s_type)", varName, ic.Name)
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
