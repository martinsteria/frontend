package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	libraryRootFolder = "~/mockLib/"
)

type variable struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	DefaultValue string `json:"defaultValue"`
}

type output struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type module struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Variables   []variable `json:"variables"`
	Outputs     []output   `json:"outputs"`
}

type moduleList struct {
	Modules []module `json:"modules"`
}

func parseModulesFromLibrary(rootFolder string) *moduleList {
	var modules moduleList

	//TODO: Run doctool to extract JSON representation for all  modules

	//MOCK
	modulesJSON, _ := ioutil.ReadFile("mockTerraform.json")
	err := json.Unmarshal(modulesJSON, &modules)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(modules.Modules[0].Description)

	return &modules
}

func main() {
	parseModulesFromLibrary("a")
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, _ := ioutil.ReadFile("mockTerraform.json")
		fmt.Fprintf(w, string(b))
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":80", nil))
}
