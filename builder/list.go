package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type listBuilder struct {
	c     Config
	lines []Row
}

func newListBuilder(c Config) *listBuilder {
	return &listBuilder{
		c:     c,
		lines: []Row{},
	}
}

func (lb *listBuilder) Add(e Entry) {
	r := Row{
		fields: e.Fields,
	}
	lb.lines = append(lb.lines, r)
}
func (lb *listBuilder) Done() {
	var rows []string

	headers := make([]string, len(lb.c.AllColumnNames))
	seps := make([]string, len(lb.c.AllColumnNames))
	for i, col := range lb.c.AllColumnNames {
		headers[i] = col
		seps[i] = "----"
	}

	rows = append(rows, colsToMarkdown(headers))
	rows = append(rows, colsToMarkdown(seps))

	sortedRows := SortRows(lb.lines, lb.c.SortCols)

	for _, l := range sortedRows {
		var row []string
		for _, colName := range lb.c.AllColumnNames {
			row = append(row, l.fields[colName])
		}

		lf := colsToMarkdown(row)

		// TODO: could allow user to specify go fmt string

		rows = append(rows, lf)
	}
	s := strings.Join(rows, "\n")
	b := []byte(s)

	outputDir := filepath.Dir(lb.c.OutputFilePath)
	os.MkdirAll(outputDir, 0700)

	err := os.WriteFile(lb.c.OutputFilePath, b, 0700)
	if err != nil {
		panic(err)
	}

}

func colsToMarkdown(cols []string) string {
	lf := strings.Join(cols, " | ")
	lf = fmt.Sprintf("| %s |", lf)
	return lf
}

func toWikiLink(s string) string {
	return fmt.Sprintf("[[%s]]", s)
}
