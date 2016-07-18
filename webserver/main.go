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

	api.AddResponse("/api/users", users.HandleUserRequests)
	api.AddResponse("/api/library", users.HandleLibraryRequests)

	api.HandleRequests(os.Args[1])
}
