package indexing

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/char/html"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/blevex/detectlang"
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

	// a specific mapping to index the description fields
	// detected language
	descriptionLangFieldMapping := bleve.NewTextFieldMapping()
	descriptionLangFieldMapping.Name = "descriptionLang"
	descriptionLangFieldMapping.Analyzer = detectlang.AnalyzerName

	basicPageMapping := bleve.NewDocumentMapping()

	// url
	basicPageMapping.AddFieldMappingsAt("URL", keywordFieldMapping)
	basicPageMapping.AddFieldMappingsAt("Path", keywordFieldMapping)

	// Title
	basicPageMapping.AddFieldMappingsAt("Title", keywordFieldMapping)

	// body
	basicPageMapping.AddFieldMappingsAt("Body",
		englishTextFieldMapping)

	// Description
	basicPageMapping.AddFieldMappingsAt("Description",
		descriptionLangFieldMapping)

	basicFeedMapping := bleve.NewDocumentMapping()

	// url
	basicFeedMapping.AddFieldMappingsAt("URL", keywordFieldMapping)
	basicFeedMapping.AddFieldMappingsAt("Path", keywordFieldMapping)

	// Title
	basicFeedMapping.AddFieldMappingsAt("Title", keywordFieldMapping)

	// Description
	basicFeedMapping.AddFieldMappingsAt("Description",
		descriptionLangFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("basicWebpage", basicPageMapping)
	indexMapping.AddDocumentMapping("basicFeed", basicFeedMapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}
