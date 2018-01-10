package indexing

import (
	"encoding/base64"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search/highlight/format/html"
	"sync"
)

type Indexer struct {
	idxs map[string]bleve.Index
	sync.RWMutex
}

func (i *Indexer) getIndex(id string) (idx bleve.Index) {
	i.RLock()
	idx = i.idxs[id]
	i.RUnlock()

	if idx != nil {
		return
	}

	i.Lock()
	idx, _ = Bleve(base64.URLEncoding.EncodeToString([]byte(id)))
	i.idxs[id] = idx
	i.Unlock()
	return
}

func (i *Indexer) AddBasicWebpage(ID, CommunityID string, wp WebpageBasic) {
	wp.Index(ID, i.getIndex(CommunityID))
}

func (i *Indexer) AddBasicFeed(ID, CommunityID string, fb FeedBasic) {
	fb.Index(ID, i.getIndex(CommunityID))
}

func (i *Indexer) GetFields(CommunityID string) ([]string, error) {
	return i.getIndex(CommunityID).Fields()
}

func (i *Indexer) Query(id, query string) (*bleve.SearchResult, error) {
	//searchRequest := bleve.NewSearchRequest(bleve.NewMatchQuery(query))
	//searchRequest := bleve.NewSearchRequest(bleve.NewFuzzyQuery(query))
	//searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	searchTerm := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(searchTerm)
	searchRequest.Fields = make([]string, 5)
	searchRequest.Fields[0] = "URL"
	searchRequest.Fields[1] = "Host"
	searchRequest.Fields[2] = "Path"
	searchRequest.Fields[3] = "Title"
	searchRequest.Fields[4] = "Description"
	searchRequest.Highlight = bleve.NewHighlightWithStyle(html.Name)
	return i.getIndex(id).Search(searchRequest)
}

func NewIndexer() Indexer {
	return Indexer{
		idxs: make(map[string]bleve.Index),
	}
}

//var bleveIdx bleve.Index
//var bleveIdxMap = make(map[string]bleve.Index)

// Bleve connect or create the index persistence
func Bleve(indexPath string) (bleve.Index, error) {

	//if bleveIdx, exists := bleveIdxMap[indexPath]; exists {
	//	return bleveIdx, nil
	//}

	// try to open de persistence file...
	bleveIdx, err := bleve.Open("bleve/" + indexPath)

	// if doesn't exists or something goes wrong...
	if err != nil {
		// create a new mapping file and create a new index
		var newMapping mapping.IndexMapping
		newMapping, err = buildIndexMapping()
		if err != nil {
			return nil, err
		}
		bleveIdx, err = bleve.New("bleve/"+indexPath, newMapping)
		if err != nil {
			return nil, err
		}
	}
	return bleveIdx, err
}
