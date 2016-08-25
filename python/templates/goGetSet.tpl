{{range .NamedMembers -}}

//export go{{$.Name}}_Set{{.Name}}
func go{{$.Name}}_Set{{.Name}}(cacheKey int64, val {{goTransitionType .Type}}) {

    goVal := {{convertGoFromC .Type "val"}}
    cache{{$.Name}}[cacheKey].{{.Name}} = goVal 

}

//export go{{$.Name}}_Get{{.Name}}
func go{{$.Name}}_Get{{.Name}}(cacheKey int64) {{goTransitionType .Type}} {

    val := cache{{$.Name}}[cacheKey].{{.Name}}
    return {{convertGoToC .Type "val"}}

}

{{end -}}
