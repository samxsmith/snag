package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Tag struct {
	ID     int
	Prefix string
}

type Match struct {
	TagID int
	Line  string
}

func FindLinesWithTags(tags []Tag, r io.Reader) ([]Match, error) {
	var m []Match
	b := bufio.NewReader(r)

	for {

		line, err := b.ReadString('\n')
		if err != nil && err != io.EOF {
			return m, fmt.Errorf("error on read: %w", err)
		}

		for _, t := range tags {
			if strings.HasPrefix(line, t.Prefix) {
				l := strings.Trim(line, "\n")
				match := Match{
					TagID: t.ID,
					Line:  l,
				}
				m = append(m, match)
			}
		}
		if err != nil && err == io.EOF {
			break
		}
	}

	return m, nil
}

type EntrySchema struct {
	Cols   []string
	Prefix string
}

type Entry struct {
	Fields map[string]string
}

func ParseEntry(line string, s EntrySchema) (Entry, error) {
	e := Entry{
		Fields: map[string]string{},
	}

	line = strings.TrimPrefix(line, s.Prefix)

	fields := strings.Split(line, ",")

	if len(fields) != len(s.Cols) {
		return e, fmt.Errorf("incorrect column count")
	}

	for i, colName := range s.Cols {
		e.Fields[colName] = strings.Trim(fields[i], " ")
	}

	return e, nil
}
