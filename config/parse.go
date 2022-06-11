package config

import (
	"fmt"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Schemas     []Schema `yaml:"schemas"`
	IgnoreRules []string
}

func (c Config) ShouldIgnore(path string) bool {
	// TODO
	return false
}

type Schema struct {
	// id for internal use -- unique reference to entry
	id int

	Tag             string          `yaml:"tag"`
	Cols            []string        `yaml:"cols"`
	AggregationType AggregationType `yaml:"aggregationType"`
	SortColumns     []SortColumn    `yaml:"sortColumns"`
	OutputFilePath  string          `yaml:"outputFilepath"`
	SourceLink      LinkType        `yaml:"sourceLink"`
}

func (e Schema) ID() int {
	return e.id
}

type AggregationType string

const (
	// AggregationTypeList lists all occurences
	AggregationTypeList AggregationType = "list"

	// AggregationTypeCount counts occurences
	AggregationTypeCount AggregationType = "count"

	// sum: sums a column
)

type SortColumn struct {
	Name          string   `yaml:"name"`
	ColType       SortType `yaml:"type"`
	SortAscending bool     `yaml:"ascending"`
}

type SortType string

const (
	// SortTypeDate
	SortTypeDate SortType = "date"
	// SortTypeNumber
	SortTypeNumber SortType = "number"
)

type LinkType string

const (
	WikiLink     LinkType = "wikiLink"
	MarkdownLink LinkType = "markdown"
)

func Parse(configData []byte, rootPath string) (Config, error) {
	c := Config{}
	err := yaml.UnmarshalStrict(configData, &c)
	if err != nil {
		return c, fmt.Errorf("unable to read config: %w", err)
	}

	for i := range c.Schemas {
		c.Schemas[i].id = i + 1
		c.Schemas[i].OutputFilePath = filepath.Join(rootPath, c.Schemas[i].OutputFilePath)
	}

	return c, nil
}
