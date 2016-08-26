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
        "externalTypes":            externalTypes,
        "camelToSnake":             camelToSnake,
        "pyModuleName":             func() string { return pyModuleName },
        "filterSupportedFunctions": filterSupportedFunctions,
        "canTranspileField":        canTranspileField,

        "goTransitionType": goTransitionType,
        "cTransitionType":  cTransitionType,
        "convertGoFromC":   convertGoFromC,
        "convertGoToC":     convertGoToC,
        "convertCFromGo":   convertCFromGo,
        "convertCToGo":     convertCToGo,
        "convertPyFromC":   convertPyFromC,
        "convertPyToC":     convertPyToC,
        "validatePyValue":  validatePyValue,

        "notLast": notLast,
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

func filterSupportedFunctions(funcs []*transpiler.FunctionMap) []*transpiler.FunctionMap {
    var ok []*transpiler.FunctionMap

    for _, f := range funcs {
        accept := true
        for _, p := range *f.Params {
            if !canTranspileField(p) {
                accept = false
                break
            }
        }
        for _, r := range *f.Results {
            if !canTranspileField(r) {
                accept = false
                break
            }
        }
        if accept {
            ok = append(ok, f)
        }
    }
    return ok
}

func canTranspileField(fm *transpiler.FieldMap) bool {
    if nil != converters[fm.Type] {
        return true
    }

    // no support for types from other packages
    if strings.Contains(fm.TypeExpr, ".") {
        return false
    }

    //no support for other anonymous types
    if strings.Contains(fm.TypeExpr, "{") {
        return false
    }

    //no support for other slice or array types
    if strings.Contains(fm.TypeExpr, "[") {
        return false
    }

    return true
}

func goTransitionType(goType string) string {
    return getConverter(goType).GoTransitionType()
}

func cTransitionType(goType string) string {
    return getConverter(goType).CTransitionType()
}

func convertGoFromC(goType, varName string) string {
    return getConverter(goType).ConvertGoFromC(varName)
}

func convertGoToC(goType, varName string) string {
    return getConverter(goType).ConvertGoToC(varName)
}

func convertCFromGo(goType, varName string) string {
    return getConverter(goType).ConvertCFromGo(varName)
}

func convertCToGo(goType, varName string) string {
    return getConverter(goType).ConvertCToGo(varName)
}

func convertPyFromC(goType, varName string) string {
    return getConverter(goType).ConvertPyFromC(varName)
}

func convertPyToC(goType, varName string) string {
    return getConverter(goType).ConvertPyToC(varName)
}

func validatePyValue(goType, varName string) string {
    return getConverter(goType).ValidatePyValue(varName)
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
