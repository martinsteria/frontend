package main

import (
	"api"
	"library"
	"os"
	"users"
)

func main() {
	library.Init()

	api.AddResponse("/api", map[string]func(api.RequestData) []byte{
		"GET": func(r api.RequestData) []byte {
			return []byte("{}")
		},
	})

	api.AddResponse("/api/users", map[string]func(api.RequestData) []byte{
		"GET": users.HandleUserRequests,
	})

	api.AddResponse("/api/library", map[string]func(api.RequestData) []byte{
		"GET": library.HandleLibraryGetRequests,
	})

	api.HandleRequests(os.Args[1])
}
