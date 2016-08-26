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

    // Reciever defines this function as a method to the given field
    Reciever    *FieldMap
    receiverMap *FieldListMap

    Params *FieldListMap

    Results *FieldListMap
}

// NewFunctionMap creates a new, empty function map
func NewFunctionMap() *FunctionMap {
    return &FunctionMap{
        Name:        "",
        Reciever:    nil,
        receiverMap: nil,
        Params:      nil,
        Results:     nil,
    }
}

// Visit ...
func (fm *FunctionMap) Visit(n ast.Node) ast.Visitor {

    switch node := n.(type) {

    case *ast.FuncDecl:
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

func (fm *FunctionMap) Finalize() {

    // pull out the reciever if there is one
    if nil != fm.receiverMap && fm.receiverMap.Count() > 0 {
        fm.Reciever = (*fm.receiverMap)[0]
    }

    // Ensure there are no nil lists
    if nil == fm.Params {
        fm.Params = NewFieldListMap()
    }
    if nil == fm.Results {
        fm.Results = NewFieldListMap()
    }
}