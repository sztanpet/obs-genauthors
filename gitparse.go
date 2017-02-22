package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
)

func (a *app) gitAuthors() ([]contributor, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	log.Printf("GIT binary found at %q", gitPath)
	log.Printf("Descending into the repo path: %q", a.config.RepoPath)
	err = os.Chdir(a.config.RepoPath)
	if err != nil {
		return nil, err
	}

	out, err := exec.Command(
		gitPath,
		"log",
		"--all",
		"--no-merges",
		"--use-mailmap",
		"--format=%aN|%cE",
	).Output()
	if err != nil {
		return nil, err
	}

	c := map[string]struct {
		email   string
		commits int
	}{}

	ls := bytes.Split(out, []byte("\n"))
	for _, l := range ls {
		l = bytes.TrimSpace(l)
		prts := bytes.Split(l, []byte("|"))
		if len(prts) == 1 {
			continue
		} else if len(prts) != 2 {
			return nil, errors.New("Unexpected separator found in git author line")
		}

		name := string(prts[0])
		email := string(prts[1])
		data := c[name]
		data.email = email
		data.commits++

		c[name] = data
	}

	ret := make([]contributor, 0, len(c))
	for name, data := range c {
		ret = append(ret, contributor{
			Name:    name,
			Email:   data.email,
			Commits: data.commits,
		})
	}

	orderContributors(ret)
	return ret, nil
}
