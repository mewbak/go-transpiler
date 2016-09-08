package transpiler

import (
    "fmt"
    "go/ast"
    "reflect"
    "strings"
)

// ExpressionToString ...
func ExpressionToString(exp ast.Expr) string {
    var members []string
    joiner := ""
    switch e := exp.(type) {

    case *ast.SelectorExpr:
        joiner = "."

    case *ast.InterfaceType:
        return "interface{}" // BUG(rydrman): should handle field list

    case *ast.Ident:
        return e.String()

    case *ast.StarExpr:
        members = append(members, "*")
        break

    case nil:
        return ""

    default:
        fmt.Printf("exp2str unhandled exp type %s\n", reflect.TypeOf(exp))
        return ""

    }
    ast.Inspect(exp, func(n ast.Node) bool {
        if n == exp {
            return true
        }
        switch n.(type) {

        case *ast.Ident:
            members = append(members, n.(*ast.Ident).String())

        case *ast.SelectorExpr:
            members = append(members,
                ExpressionToString(n.(*ast.SelectorExpr)))

        case *ast.InterfaceType:
            members = append(members, "interface{}") // BUG(rydrman): should handle field list

        case nil:
            break

        default:
            fmt.Printf("exp2str unhandled %s\n", reflect.TypeOf(n))

        }

        return true

    })

    return strings.Join(members, joiner)
}
