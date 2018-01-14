package common

import (
	"flag"
	"github.com/FreifunkSearchProjekt/Server/clientapi"
	"github.com/FreifunkSearchProjekt/Server/community-connector-api"
	"github.com/FreifunkSearchProjekt/Server/config"
	"github.com/FreifunkSearchProjekt/Server/indexing"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"net/http/pprof"
)

var pathPtr = flag.String("path", "my_registration_file.yaml", "The path to which to write the generated Registration YAML")

//var configPathPtr = flag.String("config", "config.yaml", "The path to the matrix-search config YAML")
var PprofEnabledPtr = flag.Bool("pprof", false, "Whether or not to enable Pprof debugging")

func LoadConfigs() (conf *config.Config) {
	flag.Parse()

	//if conf, err = config.LoadConfig(*configPathPtr); err != nil {
	//	panic(err)
	//}

	return
}

func Setup() (r *mux.Router) {
	idxr := indexing.NewIndexer()

	r = mux.NewRouter()
	if *PprofEnabledPtr {
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	}

	clientapi.RegisterHandler(r, idxr)

	community_connector_api.RegisterHandler(r, idxr)

	return
}

func Begin(cleanHandler http.Handler, conf *config.Config) {
	handler := cors.Default().Handler(cleanHandler)
	srv4 := &http.Server{
		Handler: handler,
		Addr:    "0.0.0.0:9999",
	}
	log.Fatal(srv4.ListenAndServe())

	srv6 := &http.Server{
		Handler: handler,
		Addr:    "[::]:9999",
	}
	log.Fatal(srv6.ListenAndServe())
}
