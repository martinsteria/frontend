package api

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	staticBasePath = "/home/martin/go-web/webserver/static/"
)

type endpoint struct {
	Url     string
	Methods map[string]func() []byte
}

var endpoints []endpoint

func AddEndpoint(url string, methods map[string]func() []byte) {
	endpoints = append(endpoints, endpoint{url, methods})
}

func HandleRequests(port string) {
	fs := http.FileServer(http.Dir(staticBasePath))
	http.Handle("/", fs)

	for _, r := range endpoints {
		go func(e endpoint) {
			http.HandleFunc(e.Url, func(w http.ResponseWriter, r *http.Request) {
				method, present := e.Methods[r.Method]
				if present {
					w.Header().Set("Content-Type", "application/json")
					w.Write(method())
					fmt.Println(string(method()))
					fmt.Println(time.Now())
					fmt.Print("\n")
				} else {
					w.Header().Set("Status", "403")
					w.Write([]byte(""))
				}
			})
		}(r)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
