package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	staticBasePath = "/home/martin/go-web/webserver/static/"
)

type RequestData struct {
	Query  map[string]string
	Body   []byte
	Method string
}

type response struct {
	Endpoint string
	Callback func(RequestData) []byte
}

var responses []response

func AddResponse(endpoint string, callback func(RequestData) []byte) {
	responses = append(responses, response{endpoint, callback})
}

func HandleRequests(port string) {
	fs := http.FileServer(http.Dir(staticBasePath))
	http.Handle("/", fs)

	for _, r := range responses {
		func(e response) {
			http.HandleFunc(e.Endpoint, func(w http.ResponseWriter, r *http.Request) {
				buffer := make([]byte, 4096)
				n, _ := r.Body.Read(buffer)

				query := make(map[string]string)
				for k, v := range r.URL.Query() {
					query[k] = strings.Join(v, "")
				}
				req := RequestData{
					Query:  query,
					Body:   buffer[:n],
					Method: r.Method,
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(e.Callback(req))
				fmt.Println(time.Now())
				fmt.Println(req.Method)
				fmt.Println(req.Query)
				fmt.Println(string(req.Body))
				fmt.Print("\n")
			})
		}(r)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
