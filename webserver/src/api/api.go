//Package api provides primitives for attaching response functions to http requests.
//All responses use the JSON format.
package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

//RequestData contains information sent alongside the http request
type RequestData struct {
	//Query is a map of the http query
	Query map[string]string

	//Body is the body of a http POST request
	Body []byte

	//Method is the method of the http request
	Method string
}

type response struct {
	Endpoint string
	Callback func(RequestData) []byte
}

var responses []response

//AddResponse attaches a callback function to an endpoint.
//The callback function returns an array of bytes on JSON format
func AddResponse(endpoint string, callback func(RequestData) []byte) {
	responses = append(responses, response{endpoint, callback})
}

//HandleRequests sets up the webserver to listen to all responses previously added
//It also serves static content through the root (/) endpoint
//It sets the content type of all responses to application/json
func HandleRequests(staticContentPath string, port string) {
	fs := http.FileServer(http.Dir(staticContentPath))
	http.Handle("/", fs)

	for _, r := range responses {
		func(e response) {
			http.HandleFunc(e.Endpoint, func(w http.ResponseWriter, r *http.Request) {
				//Read HTTP Query
				query := make(map[string]string)
				for k, v := range r.URL.Query() {
					query[k] = strings.Join(v, "")
				}

				//Read HTTP Body
				buffer := make([]byte, 4096)
				n, _ := r.Body.Read(buffer)

				req := RequestData{
					Query:  query,
					Body:   buffer[:n],
					Method: r.Method,
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(e.Callback(req))
				fmt.Println(time.Now())
				fmt.Println(r.URL)
				fmt.Println(req.Method)
				fmt.Println(req.Query)
				fmt.Println(string(req.Body))
				fmt.Print("\n")
			})
		}(r)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
