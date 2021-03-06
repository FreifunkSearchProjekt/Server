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

func (i *Indexer) getIndex(id string) (idx bleve.Index, err error) {
	i.RLock()
	idx = i.idxs[id]
	i.RUnlock()

	if idx != nil {
		return
	}

	i.Lock()
	var BleveErr error
	idx, BleveErr = Bleve(base64.URLEncoding.EncodeToString([]byte(id)))
	if BleveErr != nil {
		err = BleveErr
		return
	}
	i.idxs[id] = idx
	i.Unlock()
	return
}

func (i *Indexer) AddBasicWebpage(ID, CommunityID string, wp WebpageBasic) error {
	index, err := i.getIndex(CommunityID)
	if err != nil {
		return err
	}
	wp.Index(ID, index)
	return nil
}

func (i *Indexer) AddBasicFeed(ID, CommunityID string, fb FeedBasic) error {
	index, err := i.getIndex(CommunityID)
	if err != nil {
		return err
	}
	fb.Index(ID, index)
	return nil
}

func (i *Indexer) GetFields(CommunityID string) ([]string, error) {
	index, err := i.getIndex(CommunityID)
	if err != nil {
		return nil, err
	}
	return index.Fields()
}

func (i *Indexer) Query(id, query string) (*bleve.SearchResult, error) {
	//searchRequest := bleve.NewSearchRequest(bleve.NewMatchQuery(query))
	//searchRequest := bleve.NewSearchRequest(bleve.NewFuzzyQuery(query))
	//searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	searchTerm := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(searchTerm)
	searchRequest.Fields = make([]string, 5)
	searchRequest.Fields[0] = "url"
	searchRequest.Fields[1] = "host"
	searchRequest.Fields[2] = "path"
	searchRequest.Fields[3] = "title"
	searchRequest.Fields[4] = "description"
	searchRequest.Highlight = bleve.NewHighlightWithStyle(html.Name)
	index, err := i.getIndex(id)
	if err != nil {
		return nil, err
	}
	return index.Search(searchRequest)
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
	return bleveIdx, nil
}
