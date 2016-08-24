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

    // start with type definitions
    for _, tm := range pm.Types {

        if tm.Name == "" {
            fmt.Printf("skipping unnamed type...")
            continue
        }

        fmt.Printf("(%s) transpiling ...\n", tm.Name)

        cfiles, err := b.writeCType(tm)
        if nil != err {
            return nil, err
        }

        gofile, err := b.writeGoType(tm)
        if nil != err {
            return nil, err
        }

        created = append(created, gofile)
        created = append(created, cfiles...)
    }

    // add global module files
    mFiles, err := b.writeModuleFiles()
    if err != nil {
        return nil, err
    }
    created = append(created, mFiles...)

    return created, nil

}

func (b *Builder) writeCType(tm *transpiler.TypeMap) ([]string, error) {

    filename := path.Join(b.out, path.Base(tm.Name))
    cFileName := setExtension(filename, ".c")
    hFileName := setExtension(filename, ".h")

    hf, err := os.OpenFile(hFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return nil, err
    }
    defer hf.Close()

    err = cTemplate.ExecuteTemplate(hf, "cHeader.tpl", tm)
    if nil != err {
        return nil, err
    }

    cf, err := os.OpenFile(cFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return nil, err
    }
    defer cf.Close()

    err = cTemplate.Execute(cf, tm)
    if nil != err {
        return nil, err
    }

    return []string{hFileName, cFileName}, nil
}

func (b *Builder) writeGoType(tm *transpiler.TypeMap) (string, error) {

    filename := path.Join(b.out, path.Base(tm.Name))
    filename = setExtension(filename, ".go")

    f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return "", err
    }
    defer f.Close()

    err = goTemplate.Execute(f, tm)
    if nil != err {
        return "", err
    }
    f.Close()

    err = runGoImports(filename)
    if nil != err {
        return filename, err
    }

    return filename, nil
}

func (b *Builder) writeModuleFiles() ([]string, error) {

    var created []string
    cConversions := path.Join(b.out, "type_conversions.c")

    cf, err := os.OpenFile(cConversions, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return created, err
    }
    defer cf.Close()

    err = cTemplate.ExecuteTemplate(cf, "cConversions.tpl", converters)
    if nil != err {
        return created, err
    }
    created = append(created, cConversions)

    cConversionsH := path.Join(b.out, "type_conversions.h")

    chf, err := os.OpenFile(cConversionsH, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return created, err
    }
    defer chf.Close()

    err = cTemplate.ExecuteTemplate(chf, "cConversions.h.tpl", converters)
    if nil != err {
        return created, err
    }
    created = append(created, cConversionsH)

    goConversions := path.Join(b.out, "type_conversions.go")

    gf, err := os.OpenFile(goConversions, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return created, err
    }
    defer gf.Close()

    err = goTemplate.ExecuteTemplate(gf, "goConversions.tpl", converters)
    if nil != err {
        return created, err
    }
    created = append(created, goConversions)

    err = runGoImports(goConversions)
    if nil != err {
        return created, err
    }

    mainFile, err := b.buildMainFile()
    if nil != err {
        return created, err
    }
    created = append(created, mainFile)

    return created, nil

}

func (b *Builder) buildMainFile() (string, error) {

    filename := path.Join(b.out, "main.go")

    f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return "", err
    }
    defer f.Close()

    err = goTemplate.ExecuteTemplate(f, "goMain.tpl", b.m)
    if nil != err {
        return filename, err
    }

    err = runGoImports(filename)
    if nil != err {
        return filename, err
    }

    return filename, nil

}

func runGoImports(f string) error {
    abs, _ := filepath.Abs(f)
    fmt.Println("goimports", "-w", abs)
    cmd := exec.Command("goimports", "-w", abs)
    out, err := cmd.CombinedOutput()
    if nil != err {
        return fmt.Errorf("error calling goimports: %s\n%s", err, out)
    }
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
