package builder

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type Config struct {
	AggregationType AggregationType
	SortCols        []SortCol
	AllColumnNames  []string
	OutputFilePath  string
	SourceLink      LinkType
}

type AggregationType int

const (
	AggregationTypeList AggregationType = iota
	AggregationTypeCount
)

type SortCol struct {
	Name          string
	Type          SortColType
	SortAscending bool
}

type SortColType string

const (
	SortColTypeNum  = "number"
	SortColTypeDate = "date"
)

type LinkType int

const (
	NoLink LinkType = iota
	WikiLink
	MarkdownLink
)

type Entry struct {
	SourceFilePath string
	Fields         map[string]string
}

type Builder interface {
	Add(Entry)
	Done()
}

// BuilderWrapper implements logic applicable across all builders
type BuilderWrapper struct {
	c Config
	Builder
}

func (w BuilderWrapper) Add(e Entry) {
	e = addLinkToEntry(e, w.c.SourceLink)
	w.Builder.Add(e)
}

func NewBuilder(c Config) Builder {
	c = addLinkHeaderToConfig(c)

	var builder Builder
	if c.AggregationType == AggregationTypeList {
		builder = newListBuilder(c)
	}
	return BuilderWrapper{
		c:       c,
		Builder: builder,
	}
}

func addLinkToEntry(e Entry, lt LinkType) Entry {
	if lt != NoLink {
		e.Fields["Link"] = createLink("see file", e.SourceFilePath, lt)
	}
	return e
}
func createLink(text, path string, lt LinkType) string {
	if lt == WikiLink {
		filename := filepath.Base(path)
		return fmt.Sprintf("[[%s]]", filename)
	}
	if lt == MarkdownLink {
		return fmt.Sprintf("[%s](%s)", text, path)
	}
	return path
}
func addLinkHeaderToConfig(c Config) Config {
	if c.SourceLink != NoLink {
		c.AllColumnNames = append(c.AllColumnNames, "Link")
	}
	return c
}

type Row struct {
	fields map[string]string
}

func SortRows(rows []Row, sortOn []SortCol) []Row {
	sort.Slice(rows, func(i, j int) bool {
		a := rows[i]
		b := rows[j]

		for _, sortBy := range sortOn {
			aFieldStr := a.fields[sortBy.Name]
			bFieldStr := b.fields[sortBy.Name]

			res := isABigger(aFieldStr, bFieldStr, sortBy.Type)
			if res == same {
				// try next field
				continue
			}

			if res == bigger && sortBy.SortAscending {
				return false
			}
			if res == smaller && !sortBy.SortAscending {
				return false
			}

			// a is bigger
			return true

		}
		return false
	})
	return rows
}

type comparisonOutcome int

const (
	bigger comparisonOutcome = iota
	same
	smaller
)

func isABigger(af, bf string, fieldType SortColType) comparisonOutcome {
	if fieldType == SortColTypeDate {
		return isASooner(af, bf)
	}

	if fieldType == SortColTypeNum {
		return isANumBigger(af, bf)
	}

	return same
}

func isASooner(aDate, bDate string) comparisonOutcome {
	at, err := time.Parse("2006-01-02", aDate)
	if err != nil {
		fmt.Printf("%s is not a date \n", aDate)
		return smaller
	}
	bt, err := time.Parse("2006-01-02", bDate)
	if err != nil {
		fmt.Printf("%s is not a date \n", bDate)
		return bigger
	}
	diff := at.Sub(bt).Milliseconds()
	if diff > 0 {
		return bigger
	}
	if diff == 0 {
		return same
	}
	return smaller
}

func isANumBigger(aNum, bNum string) comparisonOutcome {
	aN, err := strconv.ParseFloat(aNum, 64)
	if err != nil {
		fmt.Printf("%s is not a number \n", aNum)
		return -1
	}
	bN, err := strconv.ParseFloat(bNum, 64)
	if err != nil {
		fmt.Printf("%s is not a number \n", bNum)
		return 1
	}

	if aN > bN {
		return bigger
	}
	if aN == bN {
		return same
	}
	return smaller
}
