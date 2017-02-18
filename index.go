package main

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func (a *app) indexGet(w http.ResponseWriter, r *http.Request) {
	err := a.htmlTpl.ExecuteTemplate(w, "layout.html", struct {
		Conf              Config
		CSRFToken         string
		AuthorTpl         string
		GitAuthors        string
		TranslatorAuthors string
		Output            string
	}{
		Conf:      a.config,
		CSRFToken: csrf.Token(r),
		AuthorTpl: a.authorTpl(),
	})

	fatalErr(err, "Could not execute index template")
}

func (a *app) indexPost(w http.ResponseWriter, r *http.Request) {
	// TODO: parse form, if there is a file, use that for the data
	// if there is no file, and the form fields gitauthors and translatorauthors
	// are not empty, use those for data, otherwise error
	// save the authortpl as an override (care about it being the same as the default?)
	// if no error, generate the output according to the template
	var cs []contributor
	var ts []translation

	out, err := a.generateOutput(cs, ts)
	// TODO show a nicely formatted error message, this is a simple user error
	fatalErr(err, "Could not generate output")

	err = a.htmlTpl.ExecuteTemplate(w, "layout.html", struct {
		Conf              Config
		CSRFToken         string
		AuthorTpl         string
		GitAuthors        string
		TranslatorAuthors string
		Output            string
	}{
		Conf:              a.config,
		CSRFToken:         csrf.Token(r),
		AuthorTpl:         a.authorTpl(),
		GitAuthors:        serializeContributors(cs),
		TranslatorAuthors: serializeTranslators(ts),
		Output:            out,
	})

	fatalErr(err, "Could not execute index template")
}
