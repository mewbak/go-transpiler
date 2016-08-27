{{range .NamedMembers -}}

//export go{{$.Name}}Set_{{.Name}}
func go{{$.Name}}Set_{{.Name}}(cacheKey int64, val {{goTransitionType .Type}}) {

    goVal := {{convertGoFromC .Type "val"}}
    cache{{$.Name}}[cacheKey].{{.Name}} = goVal 

}

//export go{{$.Name}}Get_{{.Name}}
func go{{$.Name}}Get_{{.Name}}(cacheKey int64) {{goTransitionType .Type}} {

    val := cache{{$.Name}}[cacheKey].{{.Name}}
    return {{convertGoToC .Type "val"}}

}

{{end -}}
