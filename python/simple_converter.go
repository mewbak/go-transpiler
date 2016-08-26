package python

import "fmt"

// SimpleConverter is a simple set of
// strings to define a type conversion
type SimpleConverter struct {
    Name       string
    GoTType    string
    CTType     string
    GoFromC    string
    GoToC      string
    CFromGo    string
    CToGo      string
    PyToC      string
    PyFromC    string
    PyValidate string
    PyTupleFmt string
}

// GoType returns internal type's name
func (sc *SimpleConverter) GoType() string {
    return sc.Name
}

// GoTransitionType returns the GoTType string
func (sc *SimpleConverter) GoTransitionType() string {
    return sc.GoTType
}

// CTransitionType returns the CTType string
func (sc *SimpleConverter) CTransitionType() string {
    return sc.CTType
}

// ConvertGoFromC formats the FromC string
func (sc *SimpleConverter) ConvertGoFromC(varName string) string {
    return fmt.Sprintf(sc.GoFromC, varName)
}

// ConvertGoToC formats the ToC string
func (sc *SimpleConverter) ConvertGoToC(varName string) string {
    return fmt.Sprintf(sc.GoToC, varName)
}

// ConvertCFromGo formats the FromGo string
func (sc *SimpleConverter) ConvertCFromGo(varName string) string {
    return fmt.Sprintf(sc.CFromGo, varName)
}

// ConvertCToGo formats the ToGo string
func (sc *SimpleConverter) ConvertCToGo(varName string) string {
    return fmt.Sprintf(sc.CToGo, varName)
}

// ConvertPyFromC formats the PyFromC with varName
func (sc *SimpleConverter) ConvertPyFromC(varName string) string {
    return fmt.Sprintf(sc.PyFromC, varName)
}

// ConvertPyToC formats PyTpC with the varName
func (sc *SimpleConverter) ConvertPyToC(varName string) string {
    return fmt.Sprintf(sc.PyToC, varName)
}

// ValidatePyValue formats the PyValidate string
func (sc *SimpleConverter) ValidatePyValue(varName string) string {
    return fmt.Sprintf(sc.PyValidate, varName)
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
    return `
func btoi (b bool) int {
    if b {
        return 1
    }
    return 0
}
`
}
