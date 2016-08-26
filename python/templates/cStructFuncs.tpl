{{- range .Functions}}

//{{- .Name}}() {}

{{end}}

static PyMethodDef {{.Name}}_methods[] = {
    {{- range .Functions}}
    {
        "{{camelToSnake .Name}}",
        NULL, //TODO
        0, //TODO
        "" //TODO docstring generation
    },
    {{- end}}
    {NULL}
};
