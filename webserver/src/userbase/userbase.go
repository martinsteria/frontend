package userbase

import (
	"user"
	"io/ioutil"
	"log"
)

type Userbase struct {
	Users map[string]*User
	RootDir string
}

func NewUserbase(path string) *Userbase {
	ub := &Userbase{RootDir: path}
	users = make(map[string]*User)
	files, err := ioutil.ReadDir(usersRootDir)

	if err != nil {
		log.Fatal(err)
	}

	for _, u := range files {
		userPath := ub.RootDir + "/" + u.Name()
		users[u.Name()] = newUser(userPath)
	}
}

func (ub *Userbase) AddUser(name string) []byte {
	rootDir := ub.RootDir + "/" + name
	if _, present := Users[name]; present {
		return []byte("{\"status\": \"User already exists\"}")
	}

	exec.Command("mkdir", rootDir).Output()
	ub.Users[name] = newUser(rootDir)

	return []byte("{\"status\": \"success\"}")
}

func (ub *Userbase) GetUser(name string) *User {
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
