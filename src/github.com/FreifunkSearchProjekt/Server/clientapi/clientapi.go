package clientapi

import (
	"encoding/json"
	"github.com/FreifunkSearchProjekt/Server/indexing"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RegisterHandler(r *mux.Router, idxr indexing.Indexer) {
	r.HandleFunc("/clientapi/search/{communityID}/{query}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		communityID := vars["communityID"]
		query := vars["query"]
		log.Println("Got new Search Request")

		res, queryErr := idxr.Query(communityID, query)
		if queryErr != nil {
			http.Error(w, queryErr.Error(), http.StatusInternalServerError)
			return
		}

		hits, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(hits)
	})

	r.HandleFunc("/clientapi/fields/{communityID}/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		communityID := vars["communityID"]
		log.Println("Got new Fields Request")

		res, queryErr := idxr.GetFields(communityID)
		if queryErr != nil {
			http.Error(w, queryErr.Error(), http.StatusInternalServerError)
			return
		}

		hits, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(hits)
	})
}
