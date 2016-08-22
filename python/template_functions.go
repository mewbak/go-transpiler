package python

import (
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
        "camelToSnake":     camelToSnake,
        "pythonMemberType": pythonMemberType,
        "pyArgFormat":      pyArgFormat,
    }
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

    default:
        return "PyObject*"
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
