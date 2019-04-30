package router

import (
	"strings"

	"golang.org/x/xerrors"
)

type Glob []*part

// *
// ?
// [a-z]
// [!a-z]
// []+
func Compile(g string) (Glob, error) {
	chars := "*?["

	if !strings.ContainsAny(g, chars) {
		globs := make([]*part, 1)
		globs[0] = newPart(None, g, false)
		return globs, nil
	}
	return split(g)
}

func split(g string) (Glob, error) {

	slc := make([]*part, 0, 8)
	wk := g

	//parse []
	for {
		if wk == "" {
			break
		}
		idx := strings.Index(wk, "[")
		if idx == -1 {

			p := newPart(None, wk, false)
			slc = append(slc, p)

			break
		} else {
			ldx := strings.Index(wk[idx:], "]")
			if idx == -1 {
				return nil, xerrors.Errorf(`] not found`)
			}

			p := newPart(None, wk[:idx], false)
			slc = append(slc, p)

			pt := 1
			pattern := Charactor
			if wk[idx+ldx:idx+ldx+1] == "+" {
				pt = 2
				pattern = Strings
			}

			p = newPart(pattern, wk[idx+1:idx+ldx], true)

			slc = append(slc, p)
			wk = wk[idx+ldx+pt:]
		}
	}

	//parse *,?
	rtn := make([]*part, 0, 8)
	for _, elm := range slc {
		if elm.pattern == None {
			wk := parse(elm.source, "*")
			rtn = append(rtn, wk...)
		} else {
			rtn = append(rtn, elm)
		}
	}

	return rtn, nil
}

func parse(src string, sp string) []*part {

	rtn := make([]*part, 0, 4)
	wk := strings.Split(src, sp)

	p := Strings
	if sp == "?" {
		p = Charactor
	}

	for i, buf := range wk {

		if i == 0 && buf == "" {
			p := newPart(p, sp, false)
			rtn = append(rtn, p)
		}

		if buf != "" {
			if sp != "?" {
				qslc := parse(buf, "?")
				rtn = append(rtn, qslc...)
			} else {
				rtn = append(rtn, newPart(p, buf, false))
			}
		}

		if i != len(wk)-1 && buf != "" {
			p := newPart(p, sp, false)
			rtn = append(rtn, p)
		}
	}
	return rtn
}

func (g Glob) Clear() {
	for _, glb := range g {
		glb.done = false
	}
}

func (g Glob) Match(s string) bool {

	g.Clear()

	wk := s
	ok := true

	for _, glb := range g {
		wk, ok = glb.match(wk)
		if !ok {
			return false
		}
	}

	return g.judge()
}

func (g Glob) judge() bool {
	for _, glb := range g {
		if !glb.done {
			return false
		}
	}
	return true
}

func (g Glob) String() string {
	wk := ""
	for _, elm := range g {
		buf := elm.String()

		if elm.limit {
			buf = "[" + buf + "]"
			if elm.pattern == Strings {
				buf += "+"
			}
		}
		wk += buf
	}
	return wk
}

type globPattern int

const (
	None globPattern = iota
	Strings
	Charactor
)

type part struct {
	pattern globPattern
	source  string
	limit   bool
	done    bool
}

func newPart(p globPattern, s string, limit bool) *part {
	part := part{
		pattern: p,
		source:  s,
		limit:   limit,
		done:    false,
	}
	return &part
}

func (p *part) match(t string) (string, bool) {

	buf := t

	switch p.pattern {
	case None:
	case Strings:
	case Charactor:
		if p.source != t[0:1] {
			return "", false
		}
		buf = t[1:]
	}
	return buf, true
}

func (p part) String() string {
	return p.source
}
