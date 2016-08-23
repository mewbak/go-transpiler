// auto-generated file, do not edit

#define Py_LIMITED_API
#include <Python.h>
#include "structmember.h"

{{range .Package.Types}}
{{- if .Name}}
#include "{{.Name}}.h"
{{- end}}
{{- end}}

{{if .Name}}{{template "cStruct.tpl" .}}{{end}}