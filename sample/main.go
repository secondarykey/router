package main

import (
	"github.com/secondarykey/router"

	"net/http"
)

func main() {

	r1 := router.New(router.Template, nil)
	r1.Add("/", IndexHandler)
	r1.Add("/login", LoginHandler)
	r1.TemplateDirectory = "templates"
	http.Handle("/", r1)

	r2 := router.New(router.Direct, nil)
	r2.Add("/direct/{param}", DirectHandler)
	http.Handle("/direct/", r2)

	r3 := router.New(router.JSON, nil)
	r3.Add("/api/{param}", JSONHandler)
	http.Handle("/api/", r3)

	r4 := router.New(router.Template, loggedIn)
	r4.TemplateDirectory = "templates/dashboard"
	r4.Add("/dashboard/", DashboardHandler)
	http.Handle("/dashboard/", r4)

	http.ListenAndServe(":8080", nil)
}

func IndexHandler(p *router.Parameter) error {
	p.Set("Data", "Data")
	p.Set("QUERY", p.Query("test"))
	p.Templates = []string{"index.tmpl"}
	return nil
}

func DirectHandler(p *router.Parameter) error {
	_, err := p.Res.Write([]byte("Direct:" + p.Get("param")))
	return err
}

func JSONHandler(p *router.Parameter) error {
	v := p.Get("param")
	p.Set("Data", v)
	return nil
}

func DashboardHandler(p *router.Parameter) error {
	p.Templates = []string{"index.tmpl"}
	return nil
}

func LoginHandler(p *router.Parameter) error {

	//POST
	if p.IsPOST() {
		p.Req.ParseForm()
		u := p.Req.FormValue("userid")
		pass := p.Req.FormValue("password")

		if u == "user" && pass == "p@ssword" {
			//Set Cookie
			sc := http.Cookie{
				Name:   "Login",
				Value:  "true",
				MaxAge: 60 * 60,
				Path:   "/",
			}
			http.SetCookie(p.Res, &sc)
			return p.Redirect("/dashboard/", 302)
		}
	}
	p.Templates = []string{"login.tmpl"}
	return nil
}

func loggedIn(p *router.Parameter) bool {

	//Get Cookie
	sc, err := p.Req.Cookie("Login")
	if err == nil && sc != nil {
		if sc.Value == "true" {
			return true
		}
	}

	//redirect
	p.Redirect("/login", 302)
	return false
}
