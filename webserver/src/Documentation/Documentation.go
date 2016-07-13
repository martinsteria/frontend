package Documentation

import (
	"bufio"
	"fmt"
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
	Description string     `json:"description"`
	Variables   []variable `json:"variables"`
	Outputs     []output   `json:"outputs"`
}

func BuildModule(filepath string) Module {
	fmt.Println(filepath)
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	var newModule Module
	var variables []variable
	var outputs []output
	add := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "#") {
			newModule.Name = strings.Trim(line, "#")
		} else if strings.Contains(line, "Module Description") {
			line := strings.Trim(scanner.Text(), "Module Description = ")
			for {
				if strings.Contains(scanner.Text(), "*/") {
					break
				}
				newModule.Description += line
				scanner.Scan()
				line = scanner.Text()
			}

		} else if strings.Contains(line, "variable") {
			line = strings.Trim(line, "variable { ")
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
					temp_variable.DefaultValue = strings.Trim(line, "default = \"")
				} else if strings.Contains(line, "description") {
					line = strings.Trim(line, " description = ")
					temp_variable.Description = strings.Trim(line, "\"")
				}
			}
			variables = append(variables, temp_variable)

		} else if strings.Contains(line, "output") {
			var temp_output output
			line = strings.Trim(line, "output { ")
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

	newModule.Variables = variables
	newModule.Outputs = outputs

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	return newModule
}
