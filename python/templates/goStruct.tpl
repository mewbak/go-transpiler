{{if .Struct -}}

var cache{{.Name}} = make(map[int64]*{{.Package.Name}}.{{.Name}})

func getCached{{.Name}}(item *{{.Package.Name}}.{{.Name}}) int64 {
    for key, cached := range cache{{.Name}} {
        if cached == item {
            return key
        }
    }
    id := time.Now().Unix()
    cache{{.Name}}[id] = item
    return id
}

//export create{{.Name}}
func create{{.Name}}() int64 {

    new{{.Name}} := &{{.Package.Name}}.{{.Name}}{}

    id := time.Now().Unix()
    cache{{.Name}}[id] = new{{.Name}}
    return id
}

//export free{{.Name}}
func free{{.Name}}(ref int64) {
    delete(cache{{.Name}}, ref)
}

{{template "goGetSet.tpl" .}}

{{template "goStructFuncs.tpl" .}}

{{- else -}}{{/*if .Interface*/}}
{{- /*
    interfaces act like regular types for python but
    have no associated data or functionality
*/ -}}

var cache{{.Name}} = make(map[int64]{{.Package.Name}}.{{.Name}})

func getCached{{.Name}}(item {{.Package.Name}}.{{.Name}}) int64 {
    return 0
}

//export create{{.Name}}
func create{{.Name}}() int64 { return 0 }

//export free{{.Name}}
func free{{.Name}}(ref int64) {}

{{- end}}