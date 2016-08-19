{{if .IsStruct}}
typedef struct {{.Name}} {
{{range .Members}}
    {{.Name}} -> {{.Type}}
{{end}}
}
{{end}}