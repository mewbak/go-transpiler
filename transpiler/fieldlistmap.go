package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
)

// FieldListMap ...
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
