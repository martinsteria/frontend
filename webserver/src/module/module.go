package module

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"terraform"
)

type variable struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	DefaultValue string `json:"defaultValue"`
	Value        string `json:"value"`
}

type output struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Module struct {
	Name        string               `json:"name"`
	Path        string               `json:"path"`
	Id          string               `json:"id"`
	Description string               `json:"description"`
	Provider    string               `json:"provider"`
	Deployment  terraform.Deployment `json:"deployment"`
	Variables   []variable           `json:"variables"`
	Outputs     []output             `json:"outputs"`
}

func NewModule(path string) *Module {
	m := &Module{Path: path}
	m.BuildModule()
	return m
}

func (m *Module) BuildModule() {
	files, _ := ioutil.ReadDir(m.Path)
	var variables []variable
	var outputs []output

	add := false
	for _, f := range files {
		if strings.Contains(f.Name(), ".tf") &&
			!strings.Contains(f.Name(), ".tfvars") &&
			!strings.Contains(f.Name(), ".tfstate") {
			file, err := os.Open(m.Path + "/" + f.Name())
			checkError(err)

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				line = strings.Replace(line, "\"", "", -1)
				if strings.Contains(line, "#") {
					line = strings.Replace(line, "#", "", -1)
					m.Name = strings.TrimSpace(line)

				} else if strings.Contains(line, "Module Description") {
					line := strings.Replace(line, "Module Description = ", "", -1)
					for {
						if strings.Contains(line, "*/") {
							line = strings.Trim(line, "*/")
							m.Description += strings.TrimSpace(line)
							break
						}
						m.Description += strings.TrimSpace(line)
						scanner.Scan()
						line = scanner.Text()
						line = strings.Replace(line, "\"", "", -1)
					}

				} else if strings.Contains(line, "variable") {
					line = strings.Replace(line, "variable", "", -1)
					line = strings.Replace(line, "\"", "", -1)
					line = strings.Trim(line, " { } ")

					var temp_variable variable
					temp_variable.Name = strings.TrimSpace(line)
					for {
						if strings.Contains(scanner.Text(), "}") {
							break
						}
						scanner.Scan()
						line = scanner.Text()
						line = strings.Replace(line, "\"", "", -1)
						if strings.Contains(line, "default") {
							line = strings.Replace(line, "default =", "", -1)
							temp_variable.DefaultValue = strings.TrimSpace(line)
						} else if strings.Contains(line, "description") {
							line = strings.Replace(line, "description =", "", -1)
							temp_variable.Description = strings.TrimSpace(line)
						}
					}
					variables = append(variables, temp_variable)

				} else if strings.Contains(line, "output") {
					var temp_output output
					line = strings.Replace(line, "output", "", -1)
					line = strings.Trim(line, " { } ")
					temp_output.Name = strings.TrimSpace(line)
					for {
						if strings.Contains(line, "}") || strings.Contains(scanner.Text(), "*/") {
							add = false
							break
						} else if add {
							temp_output.Description += strings.TrimSpace(line)
						} else if strings.Contains(line, "Output Description") {
							add = true
							temp_output.Description += strings.Replace(line, "Output Description = ", "", -1)
						}
						scanner.Scan()
						line = scanner.Text()
						line = strings.Replace(line, "\"", "", -1)
					}
					outputs = append(outputs, temp_output)

				} else if strings.Contains(line, "provider") {
					line = strings.Replace(line, "provider", "", -1)
					line = strings.Trim(line, " { } ")
					m.Provider = strings.TrimSpace(line)
				}
				checkError(err)
			}
			err = file.Close()
			checkError(err)
		}
	}
	pathList := strings.Split(strings.TrimSpace(m.Path), "/")
	m.Id = pathList[len(pathList)-1]
	m.Variables = variables
	m.Outputs = outputs
}

//TODO: Integrate with buildmodule
func (m *Module) updateVariableValues() {
	files, _ := ioutil.ReadDir(m.Path)

	for _, f := range files {
		if strings.Contains(f.Name(), ".tfvars") {
			file, err := os.Open(m.Path + "/" + f.Name())
			checkError(err)

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				for i := 0; i < len(m.Variables); i++ {
					if strings.Contains(line, m.Variables[i].Name) {
						line = strings.Replace(line, "\"", "", -1)
						line = strings.Replace(line, m.Variables[i].Name, "", -1)
						line = strings.Replace(line, "=", "", -1)
						m.Variables[i].Value = strings.TrimSpace(line)
					}
				}
				checkError(err)
			}
			err = file.Close()
			checkError(err)
		}
	}
}

func (m *Module) UpdateModule(varsJSON []byte) {
	file, err := os.Create(m.Path + "/terraform.tfvars")

	vars := make([]variable, 5)
	json.Unmarshal(varsJSON, &vars)
	fmt.Println(string(varsJSON))
	fmt.Println(vars)

	checkError(err)
	for i := 0; i < len(vars); i++ {
		if vars[i].Value != "" {
			variable := vars[i].Name + " = \"" + vars[i].Value + "\"" + "\n"
			file.WriteString(variable)
		}
	}

	err = file.Close()
	checkError(err)

	m.updateVariableValues()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Module) GetDocumentationJSON() []byte {
	type module struct {
		Name        string     `json:"name"`
		Id          string     `json:"id"`
		Description string     `json:"description"`
		Provider    string     `json:"provider"`
		Variables   []variable `json:"variables"`
		Outputs     []output   `json:"outputs"`
	}

	mInternal := module{
		Name:        m.Name,
		Id:          m.Id,
		Description: m.Description,
		Provider:    m.Provider,
		Variables:   m.Variables,
		Outputs:     m.Outputs,
	}

	mInternalJSON, _ := json.Marshal(mInternal)
	return mInternalJSON
}

/*
func DeleteKeys(path string, module Module){
	for i := 0; i < len(module.Variables); i++ {
		if strings.Contains(module.Variables[i].Name, "access_key") ||
		strings.Contains(module.Variables[i].Name, "secret_key"){
			module.Variables[i].Value = ""
		}
	}
	CreateTFvars(path, module.Variables)
}
*/
