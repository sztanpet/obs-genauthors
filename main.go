/***
  This file is part of obs-genauthors.

  Copyright (c) 2015 Peter Sztan <sztanpet@gmail.com>

  obs-genauthors is free software; you can redistribute it and/or modify it
  under the terms of the GNU Lesser General Public License as published by
  the Free Software Foundation; either version 3 of the License, or
  (at your option) any later version.

  obs-genauthors is distributed in the hope that it will be useful, but
  WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
  Lesser General Public License for more details.

  You should have received a copy of the GNU Lesser General Public License
  along with obs-genauthors; If not, see <http://www.gnu.org/licenses/>.
***/

package main

import (
	"context"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	ttemplate "text/template"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/codemodus/chain"
	"github.com/gorilla/csrf"
	"github.com/sztanpet/config"
	"github.com/sztanpet/obs-genauthors/data"
)

type app struct {
	htmlTpl *template.Template
	config  Config
	routes  struct {
		indexGet, indexPost http.Handler
		assets              http.Handler
	}

	once     sync.Once
	shutdown chan struct{}

	mu      sync.Mutex
	textTpl *ttemplate.Template
}

type contributor struct {
	Name, Nick, Email string
	Commits           int
}

type translation struct {
	Language    string
	Translators []contributor
}

func main() {
	a := &app{
		shutdown: make(chan struct{}),
	}

	a.setupConfig()

	err := a.setupTextTemplates()
	fatalErr(err, "Could not parse authors.tpl")

	a.setupHTMLTemplates()
	a.setupAndRunHTTP()
}

func (a *app) setupConfig() {
	a.config = Config{}
	err := config.Init(&a.config, sampleconf, "config.ini")
	fatalErr(err, "Could not parse config.ini")
}

func (a *app) setupAndRunHTTP() {
	s := http.Server{
		Addr:         a.config.Addr,
		Handler:      a,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go (func() {
		<-a.shutdown
		s.Shutdown(context.Background())
	})()

	def := chain.New(
		a.handleExtraHeaders,
		gziphandler.GzipHandler,
	)

	index := def.Append(
		csrf.Protect(
			[]byte(a.config.Secret),
			csrf.FieldName("csrf"),
			csrf.CookieName("csrf"),
			csrf.Path("/"),
			csrf.Secure(false),
		),
	)
	a.routes.indexGet = index.EndFn(a.indexGet)
	a.routes.indexPost = index.EndFn(a.indexPost)

	a.routes.assets = http.StripPrefix("/assets/", http.FileServer(data.Assets))
	a.routes.assets = gziphandler.GzipHandler(a.routes.assets)

	log.Printf("Starting server: %s\n", a.config.Addr)
	panic(s.ListenAndServe())
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isGet = r.Method == http.MethodGet
	if !isGet && r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	p := r.URL.Path

	if p == "/" {
		if isGet {
			a.routes.indexGet.ServeHTTP(w, r)
		} else {
			a.routes.indexPost.ServeHTTP(w, r)
		}
		return
	}

	if p == "/exit" {
		_, _ = w.Write([]byte("done"))
		a.once.Do(func() { close(a.shutdown) })
		return
	}

	// handle assets
	if strings.HasPrefix(p, "/assets") {
		a.routes.assets.ServeHTTP(w, r)
		return
	}

	_, _ = io.Copy(ioutil.Discard, r.Body)
	http.Error(w, "404 - Not Found", http.StatusNotFound)
}

func (a *app) handleExtraHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Frame-Options", "deny")
		h.Set("X-XSS-Protection", "1; mode=block")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Content-Security-Policy", `default-src 'self';`+
			` script-src 'self';`+
			` style-src 'self' 'unsafe-inline';`+
			` connect-src 'self';`+
			` frame-src 'self';`+
			` block-all-mixed-content;`,
		)

		h.Set("Content-Type", "text/html; charset=utf-8")

		next.ServeHTTP(w, r)
	})
}

func fatalErr(err error, str string) {
	if err != nil {
		panic(str + ": " + err.Error())
	}
}
