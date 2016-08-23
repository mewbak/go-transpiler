package transpiler

// PackageMap maps the information conatined within
// a go package
type PackageMap struct {
    Name  string
    Files []*FileMap
}

// NewPackageMap creates a new package map
// for the package with the given name
func NewPackageMap(name string) *PackageMap {
    return &PackageMap{
        Name: name,
    }
}

// AddFile appends the given file to this package
func (pm *PackageMap) AddFile(f *FileMap) {
    f.SetPackage(pm)
    pm.Files = append(pm.Files, f)
}
