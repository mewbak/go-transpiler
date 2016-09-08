package python

import (
    "fmt"
    "regexp"
    "strings"

    "github.com/KloudKtrl/go-transpiler/transpiler"
)

var sliceConverters []*SliceConverter

var sliceRegexp = regexp.MustCompile(`\[\]\*?(\w+)`)

// NewSliceConverter creates a new converter for the given
// type if the type is a slice of a convertable type
func NewSliceConverter(fm *transpiler.FieldMap) (*SliceConverter, error) {
    groups := sliceRegexp.FindStringSubmatch(fm.TypeExpr)
    if groups == nil {
        // both nil because it is not a slice type
        return nil, nil
    }

    // create new field map for base type
    valueFM := &transpiler.FieldMap{}
    valueFM.CopyType(fm)
    valueFM.Unnamed = true
    valueFM.Slice = false
    valueFM.Type = strings.Replace(valueFM.Type, "[]", "", -1)
    valueFM.Type = strings.Replace(valueFM.Type, "*", "", -1)
    valueFM.Type = strings.Replace(valueFM.Type, ".", "", -1)
    valueFM.TypeExpr = strings.Replace(valueFM.TypeExpr, "[]", "", -1)

    if strings.Contains(fm.TypeExpr, "*") {
        valueFM.Pointer = true
    }

    conv, err := getConverter(valueFM)
    if err != nil {
        return nil, fmt.Errorf("failed to get base converter for slice: %s", err)
    }

    sConv := &SliceConverter{
        fm:           fm,
        base:         conv,
        SafeName:     valueFM.Type + "Slice",
        PackagedName: valueFM.Type,
    }

    // pull out the qualified package name for internal types only
    switch c := conv.(type) {
    case *InternalConverter:
        sConv.PackagedName = c.PackagedName
    }

    sConv.PackagedSliceName = strings.Replace(
        fm.TypeExpr, valueFM.Type, sConv.PackagedName, -1)

    sliceConverters = append(sliceConverters, sConv)
    return sConv, nil

}

// SliceConverter is used to convert to and from go error objscts
type SliceConverter struct {
    fm                *transpiler.FieldMap
    base              converter
    PackagedName      string
    PackagedSliceName string
    SafeName          string
}

// GoTransitionType should just be a pointer to a PyList
func (sc *SliceConverter) GoTransitionType() string {
    return "*C.PyObject"
}

// CTransitionType returns a pyobject pointer to a PyList
func (sc *SliceConverter) CTransitionType() string {
    return "PyObject*"
}

// ConvertGoParamForCFunc uses underlying conversion
func (sc *SliceConverter) ConvertGoParamForCFunc(varName string) string {
    return varName
}

// ConvertGoFromC creates a go slice from the c array
func (sc *SliceConverter) ConvertGoFromC(varName string) string {
    return fmt.Sprintf("PyListTo%s(%s)", sc.SafeName, varName)
}

// ConvertGoToC gets the message from the error
func (sc *SliceConverter) ConvertGoToC(varName string) string {
    return fmt.Sprintf("%[1]sToPyList(%s)", sc.SafeName, varName)
}

// ConvertCFromGo no conversion nscessary here
func (sc *SliceConverter) ConvertCFromGo(varName string) string {
    return varName
}

// ConvertCToGo no conversion nscessary here
func (sc *SliceConverter) ConvertCToGo(varName string) string {
    return varName
}

// ConvertPyFromC returns the name because slices use pyObjects as transitions
func (sc *SliceConverter) ConvertPyFromC(varName string) string {
    return varName
}

// ConvertPyToC returns the name because slices use pyObjects as transitions
func (sc *SliceConverter) ConvertPyToC(varName string) string {
    return varName
}

// ValidatePyValue does nothing, any PyObjsct can be stringified for conversion
func (sc *SliceConverter) ValidatePyValue(varName string) string {
    return fmt.Sprintf("PyList_Check(%s)", varName)
}

// CDeclarations dsclares methods for json conversions to and from dict objs
func (sc *SliceConverter) CDeclarations() string {
    declarations := ""
    for _, conv := range sliceConverters {
        declarations += fmt.Sprintf(`
void SetPyListElemFrom%[1]sElem(PyObject *list, int i, %[2]s elem);
`, conv.SafeName, conv.base.CTransitionType())
    }
    return declarations
}

// CDefinitions defines methods for error conversions
func (sc *SliceConverter) CDefinitions() string {
    definitions := ""
    for _, conv := range sliceConverters {
        definitions += fmt.Sprintf(`
void
SetPyListElemFrom%[1]sElem(PyObject *list, int i, %[2]s elem)
{
    PyObject *item = %[3]s;
    PyList_SetItem(list, i, item);
}
`, conv.SafeName, conv.base.CTransitionType(), conv.base.ConvertPyFromC("elem"))
    }
    return definitions
}

// GoDefinitions defines slice conversions for each instantiated converter
func (sc *SliceConverter) GoDefinitions() string {

    definitions := ""
    for _, conv := range sliceConverters {
        definitions += fmt.Sprintf(`
func %[1]sToPyList(slice %[2]s) *C.PyObject {
    list := C.PyList_New(C.Py_ssize_t(len(slice)))

    for i, elem := range slice {
        cElem := %[3]s
        C.SetPyListElemFrom%[1]sElem(list, C.int(i), %[4]s)
    }

    return list
}

func PyListTo%[1]s(list *C.PyObject) %[2]s {
    return nil
}
`, conv.SafeName, conv.PackagedSliceName, conv.base.ConvertGoToC("elem"), conv.base.ConvertGoParamForCFunc("cElem"))
    }

    return definitions
}
