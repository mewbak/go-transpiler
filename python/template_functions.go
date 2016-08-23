package python

import (
    "reflect"
    "regexp"
    "strings"
    "text/template"

    "github.com/KloudKtrl/go-transpiler/transpiler"
)

var templateFuncs template.FuncMap

func init() {
    //template functions
    templateFuncs = template.FuncMap{
        "externalTypes": externalTypes,
        "camelToSnake":  camelToSnake,

        "cMemberType":        cMemberType,
        "goCFuncArgType":     goCFuncArgType,
        "convertFromCValue":  convertFromCValue,
        "convertToCValue":    convertToCValue,
        "convertFromGoValue": convertFromGoValue,
        "convertToGoValue":   convertToGoValue,
        "pyMemberDefType":    pyMemberDefType,
        "pyTupleTarget":      pyTupleTarget,
        "pyParseTupleArgs":   pyParseTupleArgs,

        "pyTupleFormat": pyTupleFormat,
        "notLast":       notLast,
    }
}

func externalTypes(fm *transpiler.FileMap, pm *transpiler.PackageMap) []*transpiler.TypeMap {
    var external []*transpiler.TypeMap
    for _, extFile := range pm.Files {
        for _, tm := range extFile.Types {
            isInt := false
            for _, internal := range fm.Types {
                if tm == internal {
                    isInt = true
                    break
                }
            }
            if !isInt {
                external = append(external, tm)
            }
        }
    }
    return external
}

func camelToSnake(name string) string {

    // this regex guarentees that ever GROUP is a camelCase section,
    // but each match may contain more than one of such sections. This
    // is necessary in order to catch ALLCAPS sections such as URL in baseURLString
    exp := regexp.MustCompile(
        `^(\w[a-z]+)|([A-Z][a-z]+)|([A-Z]+)([A-Z][a-z]+)|([A-Z]+)$`)
    matches := exp.FindAllStringSubmatch(name, -1)

    var parts []string
    for _, m := range matches {
        for i, g := range m {
            if i == 0 {
                continue
            }
            if g != "" {
                parts = append(parts, strings.ToLower(g))
            }
        }
    }

    return strings.Join(parts, "_")

}

func getConverter(goType string) converter {
    if nil == converters[goType] {
        converters[goType] = NewInternalConverter(goType)
    }
    return converters[goType]
}

func cMemberType(goType string) string {
    return getConverter(goType).CMemberType()
}

func goCFuncArgType(goType string) string {
    return getConverter(goType).GoCFuncArgType()
}

func convertFromCValue(goType, varName string) string {
    return getConverter(goType).ConvertFromCValue(varName)
}

func convertToCValue(goType, varName string) string {
    return getConverter(goType).ConvertToCValue(varName)
}

func convertFromGoValue(goType, varName string) string {
    return getConverter(goType).ConvertFromGoValue(varName)
}

func convertToGoValue(goType, varName string) string {
    return getConverter(goType).ConvertToGoValue(varName)
}

func pyMemberDefType(goType string) string {
    return getConverter(goType).PyMemberDefTypeEnum()
}

func pyTupleTarget(goType string, ident int) string {
    return getConverter(goType).PyTupleTarget(ident)
}

func pyParseTupleArgs(goType string, ident int) string {
    return getConverter(goType).PyParseTupleArgs(ident)
}

func pyArgTuplFormat(goType string) string {
    return getConverter(goType).PyTupleFormat()
}

var converters = map[string]converter{
    "string": &SimpleConverter{
        Name:         "string",
        CType:        "char*",
        GoCType:      "*C.Char",
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
        FromC:        "int(%s) > 0",
        ToC:          "C.char(%s)",
        FromGo:       "%s",
        ToGo:         "%s",
        PyMemberType: "T_BOOL",
        PyTupleFmt:   "f",
    },
    "time.Time": &SimpleConverter{
        Name:         "time.Time",
        CType:        "unsigned long",
        GoCType:      "int64",
        FromC:        "time.Unix(%s)",
        ToC:          "%s.Unix()",
        FromGo:       "%s",
        ToGo:         "%s",
        PyMemberType: "T_INT",
        PyTupleFmt:   "f",
    },
}

func pyTupleFormat(args []*transpiler.FieldMap) string {

    res := "|"
    for _, a := range args {

        res += pyArgTuplFormat(a.Name)

    }
    return res

}

type counter interface {
    Count() int
}

func notLast(i int, slice interface{}) bool {

    switch s := slice.(type) {
    case counter:
        return i < s.Count()-1
    }

    switch reflect.TypeOf(slice).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(slice)
        return i < s.Len()-1
    case reflect.Ptr:
        return notLast(i, reflect.ValueOf(slice).Elem())
    }

    return false
}
