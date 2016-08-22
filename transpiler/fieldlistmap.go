package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
)

// FieldListMap ...
type FieldListMap []*FieldMap

// Visit ...
func (flm *FieldListMap) Visit(n ast.Node) ast.Visitor {

    switch node := n.(type) {

    case *ast.FieldList:
        return flm

    case *ast.Field:
        field := &FieldMap{}
        if len(node.Names) == 0 {
            field.Unnamed = true
        }
        (*flm) = append(*flm, field)
        return field

    case nil:
        return nil

    default:
        fmt.Printf("FieldListMap unhandled %s\n", reflect.TypeOf(n))
        return nil

    }

}

// FieldMap ...
type FieldMap struct {
    Name    string
    Type    string
    Unnamed bool
    Pointer bool
}

// Visit ...
func (fm *FieldMap) Visit(n ast.Node) ast.Visitor {

    switch node := n.(type) {

    case *ast.StarExpr:
        fm.Pointer = true
        fm.Type += "*"
        return fm

    case *ast.Ident:
        if "" == fm.Name && !fm.Unnamed {
            fm.Name = node.String()
        } else {
            fm.Type += node.String()
        }
        return nil

    case *ast.BasicLit:
        // often string tags: ie `json:"name"`
        return nil

    case *ast.SelectorExpr:
        fm.Type += ExpressionToString(node)
        return nil

    case nil:
        return nil

    default:
        fmt.Printf("FieldMap unhandled %s\n", reflect.TypeOf(n))
        return nil

    }

}
