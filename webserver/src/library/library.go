package library

import (
	"Documentation"
	"encoding/json"
	"io/ioutil"
)

const (
	libraryURL     = "https://github.com/martinsteria/library"
	LibraryRootDir = "/home/martin/library"
	LibraryModules = LibraryRootDir + "/modules"
)

type Library struct {
	Modules []Documentation.Module
}

func (l *Library) GetModuleListJSON() []byte {
	type module struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	var ms []module

	for _, m := range l.Modules {
		ms = append(ms, module{m.Id, m.Name, m.Description})
	}

	msJSON, _ := json.Marshal(ms)
	return msJSON
}

func (l *Library) GetModuleDocumentationJSON(id string) []byte {
	var moduleJSON []byte
	for _, m := range l.Modules {
		if id == m.Id {
			moduleJSON, _ = json.Marshal(m)
		}
	}

	return moduleJSON
}

func BuildLibrary(rootDir string) Library {
	var lib Library

	files, _ := ioutil.ReadDir(rootDir)

	for _, f := range files {
		lib.Modules = append(
			lib.Modules,
			Documentation.BuildModule(rootDir+"/"+f.Name()))
	}

	return lib
}
