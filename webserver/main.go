package main

import (
	"api"
	"os"
	"users"
)

func main() {
	users.Init()

	api.AddResponse("/api", func(api.RequestData) []byte {
		return []byte("{}")
	})

	api.AddResponse("/api/users", func(r api.RequestData) []byte {
		return users.HandleUserRequests(r)
	})

	api.AddResponse("/api/library", func(r api.RequestData) []byte {
		return users.HandleLibraryRequests(r)
	})

	api.HandleRequests(os.Args[1])
}
