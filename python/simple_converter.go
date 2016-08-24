package python

import "fmt"

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

// CDeclarations returns nothing
func (sc *SimpleConverter) CDeclarations() string {
    return ""
}

// CDefinitions returns nothing
func (sc *SimpleConverter) CDefinitions() string {
    return ""
}

// GoDefinitions returns nothing
func (sc *SimpleConverter) GoDefinitions() string {
    return ""
}
