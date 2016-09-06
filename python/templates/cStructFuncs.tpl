{{/*
    cStructFuncs renders methods that have a reviever type.

    Note that golang needs to re-pack results into a python
    tuple if there is more than one argument. For this reason
    there are ..._BuildValue functions for each method. These 
    are declared in go.tpl so that they are visible to golang.
*/ -}}
{{$supportedFunctions := filterSupportedFunctions .Functions -}}
/* BEGIN TYPE FUNCTIONS */
{{- range $supportedFunctions}}
{{- $p := .Params}}
{{- $r := .Results}}

extern PyObject* go{{$.Name}}_{{.Name}}(
    long long cacheKey{{if len .Params}},{{end}}
    {{- range $i, $_ := .Params}}
    {{cTransitionType .Type}} arg_{{.Name}}{{if notLast $i $p}},{{end}}
    {{- end}}
);

static PyObject* {{$.Name}}_{{.Name}}({{$.Name}} *self, PyObject *args) {
    {{if len .Params}}
    {{- range $i, $_ := .Params}}
    PyObject* arg_{{.Name}} = NULL;
    {{- end}}

    if (!PyArg_ParseTuple(
        args, "|{{range .Params}}O{{end}}",
        {{- range $i, $_ := .Params}}
        &arg_{{.Name}}{{if notLast $i $p}},{{end}}
        {{- end}})) {
        return NULL;
    }
    {{range $i, $_ := .Params}}
    if (NULL == arg_{{.Name}}) {
        PyErr_SetString(PyExc_TypeError, "not enough parameters, expeted {{len $p}}, got {{print $i}}");
        return NULL;
    }
    else if (!{{validatePyValue .Type (print "arg_" .Name)}}) {
        PyErr_SetString(PyExc_TypeError, "invalid parameter type in position {{$i}}");
        return NULL;
    }
    {{- end}}
    {{- end}}

    {{- range .Params}}
    {{cTransitionType .Type}} argVal_{{.Name}} = {{convertPyToC .Type (print "arg_" .Name)}};
    {{- end}}

    PyObject *res = go{{$.Name}}_{{.Name}}(
        self->go{{$.Name}}{{if (len .Params)}},{{end}}
        {{- range $i, $_ := .Params}}
        argVal_{{.Name}}{{if notLast $i $p}},{{end}}
        {{- end}}
    );

    if (res == NULL) {
        Py_INCREF(Py_None);
        return Py_None;
    }

    return res;
}

// to pack tuples if necessary
PyObject *{{$.Name}}_{{.Name}}_BuildResult(
    {{- range $i, $_ := .Results}}
    {{cTransitionType .Type}} res{{print $i}}{{if notLast $i $r}},{{end}}
    {{- end}})
{
    {{- if eq 1 (len .Results)}}
    PyObject* res = {{convertPyFromC (index .Results 0).Type "res0"}};
    return res;

    {{- else if gt 1 (len .Results)}}
    {{- range $i, $_ := .Results}}
    PyObject* pyRes{{print $i}} = {{convertPyFromC .Type (print "res" $i)}};
    {{- end}}
    PyObject* res = PyTuple_Pack(
        {{print (len $r)}},
        {{- range $i, $_ := .Results}}
        pyRes{{print $i}}{{if notLast $i $r}},{{end}}
        {{- end}}
    );
    return res;

    {{- else}}
    Py_INCREF(Py_None);
    return Py_None;
    {{- end}}
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