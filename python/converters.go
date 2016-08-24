package python

func getConverter(goType string) converter {
    if nil == converters[goType] {
        converters[goType] = NewInternalConverter(goType)
    }
    return converters[goType]
}

var converters = map[string]converter{
    "string": &SimpleConverter{
        Name:         "string",
        CType:        "char*",
        GoCType:      "*C.char",
        CGoType:      "char*",
        FromC:        "C.GoString(%s)",
        ToC:          "C.CString(%s)", // BUG(rydrman): this is a memory leak if ownership is not passed to c
        FromGo:       "%s",
        ToGo:         "%s",
        PyMemberType: "T_CHAR",
        PyTupleFmt:   "s",
    },
    "int": &SimpleConverter{
        Name:         "int",
        CType:        "int",
        GoCType:      "int",
        CGoType:      "int",
        FromC:        "%s",
        ToC:          "%s",
        FromGo:       "%s",
        ToGo:         "%s",
        PyMemberType: "T_INT",
        PyTupleFmt:   "i",
    },
    "float": &SimpleConverter{
        Name:         "float",
        CType:        "float",
        GoCType:      "float",
        CGoType:      "float",
        FromC:        "%s",
        ToC:          "%s",
        FromGo:       "%s",
        ToGo:         "%s",
        PyMemberType: "T_FLOAT",
        PyTupleFmt:   "f",
    },
    "bool": &SimpleConverter{
        Name:         "bool",
        CType:        "char",
        GoCType:      "C.char",
        CGoType:      "char",
        FromC:        "int(%s) > 0",
        ToC:          "C.char(%s)",
        FromGo:       "%s",
        ToGo:         "%s",
        PyMemberType: "T_BOOL",
        PyTupleFmt:   "b",
    },
    "time.Time": &SimpleConverter{
        Name:         "Time",
        CType:        "unsigned long",
        GoCType:      "int64",
        CGoType:      "long long",
        FromC:        "time.Unix(%s)",
        ToC:          "%s.Unix()",
        FromGo:       "%s",
        ToGo:         "%s",
        PyMemberType: "T_INT",
        PyTupleFmt:   "i",
    },
    "map[string]interface{}": &DictConverter{},
}
