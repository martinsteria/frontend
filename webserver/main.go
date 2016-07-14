package main

import (
	"api"
	"library"
	"os"
)

func main() {
	library.Init()
	api.AddEndpoint("/api", map[string]func() []byte{
		"GET": func() []byte {
			return []byte("{}")
		},
	})

	api.AddEndpoint("/api/modules", map[string]func() []byte{
		"GET": library.GetModuleListJSON,
	})

	for _, v := range library.GetModuleIds() {
		func(id string) {
			api.AddEndpoint("/api/modules/"+id, map[string]func() []byte{
				"GET": func() []byte {
					return library.GetModuleDocumentationJSON(id)
				},
			})
		}(v)
	}

	api.HandleRequests(os.Args[1])
}
