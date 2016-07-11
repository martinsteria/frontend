package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", r)
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
