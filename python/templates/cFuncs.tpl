{{/*
    cFuncs renders methods that have a reviever type.

    Note that golang needs to re-pack results into a python
    tuple if there is more than one argument. For this reason
    there are ..._BuildValue functions for each method. These 
    are declared in go.tpl so that they are visible to golang.
*/ -}}
{{$supportedFunctions := filterSupportedFunctions .Functions -}}
/**BEGIN GLOBAL FUNCS**/

{{- range $supportedFunctions}}
{{- $p := .Params}}
{{- $r := .Results}}
{{- if not .Receiver}}
extern PyObject* go_{{.Name}}(
    {{- range $i, $_ := .Params}}
    {{cTransitionType .}} arg_{{.Name}}{{if notLast $i $p}},{{end}}
    {{- end}}
);

static PyObject* {{$.Name}}_{{.Name}}(PyObject *self, PyObject *args)
{
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
    else if (!{{validatePyValue . (print "arg_" .Name)}}) {
        PyErr_SetString(PyExc_TypeError, "invalid parameter type in position {{$i}}");
        return NULL;
    }
    {{- end}}
    {{- end}}

    {{- range .Params}}
    {{cTransitionType .}} argVal_{{.Name}} = {{convertPyToC . (print "arg_" .Name)}};
    {{- end}}

    PyObject *res = go_{{.Name}}(
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
    {{cTransitionType .}} res{{print $i}}{{if notLast $i $r}},{{end}}
    {{- end}})
{
    {{- if eq 1 (len .Results)}}
    PyObject* res = {{convertPyFromC (index .Results 0) "res0"}};
    return res;

    {{- else if gt (len .Results) 1}}
    {{- range $i, $_ := .Results}}
    PyObject* pyRes{{print $i}} = {{convertPyFromC . (print "res" $i)}};
    {{- end}}
    PyObject* res = PyTuple_Pack(
        {{print (len $r)}},
        {{- range $i, $_ := .Results}}
        pyRes{{print $i}}{{if notLast $i $r}},{{end}}
        {{- end}}
    );

    {{- range $i, $_ := .Results}}
    Py_DECREF(pyRes{{print $i}});
    {{- end}}

    return res;

    {{- else}}
    Py_INCREF(Py_None);
    return Py_None;
    {{- end}}
}
{{end}}
{{- end}}

static PyMethodDef {{.Name}}Methods[] = {
    {{- range $supportedFunctions}}
    {{- if not .Receiver}}
    {
        "{{camelToSnake .Name}}",
        (PyCFunction){{$.Name}}_{{.Name}},
        {{if (len .Params)}}METH_VARARGS{{else}}METH_NOARGS{{end}},
        "" //TODO docstring generation
    },
    {{- end}}
    {{- end}}
    { NULL, NULL, 0, NULL }
};

/**END GLOBAL FUNCS**/