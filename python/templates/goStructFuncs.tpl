{{$supportedFunctions := filterSupportedFunctions .Functions -}}
/* BEGIN TYPE FUNCTIONS */
{{- range $supportedFunctions}}
{{- $p := .Params}}
{{- $r := .Results}}

//export go{{$.Name}}_{{.Name}}
func go{{$.Name}}_{{.Name}}(
    cacheKey int64,
    {{- range $i, $_ := .Params}}
    {{.Name}} {{goTransitionType .Type}},
    {{- end}}
) *C.PyObject {

    {{range .Params}}
    arg{{.Name}} := {{convertGoFromC .Type .Name}}
    {{- end}}

    self := cache{{$.Name}}[cacheKey]

    {{range $i, $_ := .Results -}}
    res{{print $i}}{{if notLast $i $r}},{{end}}
    {{- end}} := self.{{.Name}}(
        {{- range .Params}}
        arg{{.Name}},
        {{- end}}
    )

    {{range $i, $_ := .Results}}
    cRes{{print $i}} := {{convertGoToC .Type (print "res" $i)}}
    {{- end}}

    return C.{{$.Name}}_{{.Name}}_BuildResult(
        {{- range $i, $_ := .Results}}
        cRes{{print $i}},
        {{- end}}
    )

}
{{end}}
/* END TYPE FUNCTIONS */