package indexing

import "github.com/blevesearch/bleve"

type Transaction struct {
	BasicWebpages []WebpageBasic `json:"basic_webpages"`
	Feeds      []FeedBasic    `json:"feeds"`
}

type WebpageBasic struct {
	URL         string `json:"url"`
	Host        string `json:"host"`
	Path        string `json:"path"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Description string `json:"description"`
}

func (wp *WebpageBasic) Type() string {
	return "basicWebpage"
}

// Index is used to add the Webpage in the bleve index.
func (wp *WebpageBasic) Index(ID string, index bleve.Index) error {
	err := index.Index(ID, wp)
	return err
}

type FeedBasic struct {
	URL         string `json:"url"`
	Host        string `json:"host"`
	Path        string `json:"path"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (wp *FeedBasic) Type() string {
	return "basicFeed"
}

// Index is used to add the Feed in the bleve index.
func (wp *FeedBasic) Index(ID string, index bleve.Index) error {
	err := index.Index(ID, wp)
	return err
}
