// main transpiler package takes the current ktrlgo
// vendor dependancy and builds a python interface via
// the python c-api and cgo
package main

import (
    "flag"
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "os"
    "path"
    "reflect"
    "strings"
)

const fileHeader = `
// auto-generated file, do not edit

#define Py_LIMITED_API
#include <Python.h>
#include "structmember.h"
`

var fileCount = 0
var pkgCount = 0

func main() {

    here := path.Clean(".")

    var pkg string
    flag.StringVar(&pkg, "p", "", "used to transpile one or more packages")

    var file string
    flag.StringVar(&file, "f", "", "used to transpile one or more files")

    var out string
    flag.StringVar(&out, "o", here, "location to output to")
    flag.Parse()

    out = path.Clean(out)
    err := os.MkdirAll(out, os.ModePerm)
    if err != nil {
        panic(
            fmt.Sprintf("could not create output directory: %s", out),
        )
    }

    packages := strings.Split(pkg, " ")
    for _, pkgPath := range packages {
        if pkgPath == "" {
            continue
        }
        err := transpilePackage(path.Clean(pkgPath), out)
        if nil != err {
            fmt.Println(err)
        }
    }

    files := strings.Split(file, " ")
    for _, filename := range files {
        if filename == "" {
            continue
        }
        err := transpileFile(path.Clean(filename), out)
        if nil != err {
            fmt.Printf("> %s\n", err)
        }
    }

    fmt.Printf("\nDone:\n")
    fmt.Printf(" > %d packages transpiled\n", pkgCount)
    fmt.Printf(" > %d files transpiled\n", fileCount)
}

func transpilePackage(packageDir, outDir string) error {

    if _, err := os.Stat(packageDir); os.IsNotExist(err) {
        return fmt.Errorf("package directory does not exist %s", packageDir)
    }

    fileSet := token.NewFileSet()

    packages, err := parser.ParseDir(fileSet, packageDir, nil, 0)
    if nil != err {
        return fmt.Errorf("failed parsing dir: %s\n", err)
    }
    if len(packages) == 0 {
        return fmt.Errorf("no package(s) found at: %s", packageDir)
    }

    for pkgName, pkg := range packages {
        fmt.Printf("transpiling package: %s\n", pkgName)
        for name := range pkg.Files {

            err := transpileFile(name, outDir)
            if nil != err {
                fmt.Printf("> %s\n", err)
            }
        }
        pkgCount++
    }
    return nil

}

func transpileFile(filepath, outDir string) error {

    fileSet := token.NewFileSet()
    file, err := parser.ParseFile(fileSet, filepath, nil, 0)
    if nil != err {
        return err
    }
    if !ast.FileExports(file) {
        return fmt.Errorf("no exported members")
    }
    contents := convertToPython(filepath, file)

    filename := path.Join(outDir, path.Base(filepath))
    filename = setExtension(filename, ".c")
    f, err := os.OpenFile(
        filename,
        os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
        os.ModePerm,
    )
    if nil != err {
        return err
    }
    defer f.Close()

    _, err = f.Write([]byte(contents))
    if nil != err {
        return err
    }
    fmt.Printf("wrote %s\n", filename)
    fileCount++
    return nil
}

func setExtension(filename, ext string) string {
    parts := strings.Split(filename, ".")
    if len(parts) == 1 {
        return filename + ext
    }
    parts = parts[:len(parts)-1]
    return strings.Join(parts, ".") + ext
}

func convertToPython(name string, file *ast.File) string {
    fmt.Printf("parsing %s...\n", name)

    _ := NewFileMap(file, name)
    return fileHeader

    /*contents := ""

      for _, decl := range file.Decls {
          switch decl.(type) {
          case *ast.GenDecl:
              genDecl := decl.(*ast.GenDecl)
              if genDecl.Tok == token.TYPE {
                  contents += convertTypeDecl(genDecl)
              }
          }
      }
      return contents*/
}

// FileMap ...
type FileMap struct {
    Name  string
    Types map[string]*TypeMap
}

// NewFileMap ...
func NewFileMap(file *ast.File, name string) *FileMap {
    fm := &FileMap{
        name,
        make(map[string]*TypeMap),
    }

    ast.Walk(fm, file)
    return fm
}

// Visit ...
func (fm *FileMap) Visit(n ast.Node) ast.Visitor {
    fmt.Printf("l")
    return fm
}

type TypeMap struct {
}

func convertTypeDecl(genDecl *ast.GenDecl) string {
    for _, spec := range genDecl.Specs {
        switch spec.(type) {
        case *ast.TypeSpec:
            tVisitor := &typeVisitor{}
            ast.Walk(tVisitor, spec)
            return tVisitor.String()
        }
    }
    return ""
}

type typeVisitor struct {
    Type      fmt.Stringer
    lastIdent string
}

func (tv *typeVisitor) String() string {
    if nil == tv.Type {
        return fmt.Sprintf(
            "\n#warning skipped type in transpiling: %s\n", tv.lastIdent)
    }
    return tv.Type.String()
}

func (tv *typeVisitor) Visit(n ast.Node) ast.Visitor {
    switch n.(type) {

    case *ast.TypeSpec:
        return tv

    case *ast.InterfaceType:
        fmt.Printf("├ type: interface (skip).\n")
        return nil

    case *ast.SelectorExpr:
        tv.Type = &vanillaTypeVisitor{
            tv.lastIdent,
            expressionToString(n.(*ast.SelectorExpr)),
        }
        return nil

    case *ast.StructType:
        v := &structVisitor{}
        v.Name = tv.lastIdent
        tv.Type = v
        return v

    case *ast.Ident:
        tv.lastIdent = n.(*ast.Ident).String()
        return nil

    case nil:
        return nil
    default:
        fmt.Printf("├ tv unhandled %s\n", reflect.TypeOf(n))
        return nil
    }
}

type vanillaTypeVisitor struct {
    Name string
    Type string
}

func (vtv *vanillaTypeVisitor) String() string {
    return fmt.Sprintf(
        "\ntypedef %s %s // %s\n",
        toCType(vtv.Type), vtv.Name, vtv.Type)
}

type structVisitor struct {
    Name    string
    Type    string
    Members []*memberVisitor
}

func (sv *structVisitor) String() string {
    str := fmt.Sprintf(
        `
typedef struct %s {
    PyObject_HEAD
`, sv.Name)
    for _, mem := range sv.Members {
        str += fmt.Sprintf("\t%s %s // %s\n", mem.CType(), mem.Name, mem.Type)
    }
    return str + "}\n"
}

func (sv *structVisitor) Visit(n ast.Node) ast.Visitor {
    switch n.(type) {
    case *ast.FieldList:
        return sv
    case *ast.Field:
        field := n.(*ast.Field)
        mem := &memberVisitor{}

        if len(field.Names) == 0 {
            mem.Unnamed = true
        }
        sv.Members = append(sv.Members, mem)
        return mem
    case nil:
        return nil
    default:
        fmt.Printf("│├ sv unhandled %s\n", reflect.TypeOf(n))
        return nil
    }
}

type memberVisitor struct {
    Name    string
    Type    string
    Pointer bool
    Unnamed bool
}

func (mv *memberVisitor) Visit(n ast.Node) ast.Visitor {
    switch n.(type) {
    case *ast.StarExpr:
        mv.Pointer = true
        mv.Type += "*"
        return mv
    case *ast.Ident:
        if "" == mv.Name && !mv.Unnamed {
            mv.Name = n.(*ast.Ident).String()
        } else {
            mv.Type += n.(*ast.Ident).String()
        }
        return nil
    case *ast.SelectorExpr:
        exp := n.(*ast.SelectorExpr)
        mv.Type += expressionToString(exp)
        return nil
    case nil:
        return nil
    default:
        fmt.Printf("││├ mv unhandled %s\n", reflect.TypeOf(n))
        return nil
    }
}

func (mv *memberVisitor) CType() string {
    if mv.Unnamed {
        return ""
    }
    return toCType(mv.Type)
}

func toCType(t string) string {
    switch t {
    case "string":
        return "char*"
    }
    return t
}

func expressionToString(exp ast.Expr) string {
    var members []string
    joiner := ""
    switch exp.(type) {
    case *ast.SelectorExpr:
        joiner = "."
    case nil:
        return ""
    default:
        fmt.Printf("exp2str unhandled exp type %s\n", reflect.TypeOf(exp))
        return ""
    }
    ast.Inspect(exp, func(n ast.Node) bool {
        if n == exp {
            return true
        }
        switch n.(type) {
        case *ast.Ident:
            members = append(members, n.(*ast.Ident).String())
        case *ast.SelectorExpr:
            members = append(members,
                expressionToString(n.(*ast.SelectorExpr)))
        case nil:
            break
        default:
            fmt.Printf("exp2str unhandled %s\n", reflect.TypeOf(n))
        }
        return true
    })
    return strings.Join(members, joiner)
}
