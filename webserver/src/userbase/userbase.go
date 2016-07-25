//package userbase provides functionality for managing a collection of users
package userbase

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"user"
)

type Userbase struct {
	Users   map[string]*user.User
	RootDir string
}

//NewUserbase creates a new userbase and initializes all existing users. The argument path is the root of the userbase directory
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

//AddUser adds a new user to the userbase
func (ub *Userbase) AddUser(name string) []byte {
	rootDir := ub.RootDir + "/" + name
	if _, present := ub.Users[name]; present {
		return []byte("{\"status\": \"User already exists\"}")
	}
	os.MkdirAll(rootDir, os.ModePerm)
	ub.Users[name] = user.NewUser(rootDir)

	return []byte("{\"status\": \"success\"}")
}

//GetUser returns a user from the userbase. If the user does not exist, it returns nil
func (ub *Userbase) GetUser(name string) *user.User {
	if user, present := ub.Users[name]; present {
		return user
	}
	return nil
}

//GetUserListJSON return a list of usernames represented as JSON
func (ub *Userbase) GetUserListJSON() []byte {
	var usernames []string
	for name, _ := range ub.Users {
		usernames = append(usernames, name)
	}
	usernamesJSON, _ := json.Marshal(usernames)
	return usernamesJSON
}
