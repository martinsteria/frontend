//package terraform provides functionality to launch a terraform command and retrieve its output incrementally in real time
package terraform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

//Command is the structure used to launch terraform commands
type Command struct {
	status string
	output bytes.Buffer
	error  bytes.Buffer
	path   string
	runner *exec.Cmd
}

//NewCommand returns a new command. The path determines the directory where commands are launched
func NewCommand(path string) *Command {
	d := &Command{path: path}
	return d
}

//Launch starts a command in its own goroutine. If a command is already running, it fails.
func (d *Command) Launch(command string) []byte {
	if d.status == "Running" {
		return []byte("{\"status:\": \"failure: busy\"}")
	}

	d.output.Reset()
	d.error.Reset()

	d.status = "Running"
	go d.commandRunner(command)

	return []byte("{\"status:\": \"success\"}")
}

//IsRunning returns whether a command is currently running or not
func (d *Command) IsRunning() bool {
	if d.status == "Running" {
		return true
	}
	return false
}

//Returns the status of the command including output so far.
func (d *Command) GetStatusJSON() []byte {
	type output struct {
		Status string `json:"status"`
		Output string `json:"output"`
		Error  string `json:"error"`
	}

	o := output{
		Status: d.status,
		Output: strings.Replace(d.output.String(), "\n", "<br />", -1),
		Error:  strings.Replace(d.error.String(), "\n", "<br />", -1),
	}

	oJSON, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}

	return oJSON
}

func (d *Command) commandRunner(command string) {
	fmt.Println("Running ", command, " in ", d.path)
	args := []string{command}

	if command == "plan" || command == "apply" || command == "destroy" {
		args = append(args, "-input=false", "-no-color")
	}

	if command == "destroy" {
		args = append(args, "-force")
	}

	d.runner = exec.Command("terraform", args...)
	d.runner.Dir = d.path

	stdout, err := d.runner.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderr, err := d.runner.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := d.runner.Start(); err != nil {
		log.Fatal(err)
	}
	defer d.runner.Wait()

	d.output.ReadFrom(stdout)
	d.error.ReadFrom(stderr)
	d.status = "Finished"
}
