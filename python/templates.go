package python

import (
    "path"
    "runtime"
    "text/template"
)

var cTemplate *template.Template
var goTemplate *template.Template

func init() {

    _, filename, _, _ := runtime.Caller(1)
    base := path.Dir(filename)

    cGlob := path.Join(base, "templates", "c*.tpl")

    tpl, err := template.New("c.tpl").Funcs(templateFuncs).ParseGlob(cGlob)
    if nil != err {
        panic(err.Error())
    }
    cTemplate = tpl

    goGlob := path.Join(base, "templates", "go*.tpl")

    tpl, err = template.New("go.tpl").Funcs(templateFuncs).ParseGlob(goGlob)
    if nil != err {
        panic(err.Error())
    }
    goTemplate = tpl

}
