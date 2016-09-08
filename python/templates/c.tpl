// auto-generated file, do not edit

#define Py_LIMITED_API
#include <Python.h>
#include "structmember.h"
#include "datetime.h"
#include "type_conversions.h"

{{range .Package.Types}}
{{- if .Name}}
#include "{{.Name}}.h"
{{- end}}
{{- end}}

{{if .Name}}{{template "cType.tpl" .}}{{end}}