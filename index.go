package main

import (
	"net/http"
	"strconv"

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
		AuthorTpl: authorTpl(),
	})

	fatalErr(err, "Could not execute index template")
}

func (a *app) indexPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10e6) // 10mb of memory
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var errStr string
	var cs []contributor
	var ts []translation

	if s := r.PostFormValue("gitauthors"); s != "" {
		cs = unserializeContributors(s)
	} else {
		cs = a.gitAuthors()
	}

	if s := r.PostFormValue("translatorauthors"); s != "" {
		ts = unserializeTranslators(s)
	} else if f, fh, err := r.FormFile("file"); err == nil {
		ss := fh.Header.Get("Content-Length")
		var size int64
		if ss == "" {
			size, err = f.Seek(0, 2)
			fatalErr(err, "Could not seek to the end of the form file")
		} else {
			size, err = strconv.ParseInt(ss, 10, 64)
			fatalErr(err, "Could not parse the length of the form file")
		}

		ts, err = parseTranslatorXls(f, size)
		if err != nil {
			errStr = "Invalid file: " + err.Error()
		}
	}

	if s := r.PostFormValue("authortpl"); s != "" {
		saveAuthorTpl([]byte(s))
		err = a.setupTextTemplates()
		if err != nil {
			errStr = "Invalid template: " + err.Error()
		}
	}

	out, err := a.generateOutput(cs, ts)
	if err != nil {
		errStr = "Invalid template: " + err.Error()
	}

	err = a.htmlTpl.ExecuteTemplate(w, "layout.html", struct {
		Conf              Config
		CSRFToken         string
		AuthorTpl         string
		GitAuthors        string
		TranslatorAuthors string
		Output            string
		Error             string
	}{
		Conf:              a.config,
		CSRFToken:         csrf.Token(r),
		AuthorTpl:         authorTpl(),
		GitAuthors:        serializeContributors(cs),
		TranslatorAuthors: serializeTranslators(ts),
		Output:            out,
		Error:             errStr,
	})

	fatalErr(err, "Could not execute index template")
}
