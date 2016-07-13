package main

import (
	"./api"
	"./library"
)

func main() {
	library.Init()
	api.AddRoute("/api", func() []byte {
		return []byte("{}")
	})
	api.AddRoute("/api/modules", library.GetModuleNames)
	api.HandleRequests()
}
