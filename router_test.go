package router_test

import (
	"github.com/secondarykey/router"

	"fmt"
	"net/http"

	"net/http/httptest"
	"testing"
)

//Router
func TestTemplate(t *testing.T) {

	r := router.New(router.Template, nil)

	r.TemplateDirectory = "test"
	r.Add("/", nil)

	req := httptest.NewRequest("GET", "/hello", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	//parameter direct
}

func TestJSON(t *testing.T) {
}

func TestLogin(t *testing.T) {
}

func TestRouterURLWithDirect(t *testing.T) {

	// Parameter values test

	r := router.New(router.Direct, nil)
	r.TemplateDirectory = "test"

	r.Add("/", testIndex)
	r.Add("/test/", test)
	r.Add("/test/{param1}", testParam)
	r.Add("/test/add", testAdd)
	r.Add("/test/{param2}/{param3}", testParam2)

	req := httptest.NewRequest("GET", "/hello", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusNotFound != rec.Code {
		t.Errorf("/hello is not found")
	}

	req = httptest.NewRequest("GET", "/", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusOK != rec.Code {
		t.Errorf("/ is found")
	}

	if rec.Body.String() != "Index" {
		t.Errorf("/ write Index")
	}

	req = httptest.NewRequest("GET", "/test/", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusOK != rec.Code {
		t.Errorf("/test/ is found")
	}

	if rec.Body.String() != "Test Directory" {
		t.Errorf("/test/ write Index")
	}

	req = httptest.NewRequest("GET", "/test/aaa", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusOK != rec.Code {
		t.Errorf("/test/aaa is found")
	}

	if rec.Body.String() != "Test Param[aaa]" {
		t.Errorf("/test/aaa write want parameter[aaa] is [%s]", rec.Body.String())
	}

	req = httptest.NewRequest("GET", "/test/bbb", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusOK != rec.Code {
		t.Errorf("/test/bbb is found")
	}

	if rec.Body.String() != "Test Param[bbb]" {
		t.Errorf("/test/bbb write parameter[bbb]")
	}

	req = httptest.NewRequest("GET", "/test/add", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusOK != rec.Code {
		t.Errorf("/test/add is found")
	}

	if rec.Body.String() != "Test add method" {
		t.Errorf("/test/add not call parameter method [%s]", rec.Body.String())
	}

	req = httptest.NewRequest("GET", "/test/aaa/bbb", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if http.StatusOK != rec.Code {
		t.Errorf("/test/aaa/bbb is found")
	}

	if rec.Body.String() != "Test Param[aaa][bbb]" {
		t.Errorf("/test/aaa/bbb double parameter")
	}
}

func testIndex(p *router.Parameter) error {
	fmt.Fprintf(p.Res, "Index")
	return nil
}

func test(p *router.Parameter) error {
	fmt.Fprintf(p.Res, "Test Directory")
	return nil
}

func testParam(p *router.Parameter) error {
	fmt.Fprintf(p.Res, "Test Param[%s]", p.Get("param1"))
	return nil
}

func testAdd(p *router.Parameter) error {
	fmt.Fprintf(p.Res, "Test add method")
	return nil
}

func testParam2(p *router.Parameter) error {

	if p.Get("param1") != "" {
		return fmt.Errorf("Error")
	}

	fmt.Fprintf(p.Res, "Test Param[%s][%s]", p.Get("param2"), p.Get("param3"))

	return nil
}
