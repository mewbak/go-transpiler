package python

import "fmt"

// DateConverter defines go->c conversions for date objects
type DateConverter struct{}

// GoType returns internal type's name
func (dc *DateConverter) GoType() string {
    return "Time"
}

// GoTransitionType returns int64
func (dc *DateConverter) GoTransitionType() string {
    return "int64"
}

// CTransitionType returns long long
func (dc *DateConverter) CTransitionType() string {
    return "long"
}

// ConvertGoFromC formats the FromC string
func (dc *DateConverter) ConvertGoFromC(varName string) string {
    return fmt.Sprintf("time.Unix(%s, 0)", varName)
}

// ConvertGoToC formats the ToC string
func (dc *DateConverter) ConvertGoToC(varName string) string {
    return fmt.Sprintf("%s.Unix()", varName)
}

// ConvertCFromGo formats the FromGo string
func (dc *DateConverter) ConvertCFromGo(varName string) string {
    return fmt.Sprintf("%s", varName)
}

// ConvertCToGo formats the ToGo string
func (dc *DateConverter) ConvertCToGo(varName string) string {
    return fmt.Sprintf("%s", varName)
}

// ConvertPyFromC formats the PyFromC with varName
func (dc *DateConverter) ConvertPyFromC(varName string) string {
    return fmt.Sprintf("PyDate_FromUnix(%s)", varName)
}

// ConvertPyToC formats PyTpC with the varName
func (dc *DateConverter) ConvertPyToC(varName string) string {
    return fmt.Sprintf("PyDate_AsUnix(%s)", varName)
}

// ValidatePyValue formats the PyValidate string
func (dc *DateConverter) ValidatePyValue(varName string) string {
    return fmt.Sprintf("PyDate_Import_Check(%s)", varName)
}

// CDeclarations returns nothing
func (dc *DateConverter) CDeclarations() string {
    return `
PyObject* PyDate_FromUnix(long unx);
long PyDate_AsUnix(PyObject* date);
int PyDate_Import_Check(PyObject *obj);
`
}

// CDefinitions returns nothing
func (dc *DateConverter) CDefinitions() string {
    return `
PyObject*
PyDate_FromUnix(long unx)
{
    PyDateTime_IMPORT;
    PyObject *argList = Py_BuildValue("(i)", unx);
    PyObject *result = PyDateTime_FromTimestamp(argList);
    Py_DECREF(argList);
    return result;
}

long
PyDate_AsUnix(PyObject* date) 
{
    PyDateTime_IMPORT;
    if (!PyDate_Check(date)) {
        PyErr_SetString(PyExc_TypeError, "expected valid datetime.datetime object");
        return 0;
    }
    PyObject *strftimeFunc = PyObject_GetAttrString(date, "strftime");
    PyObject *argList = Py_BuildValue("(s)", "%s");
    PyObject *result = PyEval_CallObject(strftimeFunc, argList);
    Py_DECREF(argList);
    Py_DECREF(strftimeFunc);
    PyObject* intObj = PyInt_FromString(PyString_AsString(result), NULL, 10);
    Py_DECREF(result);
    long res = PyInt_AsLong(intObj);
    Py_DECREF(intObj);
    return res;
}

int
PyDate_Import_Check(PyObject *obj)
{
    PyDateTime_IMPORT;
    return PyDate_Check(obj);
}

`
}

// GoDefinitions returns nothing
func (dc *DateConverter) GoDefinitions() string {
    return ""
}
