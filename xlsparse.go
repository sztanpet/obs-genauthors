package main

import (
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/tealeg/xlsx"
)

var userRex = regexp.MustCompile(`^(.+) \((.+)\)$`)

func parseTranslatorXls(r io.ReaderAt, l int64) ([]translation, error) {
	xls, err := xlsx.OpenReaderAt(r, l)
	fatalErr(err, "Could not open xls")

	if len(xls.Sheets) == 0 {
		return nil, errors.New("Invalid XLSX, no sheets found")
	}

	// map[language]deduplicatedusers
	m := map[string]map[string]int{}
	var currentLang string
	for rk, row := range xls.Sheets[0].Rows {
		// data starts from row 6 (7 if 1based)
		if rk <= 6 || len(row.Cells) < 4 {
			continue
		}

		// if the user cell is empty either its an empty row or we are at the end
		user, _ := row.Cells[2].String()
		commits, _ := row.Cells[3].Int()
		user = strings.TrimSpace(user)
		if user == "" || user == "REMOVED_USER" || user == "no data available" {
			continue
		}

		// if there is a language and is different from the current language, use it
		{
			lang, _ := row.Cells[0].String()
			lang = strings.TrimSpace(lang)
			if lang != "" && lang != currentLang {
				currentLang = lang
			}
		}

		// initialize the map at the given language if it does not exist
		if _, ok := m[currentLang]; !ok {
			m[currentLang] = map[string]int{}
		}

		commits += m[currentLang][user]
		m[currentLang][user] = commits
	}

	// now post process the users for their nicks
	// also order everything nicely
	ret := make([]translation, 0, len(m))
	for lang, users := range m {
		if len(users) == 0 {
			continue
		}

		crs := make([]contributor, 0, len(users))
		for user, commits := range users {
			var nick string
			matches := userRex.FindStringSubmatch(user)

			if len(matches) != 0 {
				user = matches[1]
				nick = matches[2]
			}

			crs = append(crs, contributor{
				Name:    user,
				Nick:    nick,
				Commits: commits,
			})
		}

		ret = append(ret, translation{
			Language:    lang,
			Translators: crs,
		})
	}

	orderTranslations(ret)
	return ret, nil
}
