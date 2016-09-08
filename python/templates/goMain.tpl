{{- $supportedFunctions := filterSupportedFunctions .Functions}}
package main

/*
#cgo pkg-config: python2
#define Py_LIMITED_API
#include <Python.h>
#include "type_conversions.h"
{{- range $supportedFunctions}}
{{- $r := .Results}}
{{- if not .Receiver}}

PyObject *{{$.Name}}_{{.Name}}_BuildResult(
    {{- range $i, $_ := .Results}}
    {{cTransitionType .}} res{{print $i}}{{if notLast $i $r}},{{end}}
    {{- end}});
{{- end}}
{{- end}}
*/
import "C"

{{template "goFuncs.tpl" .}}

func main() {}