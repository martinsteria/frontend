package users

import (
	"api"
	"fmt"
	"io/ioutil"
	"os/exec"
)

const (
	UsersRootFolder = "/home/martin/users"
)

func HandleUserRequests(r api.RequestData) []byte {

	v, present := r.Query["add"]
	if present {
		return CreateUser(v[0])
	}

	return []byte("{}")
}

func CreateUser(user string) []byte {
	fmt.Println("Creating user " + user)
	files, _ := ioutil.ReadDir(UsersRootFolder)
	for _, f := range files {
		if f.Name() == user {
			return []byte("{\"status\": \"alreadyExists\"}")
		}
	}

	exec.Command("mkdir", UsersRootFolder+"/"+user).Output()
	return []byte("{\"status\": \"success\"}")
}
