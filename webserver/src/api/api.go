package api

import (
	"fmt"
	"log"
	"net/http"
)

const (
	apiBaseURL     = "/api"
	staticBasePath = "/home/martin/go-web/webserver/static/"
)

type route struct {
	Url      string
	Function func() []byte
}

var routes []route

func AddRoute(url string, function func() []byte) {
	routes = append(routes, route{url, function})
}

func HandleRequests() {
	fs := http.FileServer(http.Dir(staticBasePath))
	http.Handle("/", fs)

	for _, r := range routes {
		go func(v route) {
			http.HandleFunc(v.Url, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				n, err := w.Write(v.Function())
				fmt.Println(n)
				fmt.Println(err)
			})
		}(r)
	}

	log.Fatal(http.ListenAndServe(":80", nil))
}
