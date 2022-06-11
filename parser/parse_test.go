package parser_test

import (
	"strings"
	"testing"

	"github.com/samxsmith/snag/parser"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	t.Run("match all lines in a file", func(t *testing.T) {
		lines := []string{
			"one line",
			"#book, The Count of Monte Cristo, 2022-06-10, 5",
			"another line                                   ",
			"#book, Walden, 2022-05-10, 4.5                 ",
			"meditations, 5/10                              ",
		}

		b := strings.Join(lines, "\n")

		r := strings.NewReader(b)
		tags := []parser.Tag{
			{
				ID:     1,
				Prefix: "#book",
			},
			{
				ID:     2,
				Prefix: "#book",
			},
			{
				ID:     3,
				Prefix: "meditations",
			},
		}

		matches, err := parser.FindLinesWithTags(tags, r)
		assert.Nil(t, err)

		assert.Equal(t, len(matches), 5)

		assert.Equal(t, matches[0].TagID, 1)
		assert.Equal(t, matches[0].Line, lines[1])

		assert.Equal(t, matches[1].TagID, 2)
		assert.Equal(t, matches[1].Line, lines[1])

		assert.Equal(t, matches[2].TagID, 1)
		assert.Equal(t, matches[2].Line, lines[3])

		assert.Equal(t, matches[3].TagID, 2)
		assert.Equal(t, matches[3].Line, lines[3])

		assert.Equal(t, matches[4].TagID, 3)
		assert.Equal(t, matches[4].Line, lines[4])
	})
}

func TestParseEntry(t *testing.T) {
	t.Run("should parse the columns", func(t *testing.T) {
		line := "memory swimming in the sea, 2022-01-01"
		s := parser.EntrySchema{
			Cols: []string{
				"description", "date",
			},
			Prefix: "memory",
		}

		entry, err := parser.ParseEntry(line, s)
		assert.Nil(t, err)
		assert.Equal(t, len(entry.Fields), 2)
	})
}
