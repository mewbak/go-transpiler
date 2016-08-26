{{$supportedFunctions := filterSupportedFunctions .Functions -}}
/* BEGIN TYPE FUNCTIONS */
{{- range $supportedFunctions}}

static PyObject* {{$.Name}}_{{.Name}}({{$.Name}} *self, PyObject *args) {
    {{if len .Params -}}
    {{range $i, $_ := .Params}}
    PyObject* arg_{{.Name}} = NULL;
    {{- end}}

    if (!PyArg_ParseTuple(
        args, "|{{range .Params}}O{{end}}",
        {{- $p := .Params}}
        {{- range $i, $_ := .Params}}
        &arg_{{.Name}}{{if notLast $i $p}},{{end}}
        {{- end}})) {
        return NULL;
    }
    {{range .Params}}
    if (!{{validatePyValue .Type (print "arg_" .Name)}}) {}
    {{- end}}
    {{- end}}
    return Py_None; //FIXME
}
{{end}}

static PyMethodDef {{.Name}}_methods[] = {
    {{- range $supportedFunctions}}
    {
        "{{camelToSnake .Name}}",
        (PyCFunction){{$.Name}}_{{.Name}},
        {{if (len .Params)}}METH_VARARGS{{else}}METH_NOARGS{{end}},
        "" //TODO docstring generation
    },
    {{- end}}
    {NULL}
};

/* END TYPE FUNCTIONS */