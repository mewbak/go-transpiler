package python

import (
    "fmt"
    "os"
    "os/exec"
    "path"
    "path/filepath"
    "strings"

    "github.com/KloudKtrl/go-transpiler/transpiler"
)

// Builder ...
type Builder struct {
    out string
    m   *transpiler.PackageMap
}

// Build ...
func (b *Builder) Build(pm *transpiler.PackageMap, outDir string) ([]string, error) {

    b.out = outDir
    b.m = pm

    var created []string

    for _, fm := range pm.Files {

        cfile, err := b.writeCFile(fm)
        if nil != err {
            return nil, err
        }

        gofile, err := b.writeGoFile(fm)
        if nil != err {
            return nil, err
        }

        created = append(created, cfile, gofile)
    }

    return created, nil

}

func (b *Builder) writeCFile(fm *transpiler.FileMap) (string, error) {

    filename := path.Join(b.out, path.Base(fm.Name))
    filename = setExtension(filename, ".c")

    f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return "", err
    }
    defer f.Close()

    err = cTemplate.Execute(f, fm)
    if nil != err {
        return "", err
    }

    return filename, nil
}

func (b *Builder) writeGoFile(fm *transpiler.FileMap) (string, error) {

    filename := path.Join(b.out, path.Base(fm.Name))
    filename = setExtension(filename, ".go")

    f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return "", err
    }
    defer f.Close()

    err = goTemplate.Execute(f, fm)
    if nil != err {
        return "", err
    }
    f.Close()

    // run go imports so that we don't have to manage that
    // in the templates
    abs, _ := filepath.Abs(filename)
    fmt.Println("goimports", "-w", abs)
    cmd := exec.Command("goimports", "-w", abs)
    err = cmd.Run()
    if nil != err {
        return filename, fmt.Errorf("error calling goimports: %s", err)
    }

    return filename, nil
}

func setExtension(filename, ext string) string {

    parts := strings.Split(filename, ".")

    if len(parts) == 1 {
        return filename + ext
    }

    parts = parts[:len(parts)-1]

    return strings.Join(parts, ".") + ext

}
