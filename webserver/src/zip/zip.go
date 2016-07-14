package zip

import (

	"os/exec"
	"io/ioutil"
	"strings"
	"fmt"
)



func Unzip(path string){

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if strings.Contains(f.Name(), ".zip"){
			out, _ := exec.Command("unzip", path + "/" + f.Name()).CombinedOutput()
			fmt.Println(string(out))
		}
	}

}

func Zip(path string){

	out, _ := exec.Command("zip", "-r", "folder{.zip,}").CombinedOutput()
	fmt.Println(string(out))

}