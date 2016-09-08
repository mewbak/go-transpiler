// auto-generated file, do not edit

package main

/*
#cgo pkg-config: python2
#define Py_LIMITED_API
#include <Python.h>
#include "type_conversions.h"
{{- if .Name}}
{{- $supportedFunctions := filterSupportedFunctions .Functions}}
{{- range $supportedFunctions}}
{{- $r := .Results}}

PyObject *{{$.Name}}_{{.Name}}_BuildResult(
    {{- range $i, $_ := .Results}}
    {{cTransitionType .}} res{{print $i}}{{if notLast $i $r}},{{end}}
    {{- end}});
{{- end}}
{{- end}}
*/
import "C"

{{if .Name}}{{template "goStruct.tpl" .}}{{end}}