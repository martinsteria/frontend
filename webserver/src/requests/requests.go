package requests

import (
	"api"
	"library"
	//"terraform"
	"encoding/json"
	"users"
)

var lib *library.Library

func Init() {
	lib = library.NewLibrary(library.LibraryModules)
	users.Init()
}

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

func HandleLibraryRequests(r api.RequestData) []byte {
	if r.Method == "GET" {
		if module, present := r.Query["module"]; present {
			return lib.GetModuleDocumentationJSON(module)
		}
	}

	return lib.GetModuleListJSON()
}

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

func HandleDeployRequests(r api.RequestData) []byte {
	if r.Method == "POST" {
		if user, present := r.Query["user"]; present {
			if module, present := r.Query["module"]; present {
				if users.GetDeployStruct(user).Deployment.Status == "Running" {
					return []byte("{\"status:\": \"Running\"}")
				}
				if command, present := r.Query["command"]; present { // DO I HAVE TO CHECK FOR BODY??
					users.GetLibrary(user).Modules[module].UpdateModule(r.Body)
					go users.GetDeployStruct(user).TerraformCommand(command, users.UsersRootDir + "/" + user + "/" + module )
					users.GetDeployStruct(user).BufferRead <- 1
					output, _ := json.Marshal(users.GetDeployStruct(user))
					users.GetDeployStruct(user).BufferRead <- 1
					<- users.GetDeployStruct(user).Deleted
					return output
				}
			}
		}
	} else if r.Method == "GET" {
		if user, present := r.Query["user"]; present {
			if _, present := r.Query["module"]; present {
				deploy := users.GetDeployStruct(user)
				deploy.BufferRead <- 1
				output, _ := json.Marshal(deploy)
				deploy.BufferRead <- 1
				<- deploy.Deleted
				return output
			}
		}
	}

	return []byte("{\"status:\": \"failed\"}")
}
