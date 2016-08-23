package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
)

// FileMap ...
type FileMap struct {
    Name        string
    Package     string
    Types       []*TypeMap
    TypesByName map[string]*TypeMap
    Functions   []*FunctionMap
}

// NewFileMap ...
func NewFileMap(file *ast.File, name string) *FileMap {
    fm := &FileMap{
        Name:        name,
        Package:     "",
        Types:       make([]*TypeMap, 0),
        TypesByName: make(map[string]*TypeMap),
        Functions:   make([]*FunctionMap, 0),
    }

    ast.Walk(fm, file)
    fm.Finalize()
    return fm
}

// Visit ...
func (fm *FileMap) Visit(n ast.Node) ast.Visitor {

    switch node := n.(type) {
    case *ast.File:
        return fm
    case *ast.DeclStmt:
        return fm
    case *ast.GenDecl:
        tm := NewTypeMap()
        tm.Package = fm.Package
        fm.Types = append(fm.Types, tm)
        return tm
    case *ast.FuncDecl:
        funcMap := NewFunctionMap()
        fm.Functions = append(fm.Functions, funcMap)
        return funcMap
    case *ast.Ident:
        if fm.Package == "" {
            fm.Package = node.String()
            return nil
        }
        return nil
    case nil:
        return nil
    default:
        fmt.Printf("FileMap unhandled: %s %s\n", reflect.TypeOf(n), n)
        return nil
    }

}

// Finalize ...
func (fm *FileMap) Finalize() {
    for _, tm := range fm.Types {
        fm.TypesByName[tm.Name] = tm
        tm.Finalize()
    }
    for _, f := range fm.Functions {
        f.Finalize()
        if f.Reciever == nil {
            continue
        }
        t := f.Reciever.TypeName
        if fm.TypesByName[t] != nil {
            fm.TypesByName[t].Functions = append(
                fm.TypesByName[t].Functions, f)
        }
    }
}
