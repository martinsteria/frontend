package userbase

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
	"user"
)

type Userbase struct {
	Users   map[string]*user.User
	RootDir string
}

func NewUserbase(path string) *Userbase {
	ub := &Userbase{RootDir: path}
	ub.Users = make(map[string]*user.User)
	files, err := ioutil.ReadDir(ub.RootDir)

	if err != nil {
		log.Fatal(err)
	}

	for _, u := range files {
		userPath := ub.RootDir + "/" + u.Name()
		ub.Users[u.Name()] = user.NewUser(userPath)
	}

	return ub
}

func (ub *Userbase) AddUser(name string) []byte {
	rootDir := ub.RootDir + "/" + name
	if _, present := ub.Users[name]; present {
		return []byte("{\"status\": \"User already exists\"}")
	}

	exec.Command("mkdir", rootDir).Output()
	ub.Users[name] = user.NewUser(rootDir)

	return []byte("{\"status\": \"success\"}")
}

func (ub *Userbase) GetUser(name string) *user.User {
	if user, present := ub.Users[name]; present {
		return user
	}
	return nil
}

func (ub *Userbase) GetUserListJSON() []byte {
	var usernames []string
	for name, _ := range ub.Users {
		usernames = append(usernames, name)
	}
	usernamesJSON, _ := json.Marshal(usernames)
	return usernamesJSON
}
