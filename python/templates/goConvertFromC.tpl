{{if eq .Type "string" -}}

C.GoString({{.Name}})

{{- else if isInternalType .Type -}}

{{.Type}}({{.Name}})

{{- else -}}

{{.Name}}

{{- end}}