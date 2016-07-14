package api

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	staticBasePath = "/home/martin/frontend/webserver/static/"
)

type RequestData struct {
	Query map[string][]string
	Body  []byte
}

type response struct {
	Endpoint string
	Methods  map[string]func(RequestData) []byte
}

var responses []response

func AddResponse(endpoint string, methods map[string]func(RequestData) []byte) {
	responses = append(responses, response{endpoint, methods})
}

func HandleRequests(port string) {
	fs := http.FileServer(http.Dir(staticBasePath))
	http.Handle("/", fs)

	for _, r := range responses {
		func(e response) {
			http.HandleFunc(e.Endpoint, func(w http.ResponseWriter, r *http.Request) {
				method, present := e.Methods[r.Method]
				if present {
					buffer := make([]byte, 4096)
					n, _ := r.Body.Read(buffer)
					req := RequestData{
						Query: r.URL.Query(),
						Body:  buffer[:n],
					}
					w.Header().Set("Content-Type", "application/json")
					w.Write(method(req))
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
