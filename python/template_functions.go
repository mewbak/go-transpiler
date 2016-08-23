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
        "cType":            cType,
        "goCType":          goCType,
        "camelToSnake":     camelToSnake,
        "pythonMemberType": pythonMemberType,
        "pyArgFormat":      pyArgFormat,
        "packagedType":     packagedType,
        "isInternalType":   isInternalType,
        "notLast":          notLast,
        "externalTypes":    externalTypes,
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

func cType(goType string) string {

    switch goType {

    case "string":
        return "char*"

    case "float":
        fallthrough
    case "int":
        return goType

    case "bool":
        return "char"

    case "time.Time":
        return "unsigned long"

    case "map[string]interface{}":
        return "PyDictObject*"

    default:
        return "PyObject*"
    }

}

func goCType(goType string) string {

    switch goType {

    case "string":
        return "*C.char"

    case "float":
        fallthrough
    case "int":
        return "C." + goType

    case "bool":
        return "C.char"

    case "unsigned long":
        "time.Time"

    default:
        return "*C.PyObject"
    }

}

func pythonMemberType(goType string) string {

    switch goType {

    case "string":
        return "T_CHAR"

    case "float":
        fallthrough
    case "bool":
        fallthrough
    case "int":
        return "T_" + strings.ToUpper(goType)

    default:
        return "T_OBJECT_EX"
    }

}

func pyArgFormat(args []*transpiler.FieldMap) string {

    res := "|"
    for _, a := range args {

        if a.Name == "" {
            continue
        }

        switch a.Type {

        case "string":
            res += "s"

        case "bool":
            res += "b"

        case "int":
            res += "i"

        default:
            res += "O!" // will check that the incoming pyObj is of the right type
        }
    }
    return res

}

func packagedType(p string, t string) string {

    // if it is already a selector, no need for a package
    if 1 < len(strings.Split(t, ".")) {
        return t
    }

    if isInternalType(t) {
        return t
    }

    return p + "." + t

}

var internalTypes = []string{
    "uint8",
    "uint16",
    "uint32",
    "uint64",

    "int",
    "int8",
    "int16",
    "int32",
    "int64",

    "float",
    "float32",
    "float64",

    "complex64",
    "complex128",

    "byte",
    "rune",

    "string",
    "bool",
}

func isInternalType(t string) bool {
    for _, it := range internalTypes {
        if t == it || "*"+it == t {
            return true
        }
    }
    return false
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
