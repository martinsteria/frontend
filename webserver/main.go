package main

import (
	"api"
	"os"
	"requests"
)

const (
	staticContentPath = "/home/martin/frontend/nettside"
	resourcesRootDir  = "/home/martin/frontend/webserver"
)

func main() {
	requests.Init(resourcesRootDir)

	api.AddResponse("/api", func(api.RequestData) []byte {
		return []byte("{}")
	})

	api.AddResponse("/api/users", requests.HandleUserRequests)
	api.AddResponse("/api/library", requests.HandleLibraryRequests)
	api.AddResponse("/api/library/copy", requests.HandleLibraryCopyRequests)
	api.AddResponse("/api/deploy", requests.HandleDeployRequests)

	api.HandleRequests(staticContentPath, os.Args[1])
}
