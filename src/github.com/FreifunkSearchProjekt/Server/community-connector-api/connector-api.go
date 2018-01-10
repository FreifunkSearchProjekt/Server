package community_connector_api

import (
	"encoding/json"
	"github.com/FreifunkSearchProjekt/Server/indexing"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RegisterHandler(r *mux.Router, idxr indexing.Indexer) {
	r.HandleFunc("/connector_api/index/{communityID}/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Got new URL to Index: ")
		vars := mux.Vars(r)
		communityID := vars["communityID"]

		var txn indexing.Transaction

		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			w.Write([]byte("{}"))
			return
		}
		err := json.NewDecoder(r.Body).Decode(&txn)
		if err != nil {
			log.Fatalf("[ERR] %s\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			w.Write([]byte("{}"))
			return
		}

		for i := range txn.BasicWebpages {
			log.Println(txn.BasicWebpages[i].URL)
			webpage := indexing.WebpageBasic{
				URL:         txn.BasicWebpages[i].URL,
				Host:        txn.BasicWebpages[i].Host,
				Path:        txn.BasicWebpages[i].Path,
				Title:       txn.BasicWebpages[i].Title,
				Body:        txn.BasicWebpages[i].Body,
				Description: txn.BasicWebpages[i].Description,
			}
			err := idxr.AddBasicWebpage(txn.BasicWebpages[i].URL+txn.BasicWebpages[i].Path, communityID, webpage)
			if err != nil {
				log.Println("[ERR] ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				w.Write([]byte("{}"))
				return
			}
		}

		for i := range txn.RssFeed {
			log.Println(txn.RssFeed[i].URL)
			rssfeed := indexing.FeedBasic{
				URL:         txn.RssFeed[i].URL,
				Host:        txn.RssFeed[i].Host,
				Path:        txn.RssFeed[i].Path,
				Title:       txn.RssFeed[i].Title,
				Description: txn.RssFeed[i].Description,
			}
			err := idxr.AddBasicFeed(txn.RssFeed[i].URL+txn.RssFeed[i].Path, communityID, rssfeed)
			if err != nil {
				log.Println("[ERR] ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				w.Write([]byte("{}"))
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
}
