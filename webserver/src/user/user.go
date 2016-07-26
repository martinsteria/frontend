//Package user contains functionality for working with users interacting with a Terraform library
package user

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"module"
	"time"
	"os"
)

//User contains information about a user
type User struct {
	//The directory where the users modules are located
	RootDir string

	//The users modules
	Modules map[string]*module.Module
}

//NewUser creates a new user and initializes it with all its modules. Dir is the users folder
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

//GetModuleListJSON returns a list of the users modules represented as JSON
//It has the following fields
//"name": The name of the module
//"id": A unique identifier for the module. Specifically, the name of the module's folder
//"description": A short description of the modules purpose
//"provider": The provider of the module
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

//GetModule returns the module structure of the argument id. If it doesn't exist, the function returns nil
func (u *User) GetModule(id string) *module.Module {
	if module, present := u.Modules[id]; present {
		return module
	}
	return nil
}

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
	//exec.Command("cp", "-r", m.Path, u.RootDir).Output()
	copyDirectory(m.Path, u.RootDir + "/" + m.Id)

	u.Modules[m.Id] = module.NewModule(u.RootDir + "/" + m.Id)
	return []byte("{\"status\": \"success\"}")
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func copyDirectory(from string, to string) {

	var dir []string
	files, err := ioutil.ReadDir(from)
	checkError(err)

	err = os.MkdirAll(to, os.ModePerm)
	checkError(err)

	checkError(err)
	for _, f := range files {
		checkError(err)
		if(f.IsDir() ){
			dir = append(dir, f.Name())
		} else {
			content, err := ioutil.ReadFile(from +"/"+ f.Name())
			checkError(err)
			os.Chdir(to)
			ioutil.WriteFile(f.Name(), content, os.ModePerm)
		}		
	}
	for d := range dir {
		copyDirectory(from + "/" + dir[d], to + "/" + dir[d])
	}
}
