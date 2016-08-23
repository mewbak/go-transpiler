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

    // Named members is populated with only members
    // that have proper field names in this type
    NamedMembers *FieldListMap

    // Package is the name of the package that this
    // type belongs too (set by the calling FileMap)
    Package string

    // Functions is a list of functions that this type
    // is a reciever for (populated by the calling FileMap)
    Functions []*FunctionMap

    // IsInterface is set to true when this type
    // defines an interface
    IsInterface bool

    // IsStruct is set to true when this type
    // defines a struct type
    IsStruct bool
}

// NewTypeMap Create a new TypeMap with default members
func NewTypeMap() *TypeMap {
    return &TypeMap{
        Name:         "",
        BaseType:     "",
        Members:      NewFieldListMap(),
        NamedMembers: NewFieldListMap(),
        Package:      "",
        Functions:    make([]*FunctionMap, 0),
        IsInterface:  false,
        IsStruct:     false,
    }
}

// Visit ...
func (tm *TypeMap) Visit(n ast.Node) ast.Visitor {

    switch node := n.(type) {

    case *ast.TypeSpec:
        return tm

    case *ast.Ident:
        if tm.Name == "" {
            tm.Name = node.String()
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
        tm.Members = NewFieldListMap()
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

// Finalize ...
func (tm *TypeMap) Finalize() {

    for _, m := range *tm.Members {
        if m.Name != "" {
            tm.NamedMembers.Add(m)
        }
    }
}
