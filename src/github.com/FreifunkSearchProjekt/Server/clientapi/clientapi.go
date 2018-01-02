package clientapi

import (
	"encoding/json"
	"github.com/FreifunkSearchProjekt/Server/indexing"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHandler(r *mux.Router, idxr indexing.Indexer) {
	r.HandleFunc("/clientapi/search/{communityID}/{query}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		communityID := vars["communityID"]
		query := vars["query"]
		res, _ := idxr.Query(communityID, query)

		hits, err := json.Marshal(res.Hits)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(hits)
	})
}
