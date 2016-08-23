// auto-generated file, do not edit

#define Py_LIMITED_API
#include <Python.h>
#include "structmember.h"
#include "datetime.h"

#ifndef {{.Name}}
#define {{.Name}}

typedef struct {
    PyObject_HEAD
    {{- range .NamedMembers}}
    {{cMemberType .Type}} {{.Name}};
    {{- end}}
    long long go{{.Name}};
} {{.Name}};

static PyTypeObject {{.Name}}_type;

#endif