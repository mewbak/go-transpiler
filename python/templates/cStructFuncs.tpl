{{- range .Functions}}

static PyObject* {{$.Name}}_{{.Name}}({{$.Name}} *self) {
    return Py_None; //FIXME
}

{{end}}

static PyMethodDef {{.Name}}_methods[] = {
    {{- range .Functions}}
    {
        "{{camelToSnake .Name}}",
        (PyCFunction){{$.Name}}_{{.Name}},
        METH_NOARGS, //TODO
        "" //TODO docstring generation
    },
    {{- end}}
    {NULL}
};
