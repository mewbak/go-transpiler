package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
)

// FieldMap maps a field definition from a go abstract
// syntax tree. The FieldMap is intended to visit ast.Field
// nodes and their children
type FieldMap struct {

    // Name is the name of this field in a struct or interface
    // and can be blank if Unnamed = true
    Name string

    // Type is the qualified typename for this field
    // eg "package.Type"
    Type string

    // TypeName is the singular name of this field type
    // eg "Type" for "package.Type" fields
    TypeName string

    // TypeExpr is the full expression denoting this type
    // eg if Pointer=true this might be "*package.Type"
    TypeExpr string

    // KeyType is only set when Map=true and denotes the
    // key type for this map eg "int" for map[int]string
    KeyType string

    // ValueType is set for Map, Array, or Slice = true and
    // denotes the type of values stored within
    // eg "string" for map[int]string
    ValueType string

    // These booleans describe the type of field
    Unnamed bool
    Pointer bool
    Array   bool
    Slice   bool
    Map     bool
}

// NewFieldMap creates a new, empty FieldMap
func NewFieldMap() *FieldMap {
    return &FieldMap{
        Name:      "",
        Type:      "",
        TypeName:  "",
        TypeExpr:  "",
        KeyType:   "",
        ValueType: "",
        Unnamed:   false,
        Pointer:   false,
        Map:       false,
    }
}

// Visit ...
func (fm *FieldMap) Visit(n ast.Node) ast.Visitor {

    switch node := n.(type) {

    case *ast.StarExpr:
        fm.Pointer = true
        fm.TypeExpr += "*"
        return fm

    case *ast.Ident:
        if "" == fm.Name && !fm.Unnamed {
            fm.Name = node.String()
        } else {
            fm.Type += node.String()
            fm.TypeExpr += node.String()
            fm.TypeName = node.String()
        }
        return nil

    case *ast.BasicLit:
        // often string tags: ie `json:"name"`
        return nil

    case *ast.MapType:
        fm.KeyType = ExpressionToString(node.Key)
        fm.ValueType = ExpressionToString(node.Value)
        fm.Map = true
        return nil

    case *ast.SelectorExpr:
        expr := ExpressionToString(node)
        fm.Type += expr
        fm.TypeExpr += expr
        return nil

    case nil:
        return nil

    default:
        fmt.Printf("FieldMap unhandled %s\n", reflect.TypeOf(n))
        return nil

    }

}
