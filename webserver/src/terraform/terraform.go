package terraform

import (
	//"log"
	"bytes"
	"io"
	"os/exec"
	"strings"
	"time" 
	"encoding/json"
	)


type Deployment struct {
	Status        string `json:"status"`
	Path          string `json:"path"`
	Output        []byte `json:"output"`
	buf           bytes.Buffer
	outputChannel chan string
	writeLock chan int
}

func NewDeployment(path string) *Deployment {
	t := &Deployment{Status: ""}
	t.Path = path
	t.outputChannel = make(chan string, 1)
	t.writeLock = make(chan int, 1)
   	return t

}

func (t *Deployment) readOutput() {
	tempOutput := ""
	for {
		temp := t.buf.String()
		if !strings.Contains(tempOutput, temp) {
			temp = strings.Replace(temp, tempOutput, "", -1)
			tempOutput += temp
			t.outputChannel <- temp
			fmt.Println(temp)
			if strings.Contains(temp, "Finished") {// MUST FIX WHEN TO STOP: SHOULD BE PUT HER
				//t.Output = []byte(tempOutput)
				t.Status = ""
				temp = ""
				return
			}
		}
		
	}

}

func (t *Deployment) getOutput() {
	out := ""
	temp := ""

	for{
		select{
			case out = <-t.outputChannel:
				temp += out
				t.Output = []byte(temp)
				if strings.Contains(temp, "Finished") {
					temp = ""
					out = ""
					return
				}
				out = ""
			case <- t.writeLock:
				<-t.writeLock
				t.Output = []byte("")
				temp = ""
				out = ""
			case <- time.After(30*time.Second):
				return
		}
	}
}

func (t *Deployment) TerraformCommand(command string, path string) {

	t.Status = "Running"
	t.getModules() // SHOULD BE PUT SOMEWHERE ELSE!!
	t.buf.Reset()

	if command == "destroy" {
		cmd := exec.Command("terraform", command, "-force")
		cmd.Dir = path
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		cmd.Start()

		defer cmd.Wait()
		go t.readOutput()
		go t.getOutput()

		io.Copy(&t.buf, stdout)
		io.Copy(&t.buf, stderr)

		t.buf.Write([]byte("\nFinished "))
		return
	}

	cmd := exec.Command("terraform", command)
	cmd.Dir = path
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	defer cmd.Wait()
	go t.readOutput()
	go t.getOutput()

	io.Copy(&t.buf, stdout)
	io.Copy(&t.buf, stderr)
	t.buf.Write([]byte("Finished"))

	//DELETE KEYS?????
}

func (t *Deployment) getModules() {

	init := exec.Command("terraform", "get")
	init.Dir = t.Path
	init.CombinedOutput()
}



func (t *Deployment) GetDeploymentJSON() []byte {

	type deploy struct {
		Status         	string `json:"status"`
		Output			string `json:"output"`

	}
	d := new(deploy)
	if t.Status == "Running"{
		t.writeLock <- 1 //readlock
		d.Status = t.Status
		d.Output = string(t.Output)
		t.writeLock <- 1
	} else {
		d.Status = t.Status
		d.Output = string(t.Output)
	}	


	deploymentJSON, _ := json.Marshal(d)
	fmt.Println(string(deploymentJSON))

	return deploymentJSON
}
