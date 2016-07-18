package library

import (
	"doc"
	"encoding/json"
	"io/ioutil"
)

const (
	libraryURL     = "https://github.com/martinsteria/library"
	LibraryRootDir = "/home/martin/library"
	LibraryModules = LibraryRootDir + "/modules"
)

type Library struct {
	Modules map[string]*doc.Module
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
	lib.Modules = make(map[string]*doc.Module)
	files, _ := ioutil.ReadDir(rootDir)

	for _, f := range files {
		lib.Modules[f.Name()].Path = (LibraryModules + "/" + f.Name())
		lib.Modules[f.Name()].BuildModule()
	}

	return lib
}
