package python

import "fmt"

// MapConverter is used to convert a map[string]interface{}
// to and from python dictionaries
type MapConverter struct{}

// GoType returns internal type's name
func (mc *MapConverter) GoType() string {
    return "map[string]interface{}"
}

// CMemberType is simply a PyObject* but should be a dictionary
func (mc *MapConverter) CMemberType() string {
    return "PyObject*"
}

// GoIncomingArgType should be the char* to a json string
func (mc *MapConverter) GoIncomingArgType() string {
    return "*C.char"
}

// COutgoingArgType returns a char* to a json string
func (mc *MapConverter) COutgoingArgType() string {
    return "const char*"
}

// ConvertFromCValue uses json to unmarshall the string
func (mc *MapConverter) ConvertFromCValue(varName string) string {
    return fmt.Sprintf("mapFromJSON(%s)", varName)
}

// ConvertToCValue uses json to create a string
func (mc *MapConverter) ConvertToCValue(varName string) string {
    return fmt.Sprintf("mapToJSON(%s)", varName)
}

// ConvertFromGoValue uses the python json module to convert to a PyDictObject*
func (mc *MapConverter) ConvertFromGoValue(varName string) string {
    return fmt.Sprintf("ParseMapJSON(%s)", varName)
}

// ConvertToGoValue uses the python json module to convert to a string
func (mc *MapConverter) ConvertToGoValue(varName string) string {
    return fmt.Sprintf("DictToJSON(%s)", varName)
}

// PyMemberDefTypeEnum returns T_OBJECT_EX because these should
// be instantiable python objects
func (mc *MapConverter) PyMemberDefTypeEnum() string {
    return "T_OBJECT_EX"
}

// PyTupleTarget is just a PyObject*
func (mc *MapConverter) PyTupleTarget(ident int) string {
    return fmt.Sprintf("PyObject *dict%d", ident)
}

// PyParseTupleArgs returns multiple args because we can leverage
// python to type-check this one for us
func (mc *MapConverter) PyParseTupleArgs(ident int) string {
    return fmt.Sprintf("&PyDict_Type, &dict%d", ident)
}

// PyTupleResult returns the name of the checked var generated for PyTupleTarget
func (mc *MapConverter) PyTupleResult(ident int) string {
    return fmt.Sprintf("dict%d", ident)
}

// PyTupleFormat returns the type for type-asserted object
func (mc *MapConverter) PyTupleFormat() string {
    return "O!"
}

// CDeclarations declares methods for json conversions to and from dict objs
func (mc *MapConverter) CDeclarations() string {
    return `
PyObject* ParseMapJSON(char *json);
const char* DictToJSON(PyObject *dict);
`
}

// CDefinitions defines methods for json conversions to and from dict objs
func (mc *MapConverter) CDefinitions() string {
    return `
PyObject*
ParseMapJSON(char *json) 
{
    PyObject *jsonMod = PyImport_ImportModuleNoBlock("json");
    PyObject *argList = Py_BuildValue("s", json);
    PyObject *loadsFunc = PyObject_GetAttrString(jsonMod, "loads");
    PyObject *result = PyEval_CallObject(loadsFunc, argList);
    Py_DECREF(loadsFunc);
    Py_DECREF(argList);
    return result;
}

const char*
DictToJSON(PyObject *dict) 
{
    if (!PyDict_Check(dict)) {
        PyErr_BadArgument();
        return NULL;
    }
    PyObject *jsonMod = PyImport_ImportModuleNoBlock("json");
    PyObject *argList = Py_BuildValue("O", dict);
    PyObject *dumpsFunc = PyObject_GetAttrString(jsonMod, "dumps");
    PyObject *result = PyEval_CallObject(dumpsFunc, argList);
    Py_DECREF(dumpsFunc);
    Py_DECREF(argList);
    const char* str = PyString_AsString(result);
    Py_DECREF(result);
    return str;
}
`
}

// GoDefinitions defines methods for json conversions to and from map objs
func (mc *MapConverter) GoDefinitions() string {
    return `
func mapFromJSON(jsonStr *C.char) map[string]interface{} {
    m := make(map[string]interface{})
    err := json.Unmarshal([]byte(C.GoString(jsonStr)), m)
    if err != nil {
        str := C.CString(err.Error())
        defer C.free(unsafe.Pointer(str))
        C.PyErr_SetString(C.PyExc_RuntimeError, str)
    }
    return m
}

func mapToJSON(m map[string]interface{}) string {
    bytes, err := json.Marshal(m)
    if err != nil {
        str := C.CString(err.Error())
        defer C.free(unsafe.Pointer(str))
        C.PyErr_SetString(C.PyExc_RuntimeError, str)
        return "{}"
    }
    return string(bytes)
}
`
}
