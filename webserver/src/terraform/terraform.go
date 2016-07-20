package terraform

import (     
	//"log"     
	"os/exec"    
	"io"      
	"bytes"     
	"strings"
	"time" 
	)

type Deployment struct {
	Status         	string `json:"status"`
	Path         	string `json:"path"`
	Output  		[]byte `json:"output"`
	buf 			bytes.Buffer
	outputChannel chan string
	BufferRead chan int
	Deleted chan int
}

func NewDeployment(path string) *Deployment {
    t := &Deployment{Status : ""}
	t.Path = path
	t.outputChannel = make(chan string, 1)
	t.BufferRead = make(chan int, 1)
	t.Deleted = make(chan int, 1)
   
   return t
}
	
func (t *Deployment) readOutput(){
	tempOutput := ""
	for{
		temp := t.buf.String()
		if !strings.Contains(tempOutput, temp) {
			temp = strings.Replace(temp, tempOutput, "", -1)
			tempOutput += temp
			t.outputChannel <- temp
		
		if strings.Contains(temp, "Finished") {// MUST FIX WHEN TO STOP: SHOULD BE PUT HER
			//t.Output = []byte(tempOutput)
			t.Status = ""
			return
		}
		}
	}

}

func (t *Deployment) getOutput(){
	out := ""
	temp := ""
	for{
		select{
			case out = <-t.outputChannel:
				temp += out
				t.Output = []byte(temp)
				if strings.Contains(temp, "Finished") {
					return
				}
				out = ""
			case <- t.BufferRead:
				<-t.BufferRead
				t.Output = []byte("")
				temp = ""
				t.Deleted <- 1
			case <- time.After(30*time.Second):
				return

		}

	}
}

func (t *Deployment)TerraformCommand(command string, path string) {

	t.Status = "Running"
	fmt.Println("I TerraformCommand")

	t.getModules() // SHOULD BE PUT SOMEWHERE ELSE!!


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

		t.buf.Write([]byte("\nFinished"))
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


