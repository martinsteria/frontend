package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"library"
	"os"
	"os/exec"
)

const (
	usersRootDir = "/home/martin/users"
)

type User struct {
	RootDir string
	Lib     *library.Library
}

var users map[string]*User

func Init() {
	users = make(map[string]*User)
	files, _ := ioutil.ReadDir(usersRootDir)
	for _, u := range files {
		userPath := usersRootDir + "/" + u.Name()
		users[u.Name()] = &User{RootDir: userPath}
		users[u.Name()].Lib = library.NewLibrary(userPath)
	}
}

func AddUser(name string) {
	rootDir := usersRootDir + "/" + name
	if _, err := os.Stat(rootDir); err == nil {
		return
	}

	exec.Command("mkdir", rootDir).Output()
	users[name] = &User{RootDir: rootDir}
	users[name].Lib.Build()
}

func AddModule(user string, modulePath string) {
	fmt.Println("MOD:", modulePath)
	fmt.Println("USER:", users[user].RootDir)
	exec.Command("cp", "-r", modulePath, users[user].RootDir).Output()
	users[user].Lib.Build()
}

func GetLibrary(user string) *library.Library {
	return users[user].Lib
}

func GetUserListJSON() []byte {
	var usernames []string
	for name, _ := range users {
		usernames = append(usernames, name)
	}
	usernamesJSON, _ := json.Marshal(usernames)
	return usernamesJSON
}
