package library

import (
	"Documentation"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

const (
	libraryURL        = "https://github.com/martinsteria/library"
	libraryRootFolder = "/home/martin/library"
)

type Library struct {
	Modules []Documentation.Module
}

var l Library

func Init() {
	l = updateLibrary()
	go libraryUpdater()

	fmt.Println("Library initialized successfully")
}

func libraryUpdater() {
	for {
		time.Sleep(24 * time.Hour)

		//UNSAFE: Trouble if someone accesses the library during the update
		l = updateLibrary()
	}
}

func updateLibrary() Library {
	currentPath, _ := os.Getwd()
	err := os.Chdir(libraryRootFolder)

	if err != nil {
		cloneLibrary()
		os.Chdir(currentPath)
	} else {
		pullLibrary()
	}
	os.Chdir(currentPath)

	return buildLibrary()
}

func cloneLibrary() {
	out, _ := exec.Command("git", "clone", libraryURL, libraryRootFolder).CombinedOutput()
	fmt.Println(string(out))
}

func pullLibrary() {
	out, _ := exec.Command("git", "pull", "origin", "master").CombinedOutput()
	fmt.Println(string(out))
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
