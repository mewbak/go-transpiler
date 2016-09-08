package main

/*
#cgo pkg-config: python2
#define Py_LIMITED_API
#include <Python.h>
#include "type_conversions.h"
*/
import "C"

{{range .}}
{{- .GoDefinitions}}
{{- end}}