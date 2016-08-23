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

        pkgMap := transpiler.NewPackageMap(pkgName)
        for name := range pkg.Files {

            f, err := transpileFile(name)
            if nil != err {
                fmt.Printf("%s\n", err)
            } else {
                pkgMap.AddFile(f)
            }
        }
        pyBuilder := &python.Builder{}
        output, err := pyBuilder.Build(pkgMap, outDir)
        if nil != err {
            fmt.Println(err)
        } else {
            fmt.Printf("wrote %d files:\n", len(output))
            for _, name := range output {
                fmt.Printf("  %s\n", name)
            }
        }
        pkgCount++
    }
    return nil

}

func transpileFile(filepath string) (*transpiler.FileMap, error) {

    fmt.Printf("> parsing %s...\n", filepath)
    if strings.Contains(filepath, "_test.go") {
        return nil, fmt.Errorf("transpiling tests is not supported")
    }

    fileSet := token.NewFileSet()
    file, err := parser.ParseFile(fileSet, filepath, nil, 0)
    if nil != err {
        return nil, err
    }
    if !ast.FileExports(file) {
        return nil, fmt.Errorf("no exported members")
    }

    mapped := transpiler.NewFileMap(file, filepath)
    fileCount++

    return mapped, nil

}
