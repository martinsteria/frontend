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
				if users.GetLibrary(user).Modules[module].Deployment.Status == "Running" {
					return []byte("{\"status:\": \"Running\"}")
				}
				if command, present := r.Query["command"]; present { // DO I HAVE TO CHECK FOR BODY??
					usr := users.GetLibrary(user).Modules[module]
					usr.UpdateModule(r.Body)
					usr.Deployment.Init(users.UsersRootDir + "/" + user + "/" + module)
					go usr.Deployment.TerraformCommand(command)
					usr.Deployment.BufferRead <- 1
					output, _ := json.Marshal(usr.Deployment)
					usr.Deployment.BufferRead <- 1
					<- usr.Deployment.Deleted
					return output
				}
			}
		}
	} else if r.Method == "GET" {
		if user, present := r.Query["user"]; present {
			if module, present := r.Query["module"]; present {
				deploy := users.GetLibrary(user).Modules[module].Deployment
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
