package indexing

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/char/html"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
)

func buildIndexMapping() (mapping.IndexMapping, error) {

	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	htmlFieldMapping := bleve.NewTextFieldMapping()
	htmlFieldMapping.Analyzer = html.Name

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	basicPageMapping := bleve.NewDocumentMapping()

	// url
	basicPageMapping.AddFieldMappingsAt("URL", keywordFieldMapping)
	basicPageMapping.AddFieldMappingsAt("Path", keywordFieldMapping)

	// body
	basicPageMapping.AddFieldMappingsAt("Body",
		englishTextFieldMapping)

	// Description
	basicPageMapping.AddFieldMappingsAt("Description",
		englishTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("basicPage", basicPageMapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}
