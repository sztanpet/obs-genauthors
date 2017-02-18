package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/sztanpet/obs-genauthors/data"
)

func (a *app) setupTextTemplates() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	tpl, err := template.New("authors").Parse(authorTpl())
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
	err := a.textTpl.ExecuteTemplate(b, "", struct {
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
