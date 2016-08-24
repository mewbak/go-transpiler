// auto-generated file, do not edit

#define Py_LIMITED_API
#include <Python.h>
#include "structmember.h"
#include "datetime.h"

#ifndef {{.Name}}_H
#define {{.Name}}_H

typedef struct {
    PyObject_HEAD
    {{- range .NamedMembers}}
    {{cMemberType .Type}} {{.Name}};
    {{- end}}
    long long go{{.Name}};
} {{.Name}};

static PyTypeObject {{.Name}}_type;

#endif