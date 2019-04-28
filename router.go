package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"golang.org/x/xerrors"
)

//
// Router is URL Pattern matching.
//
type Router struct {
	behavior Behavior
	login    LoginHandler
	patterns []*pattern

	TemplateDirectory string
	TemplateFunc      func(*Parameter) error
	JSONFunc          func(*Parameter) error
	ErrorTemplates    []string
	ErrorFunc         func(*Parameter, int, error)
}

// login handle type
type LoginHandler func(*Parameter) bool

// handle behavior
type Behavior int

const (
	Direct Behavior = iota
	Template
	JSON
)

// create router
func New(b Behavior, h LoginHandler) *Router {
	r := Router{}
	r.behavior = b
	r.login = h
	r.patterns = make([]*pattern, 0)

	r.TemplateFunc = r.SetTemplates
	r.JSONFunc = r.SetJSON
	r.ErrorFunc = r.SetError

	r.TemplateDirectory = "."
	r.ErrorTemplates = []string{"error.tmpl"}
	return &r
}

type Handler func(*Parameter) error

// url pattern add
func (r *Router) Add(key string, h Handler, methods ...string) error {

	p, err := newPattern(key, h)
	if err != nil {
		return xerrors.Errorf("Router.Add(%s): %w", key, err)
	}

	p.methods = methods
	r.patterns = append(r.patterns, p)
	return nil
}

// get url matching pattern
func (router *Router) getPattern(p *Parameter) (*pattern, error) {

	path := p.Req.URL.Path
	var pat *pattern

	for _, elm := range router.patterns {
		d, ok := elm.match(path)
		if ok {
			if d != nil {
				p.values = make(map[string]string)
				for key, v := range d {
					p.values[key] = v
				}
				if pat != nil {
					log.Printf("ducaple url [%s]\n", path)
				}
				pat = elm
			} else {
				return elm, nil
			}
		}
	}

	if pat != nil {
		return pat, nil
	}

	return nil, xerrors.Errorf("URL Pattern Not Found[%s]", path)
}

// implements http.Handle
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	p := newParameter(w, r)

	if router.login != nil {
		if !router.login(p) {
			return
		}
	}

	if router.behavior == Direct {
		p.Direct = true
	}

	pat, err := router.getPattern(p)
	if err != nil {
		router.ErrorFunc(p, http.StatusNotFound, err)
		return
	}

	if !pat.allowedMethod(r.Method) {
		router.ErrorFunc(p, http.StatusMethodNotAllowed, xerrors.Errorf("[%s] is not allowed %s method", r.URL.Path, r.Method))
		return
	}

	err = pat.handler(p)
	if err != nil {
		router.ErrorFunc(p, http.StatusInternalServerError, err)
		return
	}

	if !p.Direct {
		if router.behavior == Template {
			if p.Templates == nil || len(p.Templates) == 0 {
				err = xerrors.Errorf("error Behavior->Template but templates not setting")
			} else {
				p.Res.Header().Set("Content-Type", p.ContentType)
				err = router.TemplateFunc(p)
			}
		} else if router.behavior == JSON {
			p.Res.Header().Set("Content-Type", "application/json")
			err = router.JSONFunc(p)
		} else {
			err = xerrors.Errorf("error Behavior")
		}
	}

	if err != nil {
		router.ErrorFunc(p, http.StatusInternalServerError, err)
		return
	}
}

// Default Template Setter
func (r *Router) SetTemplates(p *Parameter) error {

	tmpls := make([]string, len(p.Templates))
	for idx, file := range p.Templates {
		tmpls[idx] = filepath.Join(r.TemplateDirectory, file)
	}

	tmpl := template.Must(template.ParseFiles(tmpls...))
	return tmpl.Execute(p.Res, p.Output)
}

// Default Error Template Setter
func (r *Router) SetError(p *Parameter, status int, err error) {

	buf := fmt.Sprintf("%+v", err)
	log.Printf("%s\n", buf)

	p.Res.WriteHeader(status)

	dto := struct {
		Status      int
		Title       string
		Description string
	}{status, http.StatusText(status), buf}

	p.Set("Error", dto)
	p.Templates = r.ErrorTemplates

	e := r.SetTemplates(p)
	if e != nil {
		log.Printf("%+v", err)
	}
}

// Default JSON Setter
func (r *Router) SetJSON(p *Parameter) error {
	res, err := json.Marshal(p.Output)
	if err != nil {
		return err
	}
	_, err = p.Res.Write(res)
	return err
}
