package terraform

import (
	//"log"
	"os/exec"
	"io"
	"bytes"
	"fmt"
)

type Deployment struct {
	Status         	string `json:"status"`
	Path         	string `json:"path"`
	Output  		[]byte `json:"description"`
	buf 			bytes.Buffer
}

func (t *Deployment)Init(path string){
	t.Path = path
}

func (t *Deployment) Print(){

		var last bytes.Buffer
		for{
			if(t.buf.String() != last.String()) {
				//t.buf.UnreadByte()
				t.Output = t.buf.Bytes()
				t.buf.Reset()
				last = t.buf
				fmt.Println(string(t.Output))

				if t.Status != "Running" {
					return
				}
			}
		}

} 

func (t *Deployment)TerraformCommand(command string) {


	t.Status = "Running"
	t.Output = []byte("Running")
	t.getModules() // SHOULD BE PUT SOMEWHERE ELSE!!


	if command == "destroy" {
		cmd := exec.Command("terraform", command, "-force")
		cmd.Dir = t.Path
		//out, _ := cmd.CombinedOutput()
    	cmd.Run()

		//t.Output = out
		return
	}

	cmd := exec.Command("terraform", command)
	cmd.Dir = t.Path

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	defer cmd.Wait()


	io.Copy(&t.buf, stdout)

	//out, err := cmd.CombinedOutput()

	/*if err != nil {
		log.Fatal(err)
	}*/
	t.Status = ""
	//t.Output = out

	//DELETE KEYS?????
}

func (t *Deployment)getModules() {

	init := exec.Command("terraform", "get")
	init.Dir = t.Path
	init.CombinedOutput()
}


