package users

import (
	"api"
	"io/ioutil"
	"library"
	"user"
)

const (
	usersRootDir = "/home/martin/users"
)

var users map[string]user.User
var lib user.User

func Init() {
	users = make(map[string]user.User)
	files, _ := ioutil.ReadDir(usersRootDir)
	for _, u := range files {
		userPath := usersRootDir + "/" + u.Name()
		users[u.Name()] = user.User{userPath, library.BuildLibrary(userPath)}
	}

	//Bad design. BuildLibrary should take the library root as parameter
	lib = user.User{
		library.LibraryModules,
		library.BuildLibrary(library.LibraryModules)}
}

func HandleUserRequests(r api.RequestData) []byte {
	v, present := r.Query["add"]
	if present {
		users[v[0]] = user.CreateUser(v[0])
	}

	return []byte("{}")
}

func HandleLibraryRequests(r api.RequestData) []byte {
	if r.Method == "GET" {
		v, present := r.Query["get"]

		if present {
			return lib.Lib.GetModuleDocumentationJSON(v[0])
		}

		v, present = r.Query["copy"]
		if present {
			v, present = r.Query["user"]
			if present {
				user := users[r.Query["user"][0]]
				user.AddModule(library.LibraryModules + "/" + r.Query["copy"][0])
			}
		}

	}
	return lib.Lib.GetModuleListJSON()
}
