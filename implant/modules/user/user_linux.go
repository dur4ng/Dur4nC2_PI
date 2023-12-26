package user

import (
	"github.com/opencontainers/runc/libcontainer/user"
)

func GetWhoami() string {
	currentUser, _ := user.CurrentUser()
	//group, _ := currentUser.CurrentGroup()

	return currentUser.Name
}
