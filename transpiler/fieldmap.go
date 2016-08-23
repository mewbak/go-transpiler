package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
)

// FieldMap maps a field definition from a go abstract
// syntax tree. The FieldMap is intended to visit ast.Field
// nodes and their children
type FieldMap struct {
    Name     string
    Type     string
    TypeName string
    TypeExpr string
    Unnamed  bool
    Pointer  bool
}

// NewFieldMap creates a new, empty FieldMap
func NewFieldMap() *FieldMap {
    return &FieldMap{
        Name:     "",
        Type:     "",
        TypeName: "",
        TypeExpr: "",
        Unnamed:  false,
        Pointer:  false,
    }
}

// Visit ...
func (fm *FieldMap) Visit(n ast.Node) ast.Visitor {

    switch node := n.(type) {

    case *ast.StarExpr:
        fm.Pointer = true
        fm.TypeExpr += "*"
        return fm

    case *ast.Ident:
        if "" == fm.Name && !fm.Unnamed {
            fm.Name = node.String()
        } else {
            fm.Type += node.String()
            fm.TypeExpr += node.String()
            fm.TypeName = node.String()
        }
        return nil

    case *ast.BasicLit:
        // often string tags: ie `json:"name"`
        return nil

    case *ast.SelectorExpr:
        expr := ExpressionToString(node)
        fm.Type += expr
        fm.TypeExpr += expr
        return nil

    case nil:
        return nil

    default:
        fmt.Printf("FieldMap unhandled %s\n", reflect.TypeOf(n))
        return nil

    }

}
