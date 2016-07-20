package users

import (
	"encoding/json"
	"io/ioutil"
	"library"
	"os"
	"os/exec"
	"terraform"
)

const (
	UsersRootDir = "D:/Users/bengelse/test/users"
)

type User struct {
	RootDir string
	Lib     *library.Library
	Deploy *terraform.Deployment
}

var users map[string]*User

func Init() {
	users = make(map[string]*User)
	files, _ := ioutil.ReadDir(UsersRootDir)
	for _, u := range files {
		userPath := UsersRootDir + "/" + u.Name()
		users[u.Name()] = &User{RootDir: userPath}
		users[u.Name()].Lib = library.NewLibrary(userPath)
		users[u.Name()].Deploy = terraform.NewDeployment(UsersRootDir)
	}

}

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

func AddModule(user string, modulePath string) []byte {
	if _, present := users[user]; !present {
		return []byte("{\"status\": \"User not found\"}")
	}

	if _, err := os.Stat(modulePath); err != nil {
		if os.IsNotExist(err) {
			return []byte("{\"status\": \"Module not found\"}")
		}
	}

	exec.Command("cp", "-r", modulePath, users[user].RootDir).Output()
	users[user].Lib.Build()
	return []byte("{\"status\": \"success\"}")
}

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

func GetUserListJSON() []byte {
	var usernames []string
	for name, _ := range users {
		usernames = append(usernames, name)
	}
	usernamesJSON, _ := json.Marshal(usernames)
	return usernamesJSON
}
