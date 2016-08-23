// auto-generated file, do not edit

package main

/*
#cgo pkg-config: python2
#define Py_LIMITED_API
#include <Python.h>

*/
import "C"

{{range .Types}}
{{if .Name}}{{template "goStruct.tpl" .}}{{end}}
{{end}}