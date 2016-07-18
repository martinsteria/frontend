package terraform

import (
	"log"
	"os/exec"
)

type Deployment struct {
	Status string `json:"status"`
	Path   string `json:"path"`
	Output string `json:"output"`
}

func (t *Deployment) Init(path string) {
	t.Path = path
}

func (t *Deployment) TerraformCommand(command string) {

	t.Status = "Running"
	t.Output = "Running"
	t.getModules() // SHOULD BE PUT SOMEWHERE ELSE!!
	if command == "destroy" {
		cmd := exec.Command("terraform", command, "-force")
		cmd.Dir = t.Path
		out, _ := cmd.CombinedOutput()
		t.Status = ""
		t.Output = string(out)
	}

	cmd := exec.Command("terraform", command)
	cmd.Dir = t.Path
	out, err := cmd.CombinedOutput()

	log.Println(command)
	log.Println(t.Path)

	if err != nil {
		log.Fatal(err)
	}
	t.Status = ""
	t.Output = string(out)

	//DELETE KEYS?????
}

func (t *Deployment) getModules() {

	init := exec.Command("terraform", "get")
	init.Dir = t.Path
	init.CombinedOutput()
}
