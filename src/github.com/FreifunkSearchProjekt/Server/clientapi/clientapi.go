package clientapi

import (
	"encoding/json"
	"github.com/FreifunkSearchProjekt/Server/database"
	"github.com/FreifunkSearchProjekt/Server/indexing"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var idxrG indexing.Indexer

var SigningKey = []byte("Jl3DyPkeWLjCytk61dXVHLPZcyr8WXwTinPLn3ttgOI6uxNtEffgZxxuMENXfVg4qK5lqgw3AjeKKBVxCTDUMWhi9uWMahPe0s2Y3BMF0x7K2bKE3zyR3DOt2eqhnbPL")

type registerRequest struct {
	CommunityID   string `json:"community_id"`
	CommunityName string `json:"community_name"`
	Homepage      string `json:"homepage"`
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		x := 0
		newString := str[0 : num+x]
		if newString[len(newString)-1:] != " " {
			for newString[len(newString)-1:] != " " {
				x += 1
				newString = str[0 : num+x]
			}
			bnoden = str[0:num+x-1] + "..."
		} else {
			bnoden = newString + "..."
		}
	}
	return bnoden
}

func RegisterHandler(r *mux.Router, idxr indexing.Indexer) {
	idxrG = idxr

	r.HandleFunc("/clientapi/search/{communityID}/{query}/max/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		communityID := vars["communityID"]
		query := vars["query"]
		log.Println("Got new Search Request")

		res, queryErr := idxr.QueryMaxSize(communityID, query)
		if queryErr != nil {
			http.Error(w, queryErr.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(res))
	}).Methods("GET")

	r.HandleFunc("/clientapi/search/{communityID}/{query}/", query).Methods("GET")
	r.HandleFunc("/clientapi/search/{communityID}/{query}/{from}/", query).Methods("GET")

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
	}).Methods("GET")

	r.HandleFunc("/clientapi/account/{communityID}/register", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		communityID := vars["communityID"]

		var txn registerRequest

		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&txn)
		if err != nil {
			log.Fatalf("[ERR] %s\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if txn.CommunityID != communityID {
			http.Error(w, "Community ID doesn't match Body", http.StatusBadRequest)
			return
		}

		used, err := database.CheckIfUserIsAlreadyRegistered(communityID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if used {
			http.Error(w, "User already taken", http.StatusConflict)
			return
		}

		/* Create the token */
		token := jwt.New(jwt.SigningMethodHS256)

		/* Create a map to store our claims */
		claims := token.Claims.(jwt.MapClaims)

		/* Set token claims */
		claims["CommunityName"] = txn.CommunityName
		claims["CommunityID"] = txn.CommunityID

		/* Sign the token with our secret */
		tokenString, _ := token.SignedString(SigningKey)

		saveErr := database.SaveNewUser(tokenString, txn.CommunityID, txn.CommunityName)
		if saveErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		/* Finally, write the token to the browser window */
		w.Write([]byte(tokenString))
	}).Methods("POST")

	/*r.HandleFunc("/clientapi/account/{communityID}/edit", func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Access-Token")
		if accessToken == "" {

		}

		vars := mux.Vars(r)
		communityID := vars["communityID"]

		http.Error(w, "Not implemented", http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	})*/
}

var query = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityID := vars["communityID"]
	query := vars["query"]
	from, found := vars["from"]
	var fromInt int
	if !found {
		fromInt = 0
	} else if found {
		var err error
		fromInt, err = strconv.Atoi(from)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	log.Println("Got new Search Request")

	res, queryErr := idxrG.Query(communityID, query, fromInt)
	if queryErr != nil {
		http.Error(w, queryErr.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range res.Hits {
		if v.Fields["Description"] != nil {
			v.Fields["Description"] = truncateString(v.Fields["Description"].(string), 260)
		}
	}

	hits, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(hits)
}
