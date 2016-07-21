//Package requests provide functions for handling specific http requests
package requests

import (
	"api"
	"library"
	//"terraform"
	"userbase"
	"user"
)

const (
	usersRootDir      = "/users"
	libraryModulesDir = "/library/modules"
)

var library *user.User
var userbase *userbase.Userbase

//Init initializes the package. Must be called before anything else
func Init(resourcesRootDir string) {
	userbase = userbase.NewUserbase(resourcesRootDir + usersRootDir)
	library = user.NewUser(resourcesRootDir + libraryModulesDir)
	users.Init()
}

//HandleUserRequests handles requests to the users endpoint
func HandleUserRequests(r api.RequestData) []byte {
	if r.Method == "GET" {
		if username, present := r.Query["user"]; present {
			if moduleId, present := r.Query["module"]; present {
				if user := userbase.GetUser(username); user != nil {
					if module := user.GetModule(moduleId); module != nil {
						return module.GetDocumentationJSON()
					}
				}
			} else {
				if user := userbase.GetUser(username); user != nil {
					return user.GetModuleListJSON()
				}
				return []byte("{\"status\": \"User not found\"}")
			}
		}

	} else if r.Method == "POST" {
		if user, present := r.Query["user"]; present {
			return usersbase.AddUser(user)
		}
	}

	return userbase.GetUserListJSON()
}

//HandleLibraryRequests handles requests to the library endpoint
func HandleLibraryRequests(r api.RequestData) []byte {
	if r.Method == "GET" {
		if module, present := r.Query["module"]; present {
			return lib.GetModuleDocumentationJSON(module)
		}
	}

	return lib.GetModuleListJSON()
}

//HandleLibraryCopyRequests handles requests to the library/copy endpoint
func HandleLibraryCopyRequests(r api.RequestData) []byte {
	if r.Method == "POST" {
		if user, present := r.Query["user"]; present {
			if module, present := r.Query["module"]; present {
				return users.AddModule(user, lib.GetRootDir()+"/"+module)
			}
		}
	}
	return []byte("{\"status:\": \"Invalid request\"}")
}

func HandleDeployRequests(r api.RequestData) []byte {
	if r.Method == "POST" {
		if user, present := r.Query["user"]; present {
			if module, present := r.Query["module"]; present {
				if command, present := r.Query["command"]; present {
					if users.GetDeployStruct(user).Status == "Running" {
						return []byte("{\"status:\": \"Running\"}")
					}
					users.GetUser(user).GetModule(module).Deploy(command)
					deploy := users.GetDeployStruct(user)
					users.GetLibrary(user).Modules[module].UpdateModule(r.Body)
					go deploy.TerraformCommand(command, users.GetLibrary(user).Modules[module].Path)
					output := deploy.GetDeploymentJSON()
					return output
				}
			}
		}
	} else if r.Method == "GET" {
		if user, present := r.Query["user"]; present {
			deploy := users.GetDeployStruct(user)
			if deploy.Status != "Running" {
				out := deploy.GetDeploymentJSON()
				deploy.Output = []byte("")
				return out
			} else {
				if _, present := r.Query["module"]; present {
					output := deploy.GetDeploymentJSON()
					return output
				}
			}
		}
	}

	return []byte("{\"status:\": \"failed\"}")
}
