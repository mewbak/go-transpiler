// auto-generated file, do not edit

#define Py_LIMITED_API
#include <Python.h>
#include "structmember.h"

{{range .Types}}
{{if .Name}}{{template "cStruct.tpl" .}}{{end}}
{{end}}