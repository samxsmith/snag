package config_test

import (
	"testing"

	"github.com/samxsmith/snag/config"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	b := []byte(
		`---
schemas:
        - tag: "#book"
          cols:
          - name
          - date
          - rating
          sortColumns:
                  - name: rating
                    type: number
                    ascending: false
                  - name: date
                    type: date
                    ascending: true

          aggregationType: list
          outputFilepath: "./aggregated/books"
          sourceLink: wikiLink

`)
	res, err := config.Parse(b, "./")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res.Schemas))

	assert.Equal(t, "#book", res.Schemas[0].Tag)
	assert.Equal(t, 3, len(res.Schemas[0].Cols))
	assert.Equal(t, res.Schemas[0].Cols[0], "name")
	assert.Equal(t, res.Schemas[0].Cols[1], "date")
	assert.Equal(t, res.Schemas[0].Cols[2], "rating")
	assert.Equal(t, res.Schemas[0].AggregationType, config.AggregationTypeList)
	assert.Equal(t, res.Schemas[0].ID(), 1)

}
