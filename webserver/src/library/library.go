package library

import (
	"Documentation"
	"encoding/json"
	"fmt"
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

func GetModuleDocumentation() []byte {
	var moduleJSON []byte
	moduleJSON, _ = json.Marshal(l.Modules[1])
	return moduleJSON
}

func buildLibrary() Library {
	var lib Library

	files, _ := ioutil.ReadDir(libraryRootFolder + "/modules")

	for _, f := range files {
		fmt.Print("\n\n\n")
		lib.Modules = append(
			lib.Modules,
			Documentation.BuildModule(libraryRootFolder+"/modules/"+f.Name()+"/main.tf"))
	}

	fmt.Println(lib)

	return lib
}
