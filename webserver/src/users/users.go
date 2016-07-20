//Package users contains functionality for working with users interacting with a Terraform library
package users

import (
	"encoding/json"
	"io/ioutil"
	"library"
	"os"
	"os/exec"
	"terraform"
	"strings"
)

const (
	UsersRootDir = "/home/martin/users"
)

//User contains information about a user
type User struct {
	//The directory where the users modules are located
	RootDir string

	//The users private library
	Lib *library.Library

	//The status of the users deployment
	Deploy *terraform.Deployment
}

var users map[string]*User
var usersRootDir string

//Init initializes the existing userbase
func Init(usersRootDir string) {
	usersRootDir = usersRootDir
	users = make(map[string]*User)
	files, _ := ioutil.ReadDir(UsersRootDir)
	for _, u := range files {
		userPath := UsersRootDir + "/" + u.Name()
		users[u.Name()] = &User{RootDir: userPath}
		users[u.Name()].Lib = library.NewLibrary(userPath)
		users[u.Name()].Deploy = terraform.NewDeployment(UsersRootDir)
	}

}

//AddUser adds a new user to the userbase. If the user already exists an error code is returned
//Returns an array of bytes representing JSON
func AddUser(name string) []byte {
	rootDir := UsersRootDir + "/" + name
	if _, err := os.Stat(rootDir); err == nil {
		return []byte("{\"status\": \"User already exists\"}")
	}

	exec.Command("mkdir", rootDir).Output()
	users[name] = &User{RootDir: rootDir}
	users[name].Lib = library.NewLibrary(rootDir)
	users[name].Lib.Build()
	users[name].Deploy = terraform.NewDeployment(UsersRootDir)

	return []byte("{\"status\": \"success\"}")
}

//AddModule copies a module from the main library to the user
//If a user already has a copy, an error is returned
func AddModule(user string, modulePath string) []byte {
	if _, present := users[user]; !present {
		return []byte("{\"status\": \"User not found\"}")
	}

	if _, err := os.Stat(modulePath); err != nil {
		if os.IsNotExist(err) {
			return []byte("{\"status\": \"Module not found\"}")
		}
	}
	splitted := strings.Split(modulePath, "/")
	module := splitted[len(splitted)]

	if _, err := os.Stat(user + "/" + module); err == nil {
		return []byte("\"status\": \"User module already exists\"")
	}

	exec.Command("cp", "-r", modulePath, users[user].RootDir).Output()
	users[user].Lib.Build()
	return []byte("{\"status\": \"success\"}")
}

//GetLibrary returns a users library
func GetLibrary(user string) *library.Library {
	if users[user] == nil {
		return nil
	}
	return users[user].Lib
}

func GetDeployStruct(user string) *terraform.Deployment {
	if users[user] == nil {
		return nil
	}
	return users[user].Deploy
}

//GetUserListJSON returns a JSON representation of all the users
func GetUserListJSON() []byte {
	var usernames []string
	for name, _ := range users {
		usernames = append(usernames, name)
	}
	usernamesJSON, _ := json.Marshal(usernames)
	return usernamesJSON
}
