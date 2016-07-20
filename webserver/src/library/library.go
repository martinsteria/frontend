//Package library contains functions for reading with a collection of Terraform modules
package library

import (
	"doc"
	"encoding/json"
	"io/ioutil"
	"log"
)

//Library contains a structered representation of a Terraform module collection
type Library struct {
	Modules map[string]*doc.Module
	RootDir string
}

//NewLibrary initializes and builds a library from a path pointing to a Terraform module collection
func NewLibrary(rootDir string) *Library {
	l := &Library{RootDir: rootDir}
	l.Build()
	return l
}

//GetRootDir returns the root directory of the library
func (l *Library) GetRootDir() string {
	return l.RootDir
}

//GetModuleListJSON returns a JSON representation of a list with all modules in the library
//Each module contains an id, a name and a description
func (l *Library) GetModuleListJSON() []byte {
	type module struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	var ms []module

	for _, m := range l.Modules {
		if m.Provider != "" {
			ms = append(ms, module{m.Id, m.Name, m.Description})
		}
	}

	msJSON, _ := json.Marshal(ms)
	return msJSON
}

//GetModuleDocumentationJSON returns a JSON representation of a module
//The representation contains the following:
//name: The Name of the module
//id: A unique identifier of the module
//description: The description of the module
//provider: The provider of the module, if it has one
//variables: A list of variables for the module
//outputs: A list of outputs for the module
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

//Build parses and imports the collection of Terraform modules located at the library's root directory.
func (l *Library) Build() {
	l.Modules = make(map[string]*doc.Module)
	files, err := ioutil.ReadDir(l.RootDir)

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		l.Modules[f.Name()] = doc.NewModule(l.RootDir + "/" + f.Name())
		l.Modules[f.Name()].BuildModule()
	}
}
