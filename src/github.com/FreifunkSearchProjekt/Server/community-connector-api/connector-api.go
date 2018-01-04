package community_connector_api

import (
	"encoding/json"
	"github.com/FreifunkSearchProjekt/Server/indexing"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type transaction struct {
	BasicWebpages []indexing.WebpageBasic `json:"basic_webpages"`
}

func RegisterHandler(r *mux.Router, idxr indexing.Indexer) {
	r.HandleFunc("/connector_api/index/{communityID}/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		communityID := vars["communityID"]

		var txn transaction

		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			w.Write([]byte("{}"))
			return
		}
		err := json.NewDecoder(r.Body).Decode(&txn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			w.Write([]byte("{}"))
			return
		}

		for i := range txn.BasicWebpages {
			log.Println(txn.BasicWebpages[i])
			webpage := indexing.WebpageBasic{
				URL:         txn.BasicWebpages[i].URL,
				Path:        txn.BasicWebpages[i].Path,
				Title:       txn.BasicWebpages[i].Title,
				Body:        txn.BasicWebpages[i].Body,
				Description: txn.BasicWebpages[i].Description,
			}
			idxr.AddBasicWebpage(txn.BasicWebpages[i].URL+txn.BasicWebpages[i].Path, communityID, webpage)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
}
