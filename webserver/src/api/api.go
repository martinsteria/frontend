package api

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	apiBaseURL     = "/api"
	staticBasePath = "/home/martin/go-web/webserver/static/"
)

type endpoint struct {
	Url      string
	Function func() []byte
}

var endpoints []endpoint

func AddEndpoint(url string, function func() []byte) {
	endpoints = append(endpoints, endpoint{url, function})
}

func HandleRequests(port string) {
	fs := http.FileServer(http.Dir(staticBasePath))
	http.Handle("/", fs)

	for _, r := range endpoints {
		go func(e endpoint) {
			http.HandleFunc(e.Url, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				n, err := w.Write(e.Function())
				fmt.Println(time.Now())
				fmt.Println(string(e.Function()))
				fmt.Println(n)
				fmt.Println(err)
				fmt.Print("\n")
			})
		}(r)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
