// auto-generated file, do not edit

package main

/*
#cgo pkg-config: python2
#define Py_LIMITED_API
#include <Python.h>

*/
import "C"

import ()

{{range .Types}}
{{template "goStruct.tpl" .}}
{{end}}