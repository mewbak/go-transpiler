{{if .Struct -}}

var cache{{.Name}} = make(map[int64]*{{.Package.Name}}.{{.Name}})

//export create{{.Name}}
func create{{.Name}}(
    {{- range .NamedMembers}}
    {{.Name}} {{goIncomingArgType .Type}},
    {{- end}}
) int64 {

    new{{.Name}} := &{{.Package.Name}}.{{.Name}}{
        {{- range .NamedMembers}}
        {{.Name}}: {{convertFromCValue .Type .Name}},
        {{- end}}
    }

    id := time.Now().Unix()
    cache{{.Name}}[id] = new{{.Name}}
    return id
}

//export free{{.Name}}
func free{{.Name}}(ref int64) {
    delete(cache{{.Name}}, ref)
}

{{- else -}}{{/*if .Interface*/}}
{{- /*interfaces act like regular types for python but have no data*/ -}}

var cache{{.Name}} = make(map[int64]{{.Package.Name}}.{{.Name}})

//export create{{.Name}}
func create{{.Name}}() int64 { return 0 }

//export free{{.Name}}
func free{{.Name}}(ref int64) {}

{{- end}}