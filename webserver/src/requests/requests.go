package requests

import (
	"api"
	"library"
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
				if users.GetDeployStruct(user).Status == "Running" {
					return []byte("{\"status:\": \"Running\"}")
				}
				if command, present := r.Query["command"]; present { 
					deploy := users.GetDeployStruct(user)
					//users.GetLibrary(user).Modules[module].UpdateModule(r.Body)
					go deploy.TerraformCommand(command, users.UsersRootDir + "/" + user + "/" + module)
					output := deploy.GetDeploymentJSON()
					return output
				}
			}
		}
	} else if r.Method == "GET" {
		if user, present := r.Query["user"]; present {
			deploy := users.GetDeployStruct(user)
			if deploy.Status != "Running"{
				out := deploy.GetDeploymentJSON()
				deploy.Output = []byte("")
				return out
			} else{
				if _, present := r.Query["module"]; present {
					output := deploy.GetDeploymentJSON()
					return output
				}
			}
		}
	}

	return []byte("{\"status:\": \"failed\"}")
}
