package library

import (
	"api"
	"doc"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
	"users"
)

const (
	libraryURL        = "https://github.com/martinsteria/library"
	libraryRootFolder = "/home/martin/library"
)

type Library struct {
	Modules []doc.Module
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

func HandleLibraryGetRequests(r api.RequestData) []byte {
	v, present := r.Query["get"]

	if present {
		return GetModuleDocumentationJSON(v[0])
	}

	v, present = r.Query["copy"]
	if present {
		v, present = r.Query["user"]
		if present {
			return CopyModule(r.Query["copy"][0], r.Query["user"][0])
		}
	}

	for k, v := range r.Query {
		if k == "get" {
			return GetModuleDocumentationJSON(v[0])
		}
	}

	return GetModuleListJSON()
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
			doc.BuildModule(libraryRootFolder+"/modules/"+f.Name()))
	}

	return lib
}

func CopyModule(id string, user string) []byte {
	exec.Command("cp", "-r", libraryRootFolder+"/modules/"+id, users.UsersRootFolder+"/"+user).Output()
	return []byte("{\"status\": \"success\"}")
}
