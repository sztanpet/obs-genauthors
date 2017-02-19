package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"sort"
)

func (a *app) gitAuthors() ([]contributor, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	err = os.Chdir(a.config.RepoPath)
	if err != nil {
		return nil, err
	}

	out, err := exec.Command(gitPath, "log", "--all", "--format=%aN|%cE").Output()
	if err != nil {
		return nil, err
	}

	c := map[string]string{}
	ls := bytes.Split(out, []byte("\n"))
	for _, l := range ls {
		l = bytes.TrimSpace(l)
		prts := bytes.Split(l, []byte("|"))
		if len(prts) == 1 {
			continue
		} else if len(prts) != 2 {
			return nil, errors.New("Unexpected separator found in git author line")
		}

		c[string(prts[0])] = string(prts[1])
	}

	ret := make([]contributor, 0, len(c))
	for name, mail := range c {
		ret = append(ret, contributor{
			Name:  name,
			Email: mail,
		})
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name < ret[j].Name
	})

	return ret, nil
}
