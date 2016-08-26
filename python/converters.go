package python

func getConverter(goType string) converter {
    if "" == goType {
        panic("empty gotype")
    }
    if nil == converters[goType] {
        converters[goType] = NewInternalConverter(goType)
    }
    return converters[goType]
}

var converters = map[string]converter{
    "string": &SimpleConverter{
        Name:    "string",
        GoTType: "*C.char",
        CTType:  "char*",
        GoFromC: "C.GoString(%s)",
        GoToC:   "C.CString(%s)",
        //this is freeing the string produced by ToC above
        CFromGo:    "PyString_FromString(%s); free(%[1]s)",
        CToGo:      "%s",
        PyToC:      "PyString_AsString(%s)",
        PyFromC:    "PyString_FromString(%s)",
        PyValidate: "PyString_Check(%s)",
    },
    "int": &SimpleConverter{
        Name:       "int",
        GoTType:    "int",
        CTType:     "long",
        GoFromC:    "%s",
        GoToC:      "%s",
        CFromGo:    "PyInt_FromLong(%s)",
        CToGo:      "%s",
        PyToC:      "PyInt_AsLong(%s)",
        PyFromC:    "PyInt_FromLong(%s)",
        PyValidate: "PyInt_Check(%s)",
    },
    "float": &SimpleConverter{
        Name:       "float",
        GoTType:    "C.double",
        CTType:     "double",
        GoFromC:    "%s",
        GoToC:      "%s",
        CFromGo:    "PyFloat_FromDouble(%s)",
        CToGo:      "%s",
        PyToC:      "PyFloat_ASDouble(%s)",
        PyFromC:    "PyFloat_FromDouble(%s)",
        PyValidate: "PyFloat_Check(%s)",
    },
    "bool": &SimpleConverter{
        Name:       "bool",
        GoTType:    "int",
        CTType:     "int",
        GoFromC:    "%s != 0",
        GoToC:      "btoi(%s)",
        CFromGo:    "%s",
        CToGo:      "%s",
        PyToC:      "%s == Py_True",
        PyFromC:    "%s == 0 ? Py_False : Py_True",
        PyValidate: "PyBool_Check(%s)",
    },
    "time.Time":              &DateConverter{},
    "map[string]interface{}": &MapConverter{},
}
