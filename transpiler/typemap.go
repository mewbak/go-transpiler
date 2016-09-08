package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
)

// TypeMap maps a type definition from a go abstract
// syntax tree. The typemap is intended to visit ast.TypeSpec
// nodes and their children
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

    // Package is the package that this type
    // belongs too (set by the calling FileMap)
    Package *PackageMap

    // Functions is a list of functions that this type
    // is a receiver for (populated by the calling FileMap)
    Functions []*FunctionMap

    // Interface is set to true when this type
    // defines an interface
    Interface bool

    // Struct is set to true when this type
    // defines a struct type
    Struct bool
}

// NewTypeMap Create a new TypeMap with default members
func NewTypeMap() *TypeMap {
    return &TypeMap{
        Name:         "",
        BaseType:     "",
        Members:      NewFieldListMap(),
        NamedMembers: NewFieldListMap(),
        Package:      nil,
        Functions:    make([]*FunctionMap, 0),
        Interface:    false,
        Struct:       false,
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
        tm.Interface = true
        tm.Struct = false
        return tm

    case *ast.StructType:
        tm.Interface = false
        tm.Struct = true
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

// SetPackage sets the package for this type for
// easier access in transpiling functions
func (tm *TypeMap) SetPackage(pm *PackageMap) {
    tm.Package = pm
    for _, m := range *tm.Members {
        m.SetPackage(pm)
    }
}

// Finalize ...
func (tm *TypeMap) Finalize() {

    tm.Members.Finalize()
    for _, m := range *tm.Members {
        if m.Name != "" {
            tm.NamedMembers.Add(m)
        }
    }
}
