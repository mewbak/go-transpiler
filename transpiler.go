// command line wrapper for transpiler package
package main

import (
    "flag"
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "os"
    "path"
    "strings"

    "github.com/KloudKtrl/go-transpiler/python"
    "github.com/KloudKtrl/go-transpiler/transpiler"
)

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

    fmt.Printf("parsing %s...\n", filepath)
    if strings.Contains(filepath, "_test.go") {
        return fmt.Errorf("transpiling tests is not supported")
    }

    fileSet := token.NewFileSet()
    file, err := parser.ParseFile(fileSet, filepath, nil, 0)
    if nil != err {
        return err
    }
    if !ast.FileExports(file) {
        return fmt.Errorf("no exported members")
    }

    mapped := transpiler.NewFileMap(file, filepath)
    pyBuilder := &python.Builder{}
    output, err := pyBuilder.Build(mapped, outDir)
    if nil != err {
        return err
    }

    fmt.Printf("> wrote %s\n", output)
    fileCount++
    return nil

}
