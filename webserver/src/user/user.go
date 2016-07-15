package user

import (
	"library"
	"os"
	"os/exec"
)

type User struct {
	Dir string
	Lib library.Library
}

func CreateUser(dirPath string) User {
	if _, err := os.Stat(dirPath); err == nil {
		return User{}
	}

	exec.Command("mkdir", dirPath).Output()
	return User{Dir: dirPath}
}

func (u *User) AddModule(modulePath string) {
	exec.Command("cp", "-r", modulePath, u.Dir).Output()
	u.Lib = library.BuildLibrary(u.Dir)
}
