package model

// original content from https://github.com/timburks/mailinglabels
// Public Domain

import (
	"github.com/jung-kurt/gofpdf"
	"io"
)

// Avery 5160
func AveryLabels(pdffile io.Writer, addresses []addressFormatter) {
	count := 0

	pdf := gofpdf.New("P", "in", "Letter", "")

	pagew, pageh, _ := pdf.PageSize(0)

	labelh := 1.0                         // 1" high
	labelw := 2.5                         // 2 1/2" wide
	marginv := (pageh - 10*labelh) / 2.0  // 1/2" from top
	marginh := (pagew - 3.0*labelw) / 4.0 // label margin
	pdf.SetFont("Helvetica", "", 10)      // "B" for bold, "I" for italic

	for _, address := range addresses {
		formatted, err := FormatAddress(address)
		if err != nil {
			continue
		}

		if count == 0 {
			pdf.AddPage()
		}
		row := count / 3
		col := count % 3
		x := marginh + float64(col)*(labelw+marginh)
		y := float64(row)*labelh + marginv
		pdf.SetXY(x, y)
		pdf.MultiCell(
			labelw,
			labelh/5,
			formatted,
			"",    // no border
			"LM",  // left justify, middle
			false) // don't fill
		count += 1
		if count == 30 {
			count = 0
		}
	}
	if err := pdf.Output(pdffile); err != nil {
		// nothing
	}
}
