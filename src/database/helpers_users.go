package database

import "strings"

func CreateNewUser(info *UserInfo) error {
	info.UserId = strings.TrimSpace(strings.ToLower(info.UserId))
	if usersInfoMap.Exists(info.UserId) {
		return ErrUserAlreadyExists
	}

	return nil
}
