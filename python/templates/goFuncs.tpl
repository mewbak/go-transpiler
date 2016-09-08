{{$supportedFunctions := filterSupportedFunctions .Functions -}}
/**BEGIN GLOBAL FUNCS**/
{{- range $supportedFunctions}}
{{- $p := .Params}}
{{- $r := .Results}}
{{- if not .Receiver}}

//export go_{{.Name}}
func go_{{.Name}}(
    {{- range $i, $_ := .Params}}
    {{.Name}} {{goTransitionType .}},
    {{- end}}
) *C.PyObject {

    {{range .Params}}
    arg{{.Name}} := {{convertGoFromC . .Name}}
    {{- end}}

    {{range $i, $_ := .Results -}}
    res{{print $i}}{{if notLast $i $r}},{{end}}
    {{- end}} := {{.Package.Name}}.{{.Name}}(
        {{- range .Params}}
        arg{{.Name}},
        {{- end}}
    )

    {{range $i, $_ := .Results}}
    cRes{{print $i}} := {{convertGoToC . (print "res" $i)}}
    {{- end}}

    return C.{{$.Name}}_{{.Name}}_BuildResult(
        {{- range $i, $_ := .Results}}
        {{convertGoParamForCFunc . (print "cRes" $i)}},
        {{- end}}
    )

}
{{end}}
{{- end}}
/**END GLOBAL FUNCS**/