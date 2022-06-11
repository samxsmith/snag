package snag

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/samxsmith/snag/builder"
	"github.com/samxsmith/snag/config"
	"github.com/samxsmith/snag/parser"
)

func Run(rootPath string) {
	const configFileName = ".snag.yml"
	configFilePath := filepath.Join(rootPath, configFileName)
	configFile, err := os.Open(configFilePath)
	if err != nil {
		// TODO: pretty handling
		panic(err)
	}
	defer configFile.Close()

	configFileB, err := io.ReadAll(configFile)
	if err != nil {
		// TODO: pretty handling
		panic(err)
	}

	cfg, err := config.Parse(configFileB, rootPath)
	if err != nil {
		// TODO: pretty handling
		panic(err)
	}

	ParseFiles(rootPath, cfg)
}

type LineMatch struct {
	TagID    int
	Line     string
	FilePath string
}

func ParseFiles(rootPath string, c config.Config) {
	wg := sync.WaitGroup{}
	wg.Add(len(c.Schemas))
	tags := []parser.Tag{}

	actorChans := map[int]chan LineMatch{}
	for _, schema := range c.Schemas {
		pipeline, ch := createActorPipeline(schema)
		actorChans[schema.ID()] = ch

		tags = append(tags, parser.Tag{
			ID:     schema.ID(),
			Prefix: schema.Tag,
		})

		go func() {
			defer wg.Done()
			pipeline.Start()
		}()

	}

	filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if c.ShouldIgnore(path) {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		matches, err := parser.FindLinesWithTags(tags, f)
		if err != nil {
			return err
		}

		for _, match := range matches {
			l := LineMatch{
				TagID:    match.TagID,
				Line:     match.Line,
				FilePath: path,
			}

			actorChans[match.TagID] <- l
		}

		return nil
	})

	for _, c := range actorChans {
		close(c)
	}
	wg.Wait()
}

type Pipeline struct {
	InputChan    chan LineMatch
	builder      builder.Builder
	parserSchema parser.EntrySchema
}

func createActorPipeline(c config.Schema) (Pipeline, chan LineMatch) {
	p := Pipeline{}

	p.InputChan = make(chan LineMatch)

	builderConfig := builder.Config{
		AggregationType: schemaToBuilderType(c.AggregationType),
		SortCols:        schemaToBuilderSortCols(c.SortColumns),
		OutputFilePath:  c.OutputFilePath,
		AllColumnNames:  c.Cols,
		SourceLink:      schemaToBuilderSourceLink(c.SourceLink),
	}

	p.builder = builder.NewBuilder(builderConfig)

	p.parserSchema = parser.EntrySchema{
		Cols:   c.Cols,
		Prefix: c.Tag,
	}

	return p, p.InputChan
}

func (p Pipeline) Start() {
	for m := range p.InputChan {
		// parse into fields
		e, err := parser.ParseEntry(m.Line, p.parserSchema)
		if err != nil {
			fmt.Printf("couldn't parse line <%s>: %s \n", m.Line, err)
			continue
		}

		builderEntry := createBuilderEntry(m, e)

		p.builder.Add(builderEntry)
	}
	p.builder.Done()
}
