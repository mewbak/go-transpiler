package python

import (
    "fmt"
    "regexp"
    "strings"

    "github.com/KloudKtrl/go-transpiler/transpiler"
)

func mustGetConverter(fieldMap *transpiler.FieldMap) converter {
    conv, err := getConverter(fieldMap)
    if err != nil {
        panic(err.Error())
    }
    return conv
}

func getConverter(fieldMap *transpiler.FieldMap) (converter, error) {
    if fieldMap.TypeExpr == "" {
        fmt.Println(fieldMap)
        return nil, fmt.Errorf("empty gotype: %s", fieldMap.Name)
    }
    conv := matchConverter(fieldMap.TypeExpr)
    if nil != conv {
        return conv, nil
    }

    conv = createConverter(fieldMap)
    if nil != conv {
        re := fmt.Sprintf("^%s$", regexp.QuoteMeta(fieldMap.TypeExpr))
        converters[re] = conv
        return conv, nil
    }

    return nil, fmt.Errorf("cannot find or create converter for go type %s", fieldMap)
}

func matchConverter(expr string) converter {
    for key, conv := range converters {
        if m, _ := regexp.MatchString(key, expr); m {
            return conv
        }
    }
    return nil
}

func createConverter(fieldMap *transpiler.FieldMap) converter {

    // no way to create converter for type from external package
    if strings.Contains(fieldMap.TypeExpr, ".") {
        return nil
    }

    sConv, err := NewSliceConverter(fieldMap)
    if nil != sConv {
        return sConv
    }
    if nil != err {
        fmt.Printf("error creating slice converter: %s\n", err)
        return nil
    }

    iConv, err := NewInternalConverter(fieldMap)
    if nil != iConv {
        return iConv
    }

    fmt.Printf("error creating internal converter: %s\n", err)

    return nil
}

var converters = map[string]converter{
    `^string$`: &SimpleConverter{
        Name:    "string",
        GoTType: "*C.char",
        CTType:  "char*",
        CParam:  "%s",
        GoFromC: "C.GoString(%s)",
        GoToC:   "C.CString(%s)",
        //this is freeing the string produced by ToC above
        CFromGo:    "PyString_FromString(%s); free(%[1]s)",
        CToGo:      "%s",
        PyToC:      "PyString_AsString(%s)",
        PyFromC:    "PyString_FromString(%s)",
        PyValidate: "PyString_Check(%s)",
    },
    `^int$`: &SimpleConverter{
        Name:       "int",
        GoTType:    "int",
        CTType:     "long",
        CParam:     "C.long(%s)",
        GoFromC:    "%s",
        GoToC:      "%s",
        CFromGo:    "PyInt_FromLong(%s)",
        CToGo:      "%s",
        PyToC:      "PyInt_AsLong(%s)",
        PyFromC:    "PyInt_FromLong(%s)",
        PyValidate: "PyInt_Check(%s)",
    },
    `^float$`: &SimpleConverter{
        Name:       "float",
        GoTType:    "C.double",
        CTType:     "double",
        CParam:     "%s",
        GoFromC:    "%s",
        GoToC:      "%s",
        CFromGo:    "PyFloat_FromDouble(%s)",
        CToGo:      "%s",
        PyToC:      "PyFloat_ASDouble(%s)",
        PyFromC:    "PyFloat_FromDouble(%s)",
        PyValidate: "PyFloat_Check(%s)",
    },
    `^byte$`: &SimpleConverter{
        Name:       "byte",
        GoTType:    "byte",
        CTType:     "char",
        CParam:     "C.char(%s)",
        GoFromC:    "%s",
        GoToC:      "%s",
        CFromGo:    "%s",
        CToGo:      "%s",
        PyToC:      "char(PyInt_AsLong(%s))",
        PyFromC:    "PyInt_FromLong((long)%s)",
        PyValidate: "PyInt_Check(%s)",
    },
    `^bool$`: &SimpleConverter{
        Name:       "bool",
        GoTType:    "int",
        CTType:     "int",
        CParam:     "C.int(%s)",
        GoFromC:    "%s != 0",
        GoToC:      "btoi(%s)",
        CFromGo:    "%s",
        CToGo:      "%s",
        PyToC:      "%s == Py_True",
        PyFromC:    "%s == 0 ? Py_False : Py_True",
        PyValidate: "PyBool_Check(%s)",
    },
    `^error$`:                     &ErrorConverter{},
    `time\.Time`:                  &DateConverter{},
    `map\[string\]interface\s*{}`: &MapConverter{},
    `^interface\s*{}`:             &MapConverter{},
}
