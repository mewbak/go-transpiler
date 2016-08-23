// create in go
extern long long create{{.Name}}(
{{- if .NamedMembers}}
    {{- range $i, $_ := .NamedMembers}}
    {{cType .Type}} {{.Name}}{{if notLast $i $.NamedMembers}},{{end}}
    {{- end}}
{{- end}}
);

//free go pointer
extern void free{{.Name}}(long long elem);

static PyMemberDef {{.Name}}_members[] = {
{{if .NamedMembers}}{{range .NamedMembers}}
    {
        "{{camelToSnake .Name}}",
        {{pythonMemberType .Type}},
        offsetof({{$.Name}}, {{.Name}}),
        0, //members are read only
        "" //TODO docstring generation
    },
{{end}}{{end}}
    {NULL}
};

static void
{{.Name}}_dealloc({{.Name}} *self)
{
    if (self->go{{.Name}}) {
        free{{.Name}}(self->go{{.Name}});
    }
    self->ob_type->tp_free((PyObject*)self);
}

static PyObject*
{{.Name}}_new(PyTypeObject *type, PyObject *args, PyObject *kwargs)
{
    {{.Name}} *self;

    self = ({{.Name}} *)type->tp_alloc(type, 0);
    if (self != NULL) {
{{if .NamedMembers}}{{range .NamedMembers}}{{if eq .Type "string"}}
        self->{{.Name}} = (char*)malloc(sizeof(char));
        memset(self->{{.Name}}, 0, sizeof(char)); //empty string (one null char)
{{end}}{{end}}{{end}}
    }

    return (PyObject *)self;
}

static int
{{.Name}}_init({{.Name}} *self, PyObject *args, PyObject *kwargs)
{
{{if .NamedMembers}}

{{- range .NamedMembers}}
    {{cType .Type}} {{.Name}};
{{- end}}

    static char *kwlist[] = {
        {{- range $i, $_ := .NamedMembers}}
        "{{camelToSnake .Name}}",
        {{- end}}
        NULL
    };

    if (!PyArg_ParseTupleAndKeywords(
        args, kwargs, "{{pyArgFormat .NamedMembers}}", kwlist
        {{- if len .NamedMembers}},{{end}} 
        {{- range $i, $_ := .NamedMembers}}
        {{if eq (cType .Type) "PyObject*"}}&{{.Name}}_type, {{end -}}
        &{{.Name}}{{if notLast $i $.NamedMembers}},{{end}}
        {{- end}})) {
        return -1;
    }

    long long ref = create{{.Name}}(
        {{- range $i, $_ := .NamedMembers}}
        {{.Name}}{{if notLast $i $.NamedMembers}},{{end}}
        {{- end}}
    );
    if (self->go{{.Name}}) {
        free{{.Name}}(self->go{{.Name}});
    }
    self->go{{.Name}} = ref;

    {{range .NamedMembers}}{{if eq .Type "string"}}
    if (self->{{.Name}} != NULL) {
        free(self->{{.Name}});
    }{{end}}{{end}}

    //FIXME free / deal with already set, non-string vars
    {{range .NamedMembers}}
    self->{{.Name}} = {{.Name}};{{end}}
{{end}}
    return 0;
}

{{template "cStructFuncs.tpl" .}}

static PyTypeObject {{.Name}}_type = {
  PyObject_HEAD_INIT(NULL)
  0,                         //ob_size
  "ktrlpy.{{.Name}}",        //tp_name
  sizeof({{.Name}}),         //tp_basicsize
  0,                         //tp_itemsize
  (destructor){{.Name}}_dealloc,//tp_dealloc
  0,                         //tp_print
  0,                         //tp_getattr
  0,                         //tp_setattr
  0,                         //tp_compare
  0,                         //tp_repr
  0,                         //tp_as_number
  0,                         //tp_as_sequence
  0,                         //tp_as_mapping
  0,                         //tp_hash
  0,                         //tp_call
  0,                         //tp_str
  0,                         //tp_getattro
  0,                         //tp_setattro
  0,                         //tp_as_buffer
  Py_TPFLAGS_DEFAULT,        //tp_flags
  "",                        //tp_doc //TODO type docstring
  0,                         //tp_traverse
  0,                         //tp_clear
  0,                         //tp_richcompare
  0,                         //tp_weaklistoffset
  0,                         //tp_iter
  0,                         //tp_iternext
  {{.Name}}_methods,         //tp_methods
  {{.Name}}_members,         //tp_members
  0,                         //tp_getset
  0,                         //tp_base
  0,                         //tp_dict
  0,                         //tp_descr_get
  0,                         //tp_descr_set
  0,                         //tp_dictoffset
  (initproc){{.Name}}_init,  //tp_init
  0,                         //tp_alloc
  {{.Name}}_new,             //tp_new
};
