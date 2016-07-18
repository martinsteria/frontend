package requests

import (
	"api"
	"library"
	//"terraform"
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
				return users.GetLibrary(user).GetModuleListJSON()
			}
		}

	} else if r.Method == "POST" {
		if user, present := r.Query["user"]; present {
			users.AddUser(user)
			return []byte("{\"status\": \"success\"}")
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
				users.AddModule(user, lib.GetRootDir()+"/"+module)
				return []byte("{\"status:\": \"success\"}")
			}
		}
	}
	return []byte("{}")
}

func HandleDeployRequests(r api.RequestData) []byte {
	if r.Method == "POST" {
		if user, present := r.Query["user"]; present {
			if module, present := r.Query["module"]; present {
				if command, present := r.Query["command"]; present { // DO I HAVE TO CHECK FOR BODY??
					user.Lib.Modules[module].UpdateModule(r.Body)
					user.Lib.Modules[module].Deployment.Init(usersRootDir + "/" + user + "/" + module)
					user.Lib.Modules[module].Deployment.TerraformCommand(command)
					return user.Lib.Modules[module].Deployment.Output
				}

			}
		}
	}

	return []byte("failed")
}
