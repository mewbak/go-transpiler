{{range .NamedMembers -}}

extern void go{{$.Name}}Set_{{.Name}}(long long cacheKey, {{cTransitionType .}} val);
extern {{cTransitionType .}} go{{$.Name}}Get_{{.Name}}(long long cacheKey);

static int
{{$.Name}}Set_{{.Name}}({{$.Name}} *self, PyObject *value, void *closure)
{
    if (value == NULL || value == Py_None) {
        PyErr_SetString(PyExc_TypeError, "Cannot delete the {{camelToSnake .Name}} attribute");
        return -1;
    }

    if (!{{validatePyValue . "value"}}) {
        PyErr_SetString(PyExc_TypeError, "invalid type assignment for attribute {{camelToSnake .Name}}");
        return -1;
    }

    {{cTransitionType .}} val = {{convertPyToC . "value"}};
    go{{$.Name}}Set_{{.Name}}(self->go{{$.Name}}, val);
    return 0;

}

static PyObject*
{{$.Name}}Get_{{.Name}}({{$.Name}}* self, void *closure)
{
    {{cTransitionType .}} val = go{{$.Name}}Get_{{.Name}}(self->go{{$.Name}});
    PyObject *obj = {{convertPyFromC . "val"}};
    return obj;
}

{{end -}}

static PyGetSetDef {{.Name}}_getseters[] = {
    {{- range .NamedMembers}}
    {
        "{{camelToSnake .Name}}",
        (getter){{$.Name}}Get_{{.Name}},
        (setter){{$.Name}}Set_{{.Name}},
        "", //TODO docstring
        NULL, //this should always be NULL as it's assumed as such in {{.Name}}_init()
    },
    {{- end}}
    {NULL}
};