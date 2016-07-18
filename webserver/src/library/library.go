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
	RootDir string
}

func (l *Library) Init(rootDir string) {
	l.RootDir = rootDir
	l.Build()
}

func (l *Library) GetRootDir() string {
	return l.RootDir
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
			break
		}
	}

	return moduleJSON
}


func (l *Library) Build() {
	l.Modules = make(map[string]*doc.Module)
	files, _ := ioutil.ReadDir(l.RootDir)

	for _, f := range files {
		l.Modules[f.Name()].Init(LibraryModules + "/" + f.Name())
		l.Modules[f.Name()].BuildModule()
	}
}
