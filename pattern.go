package router

import (
	"strings"

	"golang.org/x/xerrors"
)

// URL pattern structure
type pattern struct {
	source  string
	handler Handler
	slice   []string
	valMap  map[int]string

	methods []string
}

// create URL pattern structure
func newPattern(s string, h Handler) (*pattern, error) {

	slice := strings.Split(s, "/")
	p := pattern{
		source:  s,
		handler: h,
		slice:   slice,
		valMap:  nil,
	}

	if err := p.check(); err != nil {
		return nil, xerrors.Errorf("check error: %w", err)
	}

	return &p, nil
}

// pattern check
func (p *pattern) check() error {

	for idx, elm := range p.slice {
		if idx == 0 {
			if elm != "" {
				return xerrors.Errorf(`pattern value is first charctor "/"`)
			}
			continue
		}

		if strings.Index(elm, "{") == -1 {
			if strings.Index(elm, "}") != -1 {
				return xerrors.Errorf(`"{" not exists,"}" exists [%s]`, elm)
			}
		} else {
			if strings.Index(elm, "{") == 0 {
				if strings.Index(elm, "}") != len(elm)-1 {
					return xerrors.Errorf(`"}" last charctor[%s]`, elm)
				}

				if p.valMap == nil {
					p.valMap = make(map[int]string)
				}

				val := elm[1 : len(elm)-1]
				if val == "" {
					return xerrors.Errorf(`"{}" is bad[%s]`, elm)
				}

				p.valMap[idx] = val

			} else {
				return xerrors.Errorf(`"{" first charctor[%s]`, elm)
			}
		}
	}

	return nil
}

// url pattern match
func (p *pattern) match(url string) (map[string]string, bool) {

	if url == p.source {
		return nil, true
	}

	if p.valMap == nil {
		return nil, false
	}

	slc := strings.Split(url, "/")
	if len(p.slice) != len(slc) {
		return nil, false
	}

	flag := true
	wkMap := make(map[string]string)

	for idx, buf := range p.slice {
		if v, ok := p.valMap[idx]; ok {
			wkMap[v] = slc[idx]
		} else {
			if buf != slc[idx] {
				flag = false
				break
			}
		}
	}

	if flag {
		return wkMap, true
	}

	return nil, false
}

// allowed method judgement
func (p *pattern) allowedMethod(method string) bool {
	if p.methods == nil {
		return true
	}

	for _, elm := range p.methods {
		if elm == method {
			return true
		}
	}
	return false
}
