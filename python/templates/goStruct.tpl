var cache{{.Name}} = make(map[int64]*{{.Package.Name}}.{{.Name}})

//export create{{.Name}}
func create{{.Name}}(
    {{- range .NamedMembers}}
    {{.Name}} {{goCFuncArgType .Type}},
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