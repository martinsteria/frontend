package library

import (
	"doc"
	"encoding/json"
	"io/ioutil"
)

const (
	libraryURL     = "https://github.com/martinsteria/library"
	LibraryRootDir = "home/martin/library"
	LibraryModules = LibraryRootDir + "/modules"
)

type Library struct {
	Modules map[string]*doc.Module
	RootDir string
}

func NewLibrary(rootDir string) *Library {
	l := &Library{RootDir: rootDir}
	l.Build()
	return l
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
	if l.Modules[id] == nil {
		return []byte("{\"status\": \"Module not found\"}")
	}
	var moduleJSON []byte

	type module struct {
		Name        string         `json:"name"`
		Id          string         `json:"id"`
		Description string         `json:"description"`
		Provider    string         `json:"provider"`
		Variables   []doc.Variable `json:"variables"`
		Outputs     []doc.Output   `json:"outputs"`
	}

	m := module{
		Name:        l.Modules[id].Name,
		Id:          l.Modules[id].Id,
		Description: l.Modules[id].Description,
		Provider:    l.Modules[id].Provider,
		Variables:   l.Modules[id].Variables,
		Outputs:     l.Modules[id].Outputs,
	}

	moduleJSON, _ = json.Marshal(m)
	return moduleJSON
}

func (l *Library) Build() {
	l.Modules = make(map[string]*doc.Module)
	files, _ := ioutil.ReadDir(l.RootDir)

	for _, f := range files {
		l.Modules[f.Name()] = doc.NewModule(l.RootDir + "/" + f.Name())
		l.Modules[f.Name()].BuildModule()
	}
}
