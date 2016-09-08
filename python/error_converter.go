package python

import "fmt"

// ErrorConverter is used to convert to and from go error objects
type ErrorConverter struct{}

// GoTransitionType should be the char* to the error message
func (ec *ErrorConverter) GoTransitionType() string {
    return "*C.char"
}

// CTransitionType returns a char* to the error message
func (ec *ErrorConverter) CTransitionType() string {
    return "char*"
}

// ConvertGoParamForCFunc does nothing
func (ec *ErrorConverter) ConvertGoParamForCFunc(varName string) string {
    return varName
}

// ConvertGoFromC creates a new go error object
func (ec *ErrorConverter) ConvertGoFromC(varName string) string {
    return fmt.Sprintf("CStringToGoError(%s)", varName)
}

// ConvertGoToC gets the message from the error
func (ec *ErrorConverter) ConvertGoToC(varName string) string {
    return fmt.Sprintf("C.CString(%s.Error())", varName)
}

// ConvertCFromGo no conversion necessary here
func (ec *ErrorConverter) ConvertCFromGo(varName string) string {
    return fmt.Sprintf("%s", varName)
}

// ConvertCToGo no conversion necessary here
func (ec *ErrorConverter) ConvertCToGo(varName string) string {
    return fmt.Sprintf("%s", varName)
}

// ConvertPyFromC converts a string into a python exception instance
func (ec *ErrorConverter) ConvertPyFromC(varName string) string {
    return fmt.Sprintf("CreatePyException(%s)", varName)
}

// ConvertPyToC uses the utility function to create a string from exception
func (ec *ErrorConverter) ConvertPyToC(varName string) string {
    return fmt.Sprintf("PyExceptionToString(%s)", varName)
}

// ValidatePyValue does nothing, any PyObject can be stringified for conversion
func (ec *ErrorConverter) ValidatePyValue(varName string) string {
    return fmt.Sprintf("1")
}

// CDeclarations declares methods for json conversions to and from dict objs
func (ec *ErrorConverter) CDeclarations() string {
    return `
PyObject* CreatePyException(char *message);
char* PyExceptionToString(PyObject *pyExc);
`
}

// CDefinitions defines methods for error conversions
func (ec *ErrorConverter) CDefinitions() string {
    return `
PyObject*
CreatePyException(char *message) 
{
    PyObject *argList = Py_BuildValue("(s)", message);
    PyObject *result = PyEval_CallObject(PyExc_Exception, argList);
    Py_DECREF(argList);
    free(message);
    return result;
}

char*
PyExceptionToString(PyObject *pyExc) 
{
    const char* str = PyString_AsString(pyExc);
    char *mine = (char*)malloc(strlen(str)+1);
    Py_DECREF(pyExc);
    return strcpy(mine, str);
}
`
}

// GoDefinitions defines nothing
func (ec *ErrorConverter) GoDefinitions() string {
    return `
func CStringToGoError(str *C.char) error {
    defer C.free(unsafe.Pointer(str))
    return errors.New(C.GoString(str))
}
`
}
