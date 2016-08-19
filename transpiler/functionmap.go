package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
)

// FunctionMap ...
type FunctionMap struct {

    // Name is the name of this function
    Name string

    // RecieverType is the name of the caller type
    RecieverType string

    Params *FieldListMap

    Results *FieldListMap
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
        list := &FieldListMap{}
        if nil == fm.Params {
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

/*
type FieldListMap []*FieldMap

func (flm *FieldListMap) Visit(n ast.Node) ast.Visitor {

    switch node := n.(type) {

    case *ast.FieldList:
        return flm

    case *ast.Field:
        field := &FieldMap{}
        if len(node.Names) == 0 {
            field.Unamed = true
        } else {
            field.Name = node.Names[0]
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

type FieldMap struct {
    Name    string
    Type    string
    Unamed  bool
    Pointer bool
}

func (fm *FieldMap) Visit(n ast.Node) ast.Visitor {

    switch n.(type) {

    case *ast.StarExpr:
        fm.Pointer = true
        fm.Type += "*"
        return fm

    case *ast.Ident:
        if "" == fm.Name && !fm.Unnamed {
            fm.Name = n.(*ast.Ident).String()
        } else {
            fm.Type += n.(*ast.Ident).String()
        }
        return nil

    case *ast.SelectorExpr:
        exp := n.(*ast.SelectorExpr)
        fm.Type += expressionToString(exp)
        return nil

    case nil:
        return nil

    default:
        fmt.Printf("FieldMap unhandled %s\n", reflect.TypeOf(n))
        return nil

    }

}*/
