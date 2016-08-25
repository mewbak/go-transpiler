#define Py_LIMITED_API
#include <Python.h>
#include "structmember.h"

{{range .Types}}
{{- if .Name}}
#include "{{.Name}}.h"
{{- end}}
{{- end}}

static PyMethodDef {{.Name}}Methods[] = {
  { NULL, NULL, 0, NULL }
};

PyMODINIT_FUNC init{{pyModuleName}}(void)
{
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