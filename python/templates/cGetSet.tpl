{{range .NamedMembers -}}

extern void go{{$.Name}}_Set{{.Name}}(long long cacheKey, {{cTransitionType .Type}} val);
extern {{cTransitionType .Type}} go{{$.Name}}_Get{{.Name}}(long long cacheKey);

static int
{{$.Name}}_Set{{.Name}}({{$.Name}} *self, PyObject *value, void *closure)
{
    if (value == NULL) {
        PyErr_SetString(PyExc_TypeError, "Cannot delete the {{camelToSnake .Name}} attribute");
        return -1;
    }

    if (!{{validatePyValue .Type "value"}}) {
        PyErr_SetString(PyExc_TypeError, "invalid type assignment for attribute {{camelToSnake .Name}}");
        return -1;
    }

    {{cTransitionType .Type}} val = {{convertPyToC .Type "value"}};
    go{{$.Name}}_Set{{.Name}}(self->go{{$.Name}}, val);
    return 0;

}

static PyObject*
{{$.Name}}_Get{{.Name}}({{$.Name}}* self, void *closure)
{
    {{cTransitionType .Type}} val = go{{$.Name}}_Get{{.Name}}(self->go{{$.Name}});
    PyObject *obj = {{convertPyFromC .Type "val"}};
    return obj;
}

{{end -}}

static PyGetSetDef {{.Name}}_getseters[] = {
    {{- range .NamedMembers}}
    {
        "{{camelToSnake .Name}}",
        (getter){{$.Name}}_Get{{.Name}},
        (setter){{$.Name}}_Set{{.Name}},
        "", //TODO docstring
        NULL, //this should always be NULL as it's assumed as such in {{.Name}}_init()
    },
    {{- end}}
    {NULL}
};