package api

import (
	"../library"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	apiBaseURL     = "/api"
	staticBasePath = "/home/martin/go-web/webserver/static/"
)

var lib library.Library

func HandleRequests() {
	lib = library.ParseModulesFromLibraryMock()

	fs := http.FileServer(http.Dir(staticBasePath))
	http.Handle("/", fs)

	http.HandleFunc(apiBaseURL, apiRequestHandler)
	http.HandleFunc(apiBaseURL+"/modules", modulesRequestHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}

func apiRequestHandler(w http.ResponseWriter, r *http.Request) {

}

func modulesRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		names := lib.GetModuleNames()
		namesJSON, _ := json.Marshal(names)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(namesJSON))
	}
	fmt.Println(r.Method)
	fmt.Println(r.URL)
	fmt.Println(r.URL.Query())
	fmt.Println(r.Header)
	fmt.Println(r.Body)
}
