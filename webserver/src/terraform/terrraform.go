package terraform

import (
	"log"
	"os/exec"
)

func TerraformCommand(command string, path string) string {

	terraformInit(path) // SHOULD BE PUT SOMEWHERE ELSE!!

	if command == "destroy" {
		cmd := exec.Command("terraform", command, "-force")
		cmd.Dir = path
		out, _ := cmd.CombinedOutput()
		return (string(out))
	}

	cmd := exec.Command("terraform", command)
	cmd.Dir = path
	out, err := cmd.CombinedOutput()

	log.Println(command)
	log.Println(path)

	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func terraformInit(path string) {

	init := exec.Command("terraform", "get")
	init.Dir = path
	init.CombinedOutput()
}
