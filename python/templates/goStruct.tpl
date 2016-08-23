var cache{{.Name}} = make(map[int64]*{{.Package.Name}}.{{.Name}})

//export create{{.Name}}
func create{{.Name}}(
    {{- if .Members}}
        {{- range .Members}}
            {{- if .Name}}
    {{.Name}} {{goCType .Type}},
            {{- end}}
        {{- end}}
    {{- end}}
) int64 {

    new{{.Name}} := &{{.Package.Name}}.{{.Name}}{
    {{- if .Members}}
        {{- range .Members}}
            {{- if .Name}}
        {{.Name}}: {{template "goConvertFromC.tpl" .}},
            {{- end}}
        {{- end}}
    {{- end}}
    }

    id := time.Now().Unix()
    cache{{.Name}}[id] = new{{.Name}}
    return id
}