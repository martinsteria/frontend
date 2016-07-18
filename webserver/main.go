package main

import (
	"api"
	"os"
	"requests"
)

func main() {
	requests.Init()

	api.AddResponse("/api", func(api.RequestData) []byte {
		return []byte("{}")
	})

	api.AddResponse("/api/users", requests.HandleUserRequests)
	api.AddResponse("/api/library", requests.HandleLibraryRequests)

	api.HandleRequests(os.Args[1])
}
