package snag

import (
	"github.com/samxsmith/snag/builder"
	"github.com/samxsmith/snag/config"
	"github.com/samxsmith/snag/parser"
)

func schemaToBuilderType(ct config.AggregationType) builder.AggregationType {
	switch ct {
	case config.AggregationTypeCount:
		return builder.AggregationTypeCount
	case config.AggregationTypeList:
		return builder.AggregationTypeList
	default:
	}
	panic("unknown type: ")
}

func schemaToBuilderSortCols(cls []config.SortColumn) []builder.SortCol {
	bcols := make([]builder.SortCol, len(cls))

	for i, col := range cls {
		bc := builder.SortCol{
			Name:          col.Name,
			Type:          schemaToBuilderSortType(col.ColType),
			SortAscending: col.SortAscending,
		}
		bcols[i] = bc
	}

	return bcols
}

func schemaToBuilderSortType(t config.SortType) builder.SortColType {
	switch t {
	case config.SortTypeNumber:
		return builder.SortColTypeNum
	case config.SortTypeDate:
		return builder.SortColTypeDate

	default:
	}
	panic("unknown sort type")
}

func createBuilderEntry(m LineMatch, e parser.Entry) builder.Entry {
	return builder.Entry{
		SourceFilePath: m.FilePath,
		Fields:         e.Fields,
	}
}

func schemaToBuilderSourceLink(c config.LinkType) builder.LinkType {
	if c == config.WikiLink {
		return builder.WikiLink
	}
	if c == config.MarkdownLink {
		return builder.MarkdownLink
	}

	return builder.NoLink
}
