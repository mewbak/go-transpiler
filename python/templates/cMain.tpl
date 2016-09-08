#define Py_LIMITED_API
#include <Python.h>
#include "datetime.h"
#include "structmember.h"
#include "type_conversions.h"

{{range .Types}}
{{- if .Name}}
#include "{{.Name}}.h"
{{- end}}
{{- end}}

{{template "cFuncs.tpl" .}}

PyMODINIT_FUNC init{{pyModuleName}}(void)
{

  PyDateTime_IMPORT;

  {{range .Types}}
  {{- if .Name}}
  if(PyType_Ready(&{{.Name}}_type) < 0) {
    return;
  }
  {{- end}}
  {{- end}}

  PyObject*
  m = Py_InitModule3("{{pyModuleName}}", {{.Name}}Methods, ""); //TODO module docstring

  if (m == NULL)
      return;

  {{range .Types}}
  {{- if .Name}}
  Py_INCREF(&{{.Name}}_type);
  PyModule_AddObject(m, "{{.Name}}", (PyObject *)&{{.Name}}_type);
  {{- end}}
  {{- end}}

}