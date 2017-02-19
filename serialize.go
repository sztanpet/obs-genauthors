package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"sort"
	"strings"
)

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

		if len(r) != 2 {
			log.Printf("Invalid line in csv, pased: %v\n", r)
			continue
		}

		ret = append(ret, contributor{
			Name:  strings.TrimSpace(r[0]),
			Email: strings.TrimSpace(r[1]),
		})
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name < ret[j].Name
	})

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

		if len(r) != 3 {
			log.Printf("Invalid line in csv, pased: %v\n", r)
			continue
		}
		lang := strings.TrimSpace(r[0])

		cs := l[lang]
		cs = append(cs, contributor{
			Name: strings.TrimSpace(r[1]),
			Nick: strings.TrimSpace(r[2]),
		})

		l[lang] = cs
	}

	for k, v := range l {
		sort.Slice(v, func(i, j int) bool {
			return v[i].Name < v[j].Name
		})

		ret = append(ret, translation{
			Language:    k,
			Translators: v,
		})
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Language < ret[j].Language
	})

	return ret
}

func serializeContributors(cs []contributor) string {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	for _, c := range cs {
		if err := w.Write([]string{c.Name, c.Email}); err != nil {
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
			if err := w.Write([]string{t.Language, c.Name, c.Nick}); err != nil {
				fatalErr(err, "Could not write csv of contributors")
			}
		}
	}
	w.Flush()

	return b.String()
}
