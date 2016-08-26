package python

import "fmt"

// MapConverter is used to convert a map[string]interface{}
// to and from python dictionaries
type MapConverter struct{}

// GoTransitionType should be the char* to a json string
func (mc *MapConverter) GoTransitionType() string {
    return "*C.char"
}

// CTransitionType returns a char* to a json string
func (mc *MapConverter) CTransitionType() string {
    return "char*"
}

// ConvertGoFromC uses json to unmarshall the string
func (mc *MapConverter) ConvertGoFromC(varName string) string {
    return fmt.Sprintf("mapFromJSON(%s)", varName)
}

// ConvertGoToC uses json to create a string
func (mc *MapConverter) ConvertGoToC(varName string) string {
    return fmt.Sprintf("mapToJSON(%s)", varName)
}

// ConvertCFromGo no conversion necessary here
func (mc *MapConverter) ConvertCFromGo(varName string) string {
    return fmt.Sprintf("%s", varName)
}

// ConvertCToGo no conversion necessary here
func (mc *MapConverter) ConvertCToGo(varName string) string {
    return fmt.Sprintf("%s", varName)
}

// ConvertPyFromC parses a c string to dict with the utility function
func (mc *MapConverter) ConvertPyFromC(varName string) string {
    return fmt.Sprintf("ParseMapJSON(%s)", varName)
}

// ConvertPyToC uses the utility function to create a string from dict
func (mc *MapConverter) ConvertPyToC(varName string) string {
    return fmt.Sprintf("DictToJSON(%s)", varName)
}

// ValidatePyValue checks that the PyObject* is a dictionary
func (mc *MapConverter) ValidatePyValue(varName string) string {
    return fmt.Sprintf("PyDict_Check(%s)", varName)
}

// CDeclarations declares methods for json conversions to and from dict objs
func (mc *MapConverter) CDeclarations() string {
    return `
PyObject* ParseMapJSON(char *json);
char* DictToJSON(PyObject *dict);
`
}

// CDefinitions defines methods for json conversions to and from dict objs
func (mc *MapConverter) CDefinitions() string {
    return `
PyObject*
ParseMapJSON(char *json) 
{
    PyObject *jsonMod = PyImport_ImportModuleNoBlock("json");
    PyObject *argList = Py_BuildValue("(s)", json);
    PyObject *loadsFunc = PyObject_GetAttrString(jsonMod, "loads");
    PyObject *result = PyEval_CallObject(loadsFunc, argList);
    Py_DECREF(loadsFunc);
    Py_DECREF(argList);
    free(json);
    return result;
}

char*
DictToJSON(PyObject *dict) 
{
    if (!PyDict_Check(dict)) {
        PyErr_BadArgument();
        return NULL;
    }
    PyObject *jsonMod = PyImport_ImportModuleNoBlock("json");
    PyObject *argList = Py_BuildValue("(O)", dict);
    PyObject *dumpsFunc = PyObject_GetAttrString(jsonMod, "dumps");
    PyObject *result = PyEval_CallObject(dumpsFunc, argList);
    Py_DECREF(dumpsFunc);
    Py_DECREF(argList);
    const char* str = PyString_AsString(result);
    char *mine = (char*)malloc(strlen(str)+1);
    Py_DECREF(result);
    return strcpy(mine, str);
}
`
}

// GoDefinitions defines methods for json conversions to and from map objs
func (mc *MapConverter) GoDefinitions() string {
    return `
func mapFromJSON(jsonStr *C.char) map[string]interface{} {
    defer C.free(unsafe.Pointer(jsonStr))
    m := make(map[string]interface{})
    err := json.Unmarshal([]byte(C.GoString(jsonStr)), m)
    if err != nil {
        str := C.CString(err.Error())
        defer C.free(unsafe.Pointer(str))
        C.PyErr_SetString(C.PyExc_RuntimeError, str)
    }
    return m
}

func mapToJSON(m map[string]interface{}) *C.char {
    bytes, err := json.Marshal(m)
    if err != nil {
        str := C.CString(err.Error())
        defer C.free(unsafe.Pointer(str))
        C.PyErr_SetString(C.PyExc_RuntimeError, str)
        return C.CString("{}")
    }
    return C.CString(string(bytes))
}
`
}
