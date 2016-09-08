package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
)

// FunctionMap maps a function definition from a go abstract
// syntax tree. The function map is intended to visit ast.FuncDecl
// nodes and their children
type FunctionMap struct {

    // Name is the name of this function
    Name string

    // Receiver defines this function as a method to the given field
    Receiver    *FieldMap
    receiverMap *FieldListMap

    Params *FieldListMap

    Results *FieldListMap

    // Package is the package that this function
    // belongs too (set by the calling FileMap)
    Package *PackageMap
}

// NewFunctionMap creates a new, empty function map
func NewFunctionMap() *FunctionMap {
    return &FunctionMap{
        Name:        "",
        Receiver:    nil,
        receiverMap: nil,
        Params:      nil,
        Results:     nil,
    }
}

// Visit ...
func (fm *FunctionMap) Visit(n ast.Node) ast.Visitor {

    switch node := n.(type) {

    case *ast.FuncDecl:
        if node.Recv == nil {
            fm.receiverMap = NewFieldListMap()
        }
        return fm

    case *ast.Ident:
        if fm.Name == "" {
            fm.Name = node.String()
        }
        return nil

    case *ast.FuncType:
        return fm

    case *ast.FieldList:
        list := NewFieldListMap()
        if nil == fm.receiverMap {
            fm.receiverMap = list
        } else if nil == fm.Params {
            fm.Params = list
        } else {
            fm.Results = list
        }
        return list

    case *ast.BlockStmt:
        return nil

    case nil:
        return nil

    default:
        fmt.Printf("FunctionMap unhandled: %s\n", reflect.TypeOf(node))
        return nil

    }

    return fm

}

// SetPackage sets the package for this type for
// easier access in transpiling functions
func (fm *FunctionMap) SetPackage(pm *PackageMap) {
    fm.Package = pm
    if nil != fm.Receiver {
        fm.Receiver.SetPackage(pm)
    }
    fm.Params.SetPackage(pm)
    fm.Results.SetPackage(pm)
}

// Finalize goes over this function map and makes sure everything is set right
func (fm *FunctionMap) Finalize() {

    // pull out the receiver if there is one
    if nil != fm.receiverMap && fm.receiverMap.Count() > 0 {
        fm.Receiver = (*fm.receiverMap)[0]
    }

    // Ensure there are no nil lists
    if nil == fm.Params {
        fm.Params = NewFieldListMap()
    }
    if nil == fm.Results {
        fm.Results = NewFieldListMap()
    }

    if nil != fm.Receiver {
        fm.Receiver.Finalize()
    }
    fm.Params.Finalize()
    fm.Results.Finalize()

    // backfill param types (name1, name2 string)
    var last *FieldMap
    for i := fm.Params.Count() - 1; i >= 0; i-- {
        field := (*fm.Params)[i]
        if nil == last || field.Type != "" {
            last = field
        } else {
            field.CopyType(last)
        }
    }

}
