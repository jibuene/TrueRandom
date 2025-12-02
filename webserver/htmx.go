package webserver

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseGlob("webserver/templates/*.html")
	if err != nil {
		panic(err)
	}
}

func mainFrontEnd(w http.ResponseWriter, req *http.Request) {
	tmpl.ExecuteTemplate(w, "main.html", nil)
}

func randomNumberFrontend(w http.ResponseWriter, req *http.Request) {
	realRandom.Generate()
	num := realRandom.Number.String()
	messages := realRandom.Basis
	tmpl.ExecuteTemplate(w, "random", map[string]string{
		"Number":   num,
		"Messages": messages,
	})
}
