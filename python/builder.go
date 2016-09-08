package python

import (
    "fmt"
    "os"
    "os/exec"
    "path"
    "path/filepath"
    "reflect"
    "strings"

    "github.com/KloudKtrl/go-transpiler/transpiler"
)

var pyModuleName string

// Builder ...
type Builder struct {
    out string
    m   *transpiler.PackageMap
}

// Build ...
func (b *Builder) Build(pm *transpiler.PackageMap, outDir, name string) ([]string, error) {

    b.out = outDir
    b.m = pm
    pyModuleName = name

    if "" == pyModuleName {
        return nil, fmt.Errorf(
            "cannot build module without name")
    }

    var created []string

    // start with type definitions
    for _, tm := range pm.Types {

        if tm.Name == "" {
            fmt.Println("skipping unnamed type...")
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

    // filter converters to only one of each reflect.Type
    var cTypes []reflect.Type
    var convertersUnique []converter
    for _, cvtr := range converters {
        cType, found := reflect.TypeOf(cvtr), false
        for _, t := range cTypes {
            if t == cType {
                found = true
                break
            }
        }
        if !found {
            convertersUnique = append(convertersUnique, cvtr)
            cTypes = append(cTypes, cType)
        }
    }

    cConversions := path.Join(b.out, "type_conversions.c")

    cf, err := os.OpenFile(cConversions, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return created, err
    }
    defer cf.Close()

    err = cTemplate.ExecuteTemplate(cf, "cConversions.tpl", convertersUnique)
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

    err = cTemplate.ExecuteTemplate(chf, "cConversions.h.tpl", convertersUnique)
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

    err = goTemplate.ExecuteTemplate(gf, "goConversions.tpl", convertersUnique)
    if nil != err {
        return created, err
    }
    created = append(created, goConversions)

    err = runGoImports(goConversions)
    if nil != err {
        return created, err
    }

    mainFiles, err := b.buildMainFiles()
    if nil != err {
        return created, err
    }
    created = append(created, mainFiles...)

    return created, nil

}

func (b *Builder) buildMainFiles() ([]string, error) {

    var created []string
    filename := path.Join(b.out, "main.go")

    gf, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return created, err
    }
    defer gf.Close()
    created = append(created, filename)

    err = goTemplate.ExecuteTemplate(gf, "goMain.tpl", b.m)
    if nil != err {
        return created, err
    }

    err = runGoImports(filename)
    if nil != err {
        return created, err
    }

    filename = path.Join(b.out, "main.c")

    cf, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return created, err
    }
    defer gf.Close()
    created = append(created, filename)

    err = cTemplate.ExecuteTemplate(cf, "cMain.tpl", b.m)
    if nil != err {
        return created, err
    }

    return created, nil

}

func runGoImports(f string) error {
    abs, _ := filepath.Abs(f)
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
