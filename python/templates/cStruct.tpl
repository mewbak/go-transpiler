// create in go
extern long long create{{.Name}}(void);

//free go pointer
extern void free{{.Name}}(long long elem);

// no members are defined as we use getters and setters
static PyMemberDef {{.Name}}_members[] = {
    {NULL}
};

{{template "cGetSet.tpl" .}}

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

    self = ({{.Name}}*)type->tp_alloc(type, 0);
    if (self != NULL) {
      self->go{{.Name}} = create{{.Name}}();
    }

    return (PyObject *)self;
}

static int
{{.Name}}_init({{.Name}} *self, PyObject *args, PyObject *kwargs)
{
{{if (len .NamedMembers)}}

    {{- range $i, $_ := .NamedMembers}}
    PyObject* arg{{.Name}} = NULL;
    {{- end}}

    static char *kwlist[] = {
        {{- range $i, $_ := .NamedMembers}}
        "{{camelToSnake .Name}}",
        {{- end}}
        NULL
    };

    if (!PyArg_ParseTupleAndKeywords(
        args, kwargs, "|{{range .NamedMembers}}O{{end}}", kwlist,
        {{- range $i, $_ := .NamedMembers}}
        &arg{{.Name}}{{if notLast $i $.NamedMembers}},{{end}}
        {{- end}})) {
        return -1;
    }
    {{range .NamedMembers}}
    if(NULL != arg{{.Name}} &&
       0 != {{$.Name}}Set_{{.Name}}(self, arg{{.Name}}, NULL)){
        return -1;
    }
    {{- end}}
{{- end}}

    return 0;
}

{{template "cStructFuncs.tpl" .}}

PyTypeObject {{.Name}}_type = {
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
  {{.Name}}_getseters,       //tp_getset
  0,                         //tp_base
  0,                         //tp_dict
  0,                         //tp_descr_get
  0,                         //tp_descr_set
  0,                         //tp_dictoffset
  (initproc){{.Name}}_init,  //tp_init
  0,                         //tp_alloc
  {{.Name}}_new,             //tp_new
};
