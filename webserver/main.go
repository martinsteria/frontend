package main

import (
	/*
	"io/ioutil"
	"log"
	"net/http"*/
	"Documentation"
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

func parseModulesFromLibrary(rootFolder string) []module {
	var modules []module = make([]module, 1)

	//TODO: Run doctool to extract JSON representation for all  modules

	//MOCK
	moduleJSON, _ := ioutil.ReadFile("mockTerraform.json")
	err := json.Unmarshal(moduleJSON, &modules[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(modules[0].Description)

	return modules
}

func main() {

	/*http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {

	parseModulesFromLibrary("a")
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		b, _ := ioutil.ReadFile("mockTerraform.json")
		fmt.Fprintf(w, string(b))
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":80", nil))*/
	fmt.Printf(Documentation.ExtractDocumentation("test.tf"))
}
