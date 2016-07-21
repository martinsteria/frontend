//Package users contains functionality for working with users interacting with a Terraform library
package user

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"module"
	"os/exec"
	"time"
)

//User contains information about a user
type User struct {
	//The directory where the users modules are located
	RootDir string

	//The users modules
	Modules map[string]*module.Module
}

var usersRootDir string

func NewUser(dir string) *User {
	u := &User{RootDir: dir}

	u.Modules = make(map[string]*module.Module)
	files, err := ioutil.ReadDir(u.RootDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		u.Modules[f.Name()] = module.NewModule(u.RootDir + "/" + f.Name())
	}

	return u
}

func (u *User) GetModuleListJSON() []byte {
	type module struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Provider    string `json:"provider"`
	}
	var ms []module

	for _, m := range u.Modules {
		ms = append(ms, module{m.Id, m.Name, m.Description, m.Provider})
	}

	msJSON, _ := json.Marshal(ms)
	return msJSON
}

func (u *User) GetModule(id string) *module.Module {
	if module, present := u.Modules[id]; present {
		return module
	}
	return nil
}

//AddUser adds a new user to the userbase. If the user already exists an error code is returned
//Returns an array of bytes representing JSON

//AddModule copies a module from the main library to the user
//If a user already has a copy, an error is returned
func (u *User) AddModule(m *module.Module) []byte {
	if _, present := u.Modules[m.Id]; present {
		return []byte("\"status\": \"User module already exists\"")
	}

	m.GetCommand().Launch("get")
	for m.GetCommand().IsRunning() {
		time.Sleep(10 * time.Millisecond)
	}
	exec.Command("cp", "-r", m.Path, u.RootDir).Output()
	u.Modules[m.Id] = module.NewModule(u.RootDir + "/" + m.Id)
	return []byte("{\"status\": \"success\"}")
}
