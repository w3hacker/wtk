package wtk

import (
	"html/template"
	"io/ioutil"
)

var tplFuncMap template.FuncMap
var tplVars map[string]interface{}

func init() {
	tplFuncMap = make(template.FuncMap)
	tplVars = make(map[string]interface{})
}

func AddTemplateFunc(name string, tplFunc interface{}) {
	tplFuncMap[name] = tplFunc
}

func SetTemplateVar(name string, value interface{}) {
	tplVars[name] = value
}

type Template struct {
	hdlr      *Handler
	tpl       *template.Template
	vars      map[string]interface{}
	tplResult *wtkTemplateResult
}

func (this *Template) SetVar(name string, value interface{}) {
	if this.vars == nil {
		this.vars = make(map[string]interface{})
	}
	this.vars[name] = value
}

func (this *Template) GetVar(name string) interface{} {
	if this.vars == nil {
		this.vars = make(map[string]interface{})
	}
	if val, ok := this.vars[name]; ok {
		return val
	}
	return nil
}

func (this *Template) SetTemplateString(str string) bool {
	this.tpl = template.New("")
	this.tpl.Funcs(tplFuncMap).Parse(str)
	return true
}

func (this *Template) SetTemplateFile(filename string) bool {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return false
	}
	return this.SetTemplateString(string(content))
}

func (this *Template) SetSubTemplateString(name, str string) bool {
	if this.tpl == nil {
		return false
	}
	tpl := this.tpl.New(name)
	tpl.Funcs(tplFuncMap).Parse(`{{define "` + name + `"}}` + str + `{{end}}`)
	return true
}

func (this *Template) SetSubTemplateFile(name, filename string) bool {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return false
	}
	return this.SetSubTemplateString(name, string(content))
}

func (this *Template) Parse() bool {
	if this.tpl == nil {
		return false
	}
	if this.tplResult != nil {
		return false
	}

	this.hdlr.callHandlerHook("BeforeRender")
	if this.hdlr.Context.response.Finished {
		return true
	}

	if this.vars == nil {
		this.vars = make(map[string]interface{})
	}
	for n, v := range tplVars {
		if _, ok := this.vars[n]; !ok {
			this.vars[n] = v
		}
	}
	this.tplResult = &wtkTemplateResult{data: []byte{}}
	err := this.tpl.Execute(this.tplResult, this.vars)
	if err != nil {
		return false
	}

	this.hdlr.callHandlerHook("AfterRender")
	return true
}

func (this *Template) GetResult() []byte {
	if this.tplResult == nil {
		return []byte{}
	}
	return this.tplResult.Bytes()
}

func (this *Template) GetResultString() string {
	return string(this.GetResult())
}

func (this *Template) SetResult(p []byte) {
	if this.tplResult == nil {
		return
	}
	this.tplResult.SetBytes(p)
}

func (this *Template) SetResultString(s string) {
	this.SetResult([]byte(s))
}

type wtkTemplateResult struct {
	data []byte
}

func (this *wtkTemplateResult) Write(p []byte) (n int, err error) {
	this.data = append(this.data, p...)
	return len(p), nil
}

func (this *wtkTemplateResult) SetBytes(p []byte) {
	this.data = p
}

func (this *wtkTemplateResult) String() string {
	return string(this.data)
}

func (this *wtkTemplateResult) Bytes() []byte {
	return this.data
}
