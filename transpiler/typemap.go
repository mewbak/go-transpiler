package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
)

// TypeMap ...
type TypeMap struct {

    // Name is the name of this type
    Name string

    // BaseName is the name of the inherited type
    // if this is neither a struct nor interface
    BaseType string

    // Members is the list of members if this is
    // a struct or interface type
    Members *FieldListMap

    IsInterface bool

    IsStruct bool
}

// Visit ...
func (tm *TypeMap) Visit(n ast.Node) ast.Visitor {

    switch node := n.(type) {

    case *ast.TypeSpec:
        return tm

    case *ast.Ident:
        if tm.Name == "" {
            tm.Name = node.String()
            fmt.Printf("type %s\n", tm.Name)
        }
        return nil

    case *ast.InterfaceType:
        tm.IsInterface = true
        tm.IsStruct = false
        return tm

    case *ast.StructType:
        tm.IsInterface = false
        tm.IsStruct = true
        return tm

    case *ast.FieldList:
        tm.Members = &FieldListMap{}
        return tm.Members

    case *ast.SelectorExpr:
        tm.BaseType = ExpressionToString(node)

    case nil:
        return nil

    default:
        fmt.Printf("TypeMap unhandled: %s\n", reflect.TypeOf(n))
        return nil

    }

    return tm

}
