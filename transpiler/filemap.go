package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
)

// FileMap ...
type FileMap struct {
    Name      string
    Types     []*TypeMap
    Functions []*FunctionMap
}

// NewFileMap ...
func NewFileMap(file *ast.File, name string) *FileMap {
    fm := &FileMap{
        name,
        make([]*TypeMap, 0),
        make([]*FunctionMap, 0),
    }

    ast.Walk(fm, file)
    return fm
}

// Visit ...
func (fm *FileMap) Visit(n ast.Node) ast.Visitor {

    switch n.(type) {
    case *ast.File:
        return fm
    case *ast.DeclStmt:
        return fm
    case *ast.GenDecl:
        tm := &TypeMap{}
        fm.Types = append(fm.Types, tm)
        return tm
    case *ast.FuncDecl:
        funcMap := &FunctionMap{}
        return funcMap
    case nil:
        return nil
    default:
        fmt.Printf("FileMap unhandled: %s %s\n", reflect.TypeOf(n), n)
        return nil
    }

}
