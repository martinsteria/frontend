package library

import (
	"Documentation"
	"encoding/json"
	"io/ioutil"
)

const (
	libraryRootFolder = "/home/martin/terraform/martin"
)

type Library struct {
	Modules []Documentation.Module
}

var l Library

func Init() {
	l = buildLibrary()
}

func GetModuleIds() []string {
	var Ids []string
	for _, m := range l.Modules {
		Ids = append(Ids, m.Id)
	}
	return Ids
}

func GetModuleListJSON() []byte {
	type module struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	var ms []module

	for _, m := range l.Modules {
		ms = append(ms, module{m.Id, m.Name})
	}

	msJSON, _ := json.Marshal(ms)
	return msJSON
}

func GetModuleDocumentationJSON(id string) []byte {
	var moduleJSON []byte
	for _, m := range l.Modules {
		if id == m.Id {
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
			Documentation.BuildModule(libraryRootFolder+"/modules/"+f.Name()))
	}

	return lib
}
