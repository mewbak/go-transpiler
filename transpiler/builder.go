package transpiler

type Builder interface {
    Build(file *FileMap, outDir string)
}
