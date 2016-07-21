//Package requests provide functions for handling specific http requests
package requests

import (
	"api"
	"library"
	//"terraform"
	"users"
)

const (
	usersRootDir      = "/users"
	libraryModulesDir = "/library/modules"
)

var lib *library.Library

//Init initializes the package. Must be called before anything else
func Init(resourcesRootDir string) {
	lib = library.NewLibrary(resourcesRootDir + libraryModulesDir)
	users.Init(resourcesRootDir + usersRootDir)
}

//HandleUserRequests handles requests to the users endpoint
func HandleUserRequests(r api.RequestData) []byte {
	if r.Method == "GET" {
		if user, present := r.Query["user"]; present {
			if module, present := r.Query["module"]; present {
				return users.GetLibrary(user).GetModuleDocumentationJSON(module)
			} else {
				if lib := users.GetLibrary(user); lib != nil {
					return lib.GetModuleListJSON()
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
				if users.GetDeployStruct(user).Status == "Running" {
					return []byte("{\"status:\": \"Running\"}")
				}
				if command, present := r.Query["command"]; present {
					deploy := users.GetDeployStruct(user)
					users.GetLibrary(user).Modules[module].UpdateModule(r.Body)
					go deploy.TerraformCommand(command, lib.GetRootDir()+"/"+user+"/"+module)
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
