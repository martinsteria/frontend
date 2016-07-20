//Package requests provide functions for handling specific http requests
package requests

import (
	"api"
	"library"
	//"terraform"
	"encoding/json"
	"fmt"
	"users"
	"fmt"
)

var lib *library.Library

//Init initializes the package. Must be called before anything else
func Init(resourcesRootDir string) {
	lib = library.NewLibrary(resourcesRootDir + libraryModulesDir)
	users.Init(resourcesRootDir + UsersRootDir)
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
	return []byte("{}")
}

//HandleDeployRequests handles requests to the deploy endpoint
func HandleDeployRequests(r api.RequestData) []byte {
	if r.Method == "POST" {
		if user, present := r.Query["user"]; present {
			if module, present := r.Query["module"]; present {
				if users.GetDeployStruct(user).Status == "Running" {
					return []byte("{\"status:\": \"Running\"}")
				}
<<<<<<< HEAD
				if command, present := r.Query["command"]; present { 
					deploy := users.GetDeployStruct(user)
					users.GetLibrary(user).Modules[module].UpdateModule(r.Body)
					go deploy.TerraformCommand(command, users.UsersRootDir + "/" + user + "/" + module)
					output := deploy.GetDeploymentJSON()
					fmt.Println()
=======
				if command, present := r.Query["command"]; present { // DO I HAVE TO CHECK FOR BODY??
					users.GetLibrary(user).Modules[module].UpdateModule(r.Body)
					module = module
					go users.GetDeployStruct(user).TerraformCommand(command, users.UsersRootDir+"/"+user+"/"+module)
					users.GetDeployStruct(user).BufferRead <- 1
					output, _ := json.Marshal(users.GetDeployStruct(user))
					users.GetDeployStruct(user).BufferRead <- 1
					<-users.GetDeployStruct(user).Deleted
>>>>>>> f9c7eb2529773a4f450387b85abcca12a1f309ca
					return output
				}
			}
		}
	} else if r.Method == "GET" {
		if user, present := r.Query["user"]; present {
			deploy := users.GetDeployStruct(user)
			fmt.Println(deploy.Status)
			if deploy.Status != "Running" {
				fmt.Println("STATUS: " + string(deploy.Status) + "\noutput: " + string(deploy.Output))
				output, _ := json.Marshal(deploy)
				return output
			} else {
				if _, present := r.Query["module"]; present {
					deploy.BufferRead <- 1
					output, _ := json.Marshal(deploy)
					deploy.BufferRead <- 1 // Bufferen blir fullt
					<-deploy.Deleted
					return output
				}
			}
		}
	}

	return []byte("{\"status:\": \"failed\"}")
}
