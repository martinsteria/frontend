package main

import (
	"api"
	"library"
	"os"
)

func main() {
	library.Init()
	api.AddEndpoint("/api", func() []byte {
		return []byte("{}")
	})

	api.AddEndpoint("/api/modules", library.GetModuleListJSON)

	for _, v := range library.GetModuleIds() {
		func(id string) {
			api.AddEndpoint("/api/modules/"+id, func() []byte {
				return library.GetModuleDocumentationJSON(id)
			})
		}(v)
	}

	api.HandleRequests(os.Args[1])
}
