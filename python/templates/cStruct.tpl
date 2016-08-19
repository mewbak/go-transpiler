{{if .IsStruct}}
typedef struct {{.Name}} {
{{range .Fields}}
    {{.Name}}
{{end}}
}
{{end}}