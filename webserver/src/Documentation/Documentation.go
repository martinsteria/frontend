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
		if strings.Contains(f.Name(), ".tf") && !strings.Contains(f.Name(), ".tfvars") {
			file, err := os.Open(f.Name())
			if err != nil {
				log.Fatal(err)
			}
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, "#") {
					newModule.Name = strings.Replace(line, "#", "", -1)
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
					}

				} else if strings.Contains(line, "variable") {
					line = strings.Trim(line, "variable { } ")
					line = strings.Trim(line, "\"")

					var temp_variable variable
					temp_variable.Name = line
					for {
						if strings.Contains(scanner.Text(), "}") {
							break
						}
						scanner.Scan()
						line = scanner.Text()
						if strings.Contains(line, "default") {
							line = strings.Replace(line, "default = ", "", -1)
							temp_variable.DefaultValue = strings.Trim(line, " \"")
						} else if strings.Contains(line, "description") {
							line = strings.Trim(line, " description = ")
							temp_variable.Description = strings.Trim(line, "\"")
						}
					}
					variables = append(variables, temp_variable)

				} else if strings.Contains(line, "output") {
					var temp_output output
					line = strings.Trim(line, "output { } ")
					line = strings.Trim(line, "\"")
					temp_output.Name = strings.Trim(line, "\"")
					for {
						if strings.Contains(scanner.Text(), "}") || strings.Contains(scanner.Text(), "*/") {
							add = false
							break
						} else if add {
							temp_output.Description += strings.Trim(scanner.Text(), "Output Description = ")
						} else if strings.Contains(scanner.Text(), "Output Description") {
							add = true
							temp_output.Description += scanner.Text()
						}
						scanner.Scan()
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
