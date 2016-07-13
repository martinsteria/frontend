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

	api.AddEndpoint("/api/modules", library.GetModuleNamesJSON)

	for _, m := range library.GetModuleNames() {
		func(moduleName string) {
			api.AddEndpoint("/api/modules/"+moduleName, func() []byte {
				return library.GetModuleDocumentation(moduleName)
			})
		}(m)
	}

	api.HandleRequests(os.Args[1])
}
