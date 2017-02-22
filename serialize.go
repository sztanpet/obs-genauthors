package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
)

func orderTranslations(ts []translation) {
	for _, v := range ts {
		orderContributors(v.Translators)
	}
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Language < ts[j].Language
	})
}
func orderContributors(cs []contributor) {
	sort.Slice(cs, func(i, j int) bool {
		if cs[i].Commits != cs[j].Commits {
			return cs[i].Commits > cs[j].Commits
		}

		is := strings.ToLower(cs[i].Name)
		js := strings.ToLower(cs[j].Name)
		return is < js
	})
}

func unserializeContributors(s string) []contributor {
	ret := []contributor{}
	r := csv.NewReader(bytes.NewBufferString(s))
	for {
		r, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fatalErr(err, "Could not read contributors string as csv")
		}

		if len(r) != 3 {
			log.Printf("Invalid line in csv, pased: %v\n", r)
			continue
		}

		commits, _ := strconv.Atoi(r[2])

		ret = append(ret, contributor{
			Name:    strings.TrimSpace(r[0]),
			Email:   strings.TrimSpace(r[1]),
			Commits: commits,
		})
	}

	orderContributors(ret)
	return ret
}

func unserializeTranslators(s string) []translation {
	ret := []translation{}
	l := map[string][]contributor{}

	r := csv.NewReader(bytes.NewBufferString(s))
	for {
		r, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fatalErr(err, "Could not read contributors string as csv")
		}

		if len(r) != 4 {
			log.Printf("Invalid line in csv, pased: %v\n", r)
			continue
		}
		lang := strings.TrimSpace(r[0])
		commits, _ := strconv.Atoi(r[3])

		cs := l[lang]
		cs = append(cs, contributor{
			Name:    strings.TrimSpace(r[1]),
			Nick:    strings.TrimSpace(r[2]),
			Commits: commits,
		})

		l[lang] = cs
	}

	for k, v := range l {
		ret = append(ret, translation{
			Language:    k,
			Translators: v,
		})
	}

	orderTranslations(ret)
	return ret
}

func serializeContributors(cs []contributor) string {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	for _, c := range cs {
		if err := w.Write([]string{c.Name, c.Email, strconv.Itoa(c.Commits)}); err != nil {
			fatalErr(err, "Could not write csv of contributors")
		}
	}
	w.Flush()

	return b.String()
}

func serializeTranslators(ts []translation) string {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	for _, t := range ts {
		for _, c := range t.Translators {
			if err := w.Write([]string{t.Language, c.Name, c.Nick, strconv.Itoa(c.Commits)}); err != nil {
				fatalErr(err, "Could not write csv of contributors")
			}
		}
	}
	w.Flush()

	return b.String()
}
