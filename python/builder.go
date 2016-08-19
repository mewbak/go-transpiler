package python

import (
    "fmt"
    "os"
    "path"
    "strings"

    "github.com/KloudKtrl/go-transpiler/transpiler"
)

// Builder ...
type Builder struct {
    out string
    m   *transpiler.FileMap
}

// Build ...
func (b *Builder) Build(fm *transpiler.FileMap, outDir string) ([]string, error) {

    b.out = outDir
    b.m = fm

    cfile, err := b.writeCFile()
    if nil != err {
        return nil, err
    }

    gofile, err := b.writeGoFile()
    if nil != err {
        return nil, err
    }

    return []string{cfile, gofile}, nil

}

func (b *Builder) writeCFile() (string, error) {

    filename := path.Join(b.out, path.Base(b.m.Name))
    filename = setExtension(filename, ".c")

    f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return "", err
    }
    defer f.Close()

    err = cTemplate.Execute(f, b.m)
    if nil != err {
        return "", err
    }

    return filename, nil
}

func (b *Builder) writeGoFile() (string, error) {

    filename := path.Join(b.out, path.Base(b.m.Name))
    filename = setExtension(filename, ".go")

    f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
    if nil != err {
        return "", err
    }
    defer f.Close()

    for _, t := range b.m.Types {
        f.Write([]byte(fmt.Sprintf("%s\n", t.Name)))
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
