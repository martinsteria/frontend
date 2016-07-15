package users

import (
	"api"
	"doc"
	"fmt"
	"io/ioutil"
	"library"
	"terraform"
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

	lib = user.User{
		library.LibraryModules,
		library.BuildLibrary(library.LibraryModules)}
}

func HandleUserRequests(r api.RequestData) []byte {
	if r.Method == "GET" {
		_, present := r.Query["add"]
		if present {
			users[r.Query["add"]] = user.CreateUser(usersRootDir + "/" + r.Query["add"])
		}

		_, present = r.Query["user"]
		if present {
			user := users[r.Query["user"]]
			_, present = r.Query["get"]
			if present {
				return user.Lib.GetModuleDocumentationJSON(r.Query["get"])
			} else {
				return user.Lib.GetModuleListJSON()
			}
		}

	} else if r.Method == "POST" {
		_, present := r.Query["user"]
		if present {
			user := users[r.Query["user"]]
			_, present = r.Query["module"]
			if present {
				_, present = r.Query["terraform"]
				if present {
					modulePath := users[r.Query["user"]].Dir + "/" + r.Query["module"]
					doc.CreateTFvars(modulePath, r.Body)
					user.Lib.Modules[r.Query["module"]] = doc.ReadVariableValues(modulePath, user.Lib.Modules[r.Query["module"]])
					fmt.Println(user.Lib.Modules[r.Query["module"]])
					fmt.Println(terraform.TerraformCommand(r.Query["terraform"], modulePath))

				}
			}
		}
	}
	return []byte("{}")
}

func HandleLibraryRequests(r api.RequestData) []byte {
	if r.Method == "GET" {
		v, present := r.Query["get"]

		if present {
			return lib.Lib.GetModuleDocumentationJSON(v)
		}

		v, present = r.Query["copy"]
		if present {
			v, present = r.Query["user"]
			if present {
				user := users[r.Query["user"]]
				user.AddModule(library.LibraryModules + "/" + r.Query["copy"])
			}
		}
	}

	return lib.Lib.GetModuleListJSON()
}
