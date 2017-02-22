package main

import (
	"bytes"
	htemplate "html/template"
	"io"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/sztanpet/obs-genauthors/data"
)

func (a *app) setupHTMLTemplates() {
	f, err := data.Assets.Open("tpl/layout.html")
	fatalErr(err, "Could not open layout.html")

	l, err := htemplate.New("layout.html").Parse(readFile(f))
	fatalErr(err, "Could not parse layout.html")

	f, err = data.Assets.Open("tpl/index.html")
	fatalErr(err, "Could not open index.html")

	tpl, err := htemplate.Must(l.Clone()).Parse(readFile(f))
	fatalErr(err, "Could not parse index.html")

	a.htmlTpl = tpl
}

func (a *app) setupTextTemplates() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	s := authorTpl()
	tpl, err := template.New("authors").Parse(s)
	if err != nil {
		return err
	}

	a.textTpl = tpl
	return nil
}

func (a *app) generateOutput(cs []contributor, ts []translation) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	b := &bytes.Buffer{}
	err := a.textTpl.ExecuteTemplate(b, "authors", struct {
		Contributors []contributor
		Translations []translation
	}{
		Contributors: cs,
		Translations: ts,
	})

	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func overridePath() string {
	p, err := os.Executable()
	fatalErr(err, "unable to get executable path")

	return path.Dir(p) + "/authors.tpl"
}

func saveAuthorTpl(data []byte) {
	err := ioutil.WriteFile(overridePath(), data, 0644)
	fatalErr(err, "Could not write override authors.tpl")
}

func delAuthorTpl() {
	_ = os.Remove(overridePath())
}

func readFile(r io.Reader) string {
	b := &bytes.Buffer{}
	_, err := io.Copy(b, r)
	fatalErr(err, "failed to copy authors.tpl")

	return b.String()
}

func authorTpl() string {
	p := overridePath()

	// is there an override?
	i, err := os.Stat(p)
	if err == nil && i.Size() != 0 {
		o, err := os.Open(p)
		fatalErr(err, "could not open override authors.tpl")

		defer o.Close()
		return readFile(o)
	}

	f, err := data.Assets.Open("tpl/authors.tpl")
	fatalErr(err, "unable to open embedded authors.tpl")
	defer f.Close()

	return readFile(f)
}
