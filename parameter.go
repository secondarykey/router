package router

import (
	"net/http"
)

//
// Parameter is Web Request,Response data
//
type Parameter struct {
	Res http.ResponseWriter
	Req *http.Request

	values    map[string]string
	Output    map[string]interface{}
	Templates []string

	//Headers
	Direct      bool
	ContentType string
}

// create parameter
func newParameter(w http.ResponseWriter, r *http.Request) *Parameter {
	p := Parameter{}
	p.values = nil

	p.Req = r
	p.Res = w
	p.Output = nil
	p.Templates = nil

	p.Direct = false
	p.ContentType = "text/html"
	return &p
}

func (p *Parameter) IsPOST() bool {
	return p.Req.Method == http.MethodPost
}

func (p *Parameter) IsGET() bool {
	return p.Req.Method == http.MethodGet
}

// Redirect
func (p *Parameter) Redirect(path string, status int) error {
	p.Direct = true
	http.Redirect(p.Res, p.Req, path, status)
	return nil
}

// Output values set
func (p *Parameter) Set(key string, v interface{}) {
	if p.Output == nil {
		p.Output = make(map[string]interface{})
	}
	p.Output[key] = v
}

// query value get
func (p *Parameter) Query(key string) []string {
	q := p.Req.URL.Query()
	return q[key]
}

// values get
func (p *Parameter) Get(key string) string {
	if p.values == nil {
		return ""
	}
	return p.values[key]
}
