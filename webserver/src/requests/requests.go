//Package requests provide functions for handling specific http requests
package requests

import (
	"api"
	//"terraform"
	"user"
	"userbase"
)

const (
	usersRootDir      = "/users"
	libraryModulesDir = "/library/modules"
)

var library *user.User
var users *userbase.Userbase

//Init initializes the package. Must be called before anything else
func Init(resourcesRootDir string) {
	users = userbase.NewUserbase(resourcesRootDir + usersRootDir)
	library = user.NewUser(resourcesRootDir + libraryModulesDir)
}

//HandleUserRequests handles requests to the users endpoint
func HandleUserRequests(r api.RequestData) []byte {
	if r.Method == "GET" {
		if username, present := r.Query["user"]; present {
			if moduleId, present := r.Query["module"]; present {
				if user := users.GetUser(username); user != nil {
					if module := user.GetModule(moduleId); module != nil {
						return module.GetDocumentationJSON()
					}
				}
			} else {
				if user := users.GetUser(username); user != nil {
					return user.GetModuleListJSON()
				}
				return []byte("{\"status\": \"User not found\"}")
			}
		}

	} else if r.Method == "POST" {
		if user, present := r.Query["user"]; present {
			return users.AddUser(user)
		}
	}

	return users.GetUserListJSON()
}

//HandleLibraryRequests handles requests to the library endpoint
func HandleLibraryRequests(r api.RequestData) []byte {
	if r.Method == "GET" {
		if moduleId, present := r.Query["module"]; present {
			if module := library.GetModule(moduleId); module != nil {
				return module.GetDocumentationJSON()
			}
		}
	}

	return library.GetModuleListJSON()
}

//HandleLibraryCopyRequests handles requests to the library/copy endpoint
func HandleLibraryCopyRequests(r api.RequestData) []byte {
	if r.Method == "POST" {
		if username, present := r.Query["user"]; present {
			if moduleId, present := r.Query["module"]; present {
				if user := users.GetUser(username); user != nil {
					if module := library.GetModule(moduleId); module != nil {
						return user.AddModule(module)
					}
				}
			}
		}
	}
	return []byte("{\"status:\": \"Invalid request\"}")
}
