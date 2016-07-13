package library

import (
	"Documentation"
	"encoding/json"
	"io/ioutil"
)

const (
	libraryRootFolder = "/home/martin/mockLib/"
)

type Library struct {
	Modules []Documentation.Module
}

var l Library

func Init() {
	l = buildLibrary()
}

func GetModuleNames() []string {
	names := make([]string, len(l.Modules))
	for i := range l.Modules {
		names[i] = l.Modules[i].Name
	}
	return names
}

func GetModuleNamesJSON() []byte {
	namesJSON, _ := json.Marshal(GetModuleNames())
	return namesJSON
}

func GetModuleDocumentation(name string) []byte {
	var moduleJSON []byte
	for _, m := range l.Modules {
		if name == m.Name {
			moduleJSON, _ = json.Marshal(m)
		}
	}

	return moduleJSON
}

func buildLibrary() Library {
	var lib Library

	files, _ := ioutil.ReadDir(libraryRootFolder + "/modules")

	for _, f := range files {
		lib.Modules = append(
			lib.Modules,
			Documentation.BuildModule(libraryRootFolder+"/modules/"+f.Name()+"/main.tf"))
	}

	return lib
}
