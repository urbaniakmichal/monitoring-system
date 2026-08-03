//go:build windows

package users

import (
	"os/user"
)

func (*SystemUsers)RetrieveUsersInfo() ([]UserInformation, error) {
	curr, err := user.Current()
	if err != nil {
		return nil, err
	}
	return []UserInformation{
		{
			Username: curr.Username,
			Terminal: "Console",
			Host:     curr.Name,
		},
	}, nil
}