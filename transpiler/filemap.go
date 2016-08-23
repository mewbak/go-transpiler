package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
)

// FileMap maps an entire go file's abstract syntax tree.
type FileMap struct {
    Name        string
    Package     *PackageMap
    Types       []*TypeMap
    TypesByName map[string]*TypeMap
    Functions   []*FunctionMap
}

// NewFileMap creates a new FileMap instance with the file mapping
// of the given ast.File object and file name
func NewFileMap(file *ast.File, name string) *FileMap {

    fm := &FileMap{
        Name:        name,
        Package:     nil,
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

    switch n.(type) {

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
        return nil

    case nil:
        return nil

    default:
        fmt.Printf("FileMap unhandled: %s %s\n", reflect.TypeOf(n), n)
        return nil

    }

}

// SetPackage sets the package for this file and all
// underlying type definitions for easier access in
// transpiling functions
func (fm *FileMap) SetPackage(pm *PackageMap) {
    fm.Package = pm
    for _, tm := range fm.Types {
        tm.SetPackage(pm)
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
