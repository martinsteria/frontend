package main

import (
	"api"
	"fmt"
	"library"
)

func main() {
	library.Init()
	api.AddRoute("/api", func() []byte {
		return []byte("{}")
	})
	api.AddRoute("/api/modules", library.GetModuleNamesJSON)
	for _, m := range library.GetModuleNames() {
		api.AddRoute("/api/modules/"+m, library.GetModuleDocumentation)
	}
	api.HandleRequests()
}
