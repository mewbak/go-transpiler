package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
)

// FieldListMap maps a list of fields from a go abstract
// syntax tree. The fieldlistmap is intended to visit ast.FieldList
// nodes and their children
type FieldListMap []*FieldMap

// NewFieldListMap creates an empty field list map
func NewFieldListMap() *FieldListMap {
    return &FieldListMap{}
}

// Visit ...
func (flm *FieldListMap) Visit(n ast.Node) ast.Visitor {

    switch node := n.(type) {

    case *ast.FieldList:
        return flm

    case *ast.Field:
        field := NewFieldMap()
        if len(node.Names) == 0 {
            field.Unnamed = true
        }
        // multiple names share the type (one, two string)
        for i, name := range node.Names {
            if i == len(node.Names)-1 {
                break
            }
            f := NewFieldMap()
            f.Name = name.String()
            flm.Add(f)
        }
        flm.Add(field)
        return field

    case nil:
        return nil

    default:
        fmt.Printf("FieldListMap unhandled %s\n", reflect.TypeOf(n))
        return nil

    }

}

// Add is a shotform to append a FieldMap to this list/slice
func (flm *FieldListMap) Add(field *FieldMap) {
    (*flm) = append(*flm, field)
}

// Count is a shotform to get the length of this list
func (flm *FieldListMap) Count() int {
    return len(*flm)
}

// Finalize ...
func (flm *FieldListMap) Finalize() {

    for _, fm := range *flm {
        fm.Finalize()
    }

}
