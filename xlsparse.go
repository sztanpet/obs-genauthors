package main

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/tealeg/xlsx"
)

func (a *app) parseTranslatorXls(r *bytes.Reader) ([]translation, error) {
	xls, err := xlsx.OpenReaderAt(r, int64(r.Len()))
	fatalErr(err, "Could not open xls")

	if len(xls.Sheets) == 0 {
		return nil, errors.New("Invalid XLSX, no sheets found")
	}

	for rk, row := range xls.Sheets[0].Rows {
		for ck, cell := range row.Cells {
			v, err := cell.String()
			fmt.Printf("rk %v| ck: %v | v: %v| err: %v\n", rk, ck, v, err)
		}
	}
	return nil, nil
}
