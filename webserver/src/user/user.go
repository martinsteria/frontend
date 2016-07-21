//Package users contains functionality for working with users interacting with a Terraform library
package user

import (
	"encoding/json"
	"io/ioutil"
	"library"
	"log"
	"os"
	"os/exec"
	"strings"
	"terraform"
	"module"
)

//User contains information about a user
type User struct {
	//The directory where the users modules are located
	RootDir string

	//The status of the users deployment
	Deploy *terraform.Deployment

	//The users modules
	Modules map[string]*module.Module
}

var usersRootDir string

func newUser(dir string) *User {
	u := &User{RootDir: userPath}

	u.Modules = make(map[string]*doc.Module)
	files, err := ioutil.ReadDir(u.RootDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		u.Modules[f.Name()] = doc.NewModule(u.RootDir + "/" + f.Name())
	}

	u.Deploy = terraform.NewDeployment(usersRootDir)

	return u
}

func (u *User) GetModuleListJSON() []byte {
	type module struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	var ms []module

	for _, m := range u.Modules {
		ms = append(ms, module{m.Id, m.Name, m.Description})
	}

	msJSON, _ := json.Marshal(ms)
	return msJSON
}

func (u *User) GetModule(id string) *Module {
	if module, present := u.Modules[id]; present {
		return module
	}
	return nil
}

//AddUser adds a new user to the userbase. If the user already exists an error code is returned
//Returns an array of bytes representing JSON

//AddModule copies a module from the main library to the user
//If a user already has a copy, an error is returned
func (u *User) AddModule(modulePath string) []byte {
	if _, err := os.Stat(modulePath); err != nil {
		if os.IsNotExist(err) {
			return []byte("{\"status\": \"Module not found\"}")
		}
	}

	splitted := strings.Split(modulePath, "/")
	module := splitted[len(splitted)-1]

	if _, err := os.Stat(u.RootDir + "/" + module); err == nil {
		return []byte("\"status\": \"User module already exists\"")
	}

	exec.Command("cp", "-r", modulePath, u.RootDir).Output()
	u.Modules[f.Name()] = doc.NewModule(u.RootDir + "/" + f.Name())
	return []byte("{\"status\": \"success\"}")
}
