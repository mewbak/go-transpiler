#define Py_LIMITED_API
#include <Python.h>
#include "structmember.h"
#include "datetime.h"

{{range .}}
{{- .CDefinitions}}
{{- end}}
