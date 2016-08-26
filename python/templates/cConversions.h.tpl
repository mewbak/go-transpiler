#define Py_LIMITED_API
#include <Python.h>
#include "structmember.h"
#include "datetime.h"

#ifndef CONVERSIONS_H
#define CONVERSIONS_H

{{range .}}
{{- .CDeclarations}}
{{- end}}

#endif