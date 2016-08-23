// auto-generated file, do not edit

#define Py_LIMITED_API
#include <Python.h>
#include "structmember.h"

typedef struct {
    PyObject_HEAD
    {{- range .NamedMembers}}
    {{cType .Type}} {{.Name}};
    {{- end}}
    long long go{{.Name}};
} {{.Name}};

static PyTypeObject {{.Name}}_type;
