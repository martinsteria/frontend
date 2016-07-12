package library

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	libraryRootFolder = "~/mockLib/"
)

type variable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     string `json:"default"`
}

type output struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type module struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Variables   []variable `json:"variables"`
	Outputs     []output   `json:"outputs"`
}

type Library struct {
	Modules []module
}

func ParseModulesFromLibraryMock() Library {
	var lib Library
	lib.Modules = make([]module, 1)

	moduleJSON, _ := ioutil.ReadFile("../mockTerraform.json")
	err := json.Unmarshal(moduleJSON, &lib.Modules[0])
	if err != nil {
		log.Fatal(err)
	}

	return lib
}

func (l *Library) GetModuleNames() []string {
	names := make([]string, len(l.Modules))
	for i := range l.Modules {
		names[i] = l.Modules[i].Name
	}

	return names
}
