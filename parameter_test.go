package router_test

import (
	"github.com/secondarykey/router"

	"net/http"

	"net/http/httptest"
	"testing"
)

func TestParameter(t *testing.T) {

	req := httptest.NewRequest("GET", "/test/aaa?q=query", nil)
	rec := httptest.NewRecorder()

	p := router.NewParameter(rec, req)
	if p == nil {
		t.Errorf("error that p is nil")
	}
	if p.Output != nil {
		t.Errorf(`error that p.Output is nil`)
	}
	if p.Templates != nil {
		t.Errorf(`error that p.Templates is nil`)
	}
	if p.Direct {
		t.Errorf(`error that p.Direct is false`)
	}
	if p.Req != req {
		t.Errorf(`error that p.Req equal req`)
	}
	if p.Res != rec {
		t.Errorf(`error that p.Res is rec`)
	}
	if p.ContentType != "text/html" {
		t.Errorf(`error that p.ContentType is text/html`)
	}

	if !p.IsGET() {
		t.Errorf(`error that p.IsGET() is true`)
	}
	if p.IsPOST() {
		t.Errorf(`error that p.IsPOST() is false`)
	}

	p.Set("test", "set test")
	if p.Output["test"] != "set test" {
		t.Errorf(`error that p.Get("test") is "set test"`)
	}

	q := p.Query("q")
	if len(q) != 1 {
		t.Errorf(`error that p.Query() length 1`)
	}

	if q[0] != "query" {
		t.Errorf(`error that p.Query("q")[0] is query`)
	}

	err := p.Redirect("/aaa", http.StatusFound)
	if err != nil {
		t.Errorf(`error that p.Rediret return is nil`)
	}
	if !p.Direct {
		t.Errorf(`error that p.Direct is true`)
	}

	if http.StatusFound != rec.Code {
		t.Errorf(`error that Redirect code equals StatusFound`)
	}

	req = httptest.NewRequest("POST", "/test/aaa", nil)
	rec = httptest.NewRecorder()
	p = router.NewParameter(rec, req)
	if p.IsGET() {
		t.Errorf(`error that p.IsGET() is false`)
	}
	if !p.IsPOST() {
		t.Errorf(`error that p.IsPOST() is true`)
	}
}
