package Documentation

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

type Module struct {
	Name        string     `json:"name"`
	Id          string     `json:"id"`
	Description string     `json:"description"`
	Variables   []variable `json:"variables"`
	Outputs     []output   `json:"outputs"`
}

func BuildModule(path string) Module {
	files, _ := ioutil.ReadDir(path)

	var newModule Module
	var variables []variable
	var outputs []output
	add := false
	for _, f := range files {
		if strings.Contains(f.Name(), ".tf") &&
			!strings.Contains(f.Name(), ".tfvars") &&
			!strings.Contains(f.Name(), ".tfstate") {
			file, err := os.Open(path + "/" + f.Name())

			if err != nil {
				log.Fatal(err)
			}

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				line = strings.Replace(line, "\"", "", -1)
				if strings.Contains(line, "#") {
					line = strings.Replace(line, "#", "", -1)
					newModule.Name = strings.TrimSpace(line)

				} else if strings.Contains(line, "Module Description") {
					line := strings.Replace(line, "Module Description = ", "", -1)
					for {
						if strings.Contains(line, "*/") {
							line = strings.Trim(line, "*/")
							newModule.Description += strings.TrimSpace(line)
							break
						}
						newModule.Description += strings.TrimSpace(line) //To get rid of blank enters.
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
				}

				if err = scanner.Err(); err != nil {
					log.Fatal(err)
				}

			}
			err = file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	newModule.Id = strings.Replace(newModule.Name, " ", "", -1)
	newModule.Variables = variables
	newModule.Outputs = outputs

	return newModule
}
