package database

import (
	"OnlineExams/src/core/appConfig"
	"OnlineExams/src/core/appValues"
	"OnlineExams/src/core/utils/hashing"
	"OnlineExams/src/core/utils/logging"
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
)

func GetUserInfoByAuthHash(userId, authHash string) *UserInfo {
	userId = appValues.NormalizeUserId(userId)
	if userId == "" || authHash == "" {
		return nil
	}

	info := usersInfoMap.Get(userId)
	if info == valueUserNotFound {
		return nil
	} else if info == nil {
		info, err := getUserFromDB(userId)
		if err != nil {
			if err == ErrUserNotFound {
				usersInfoMap.Add(userId, valueUserNotFound)
			}
			return nil
		}

		usersInfoMap.Add(userId, info)
	}

	if info.AuthHash != authHash {
		return nil
	}

	return info
}

func GetUserInfoByPass(userId, password string) *UserInfo {
	userId = appValues.NormalizeUserId(userId)
	if userId == "" || password == "" {
		return nil
	} else if appConfig.IsOwnerUsername(userId) {
		if !appConfig.IsOwner(userId, password) {
			return nil
		}

		user := usersInfoMap.Get(userId)
		if user != nil && user.UserId == userId {
			return user
		}

		// add an artificial user
		user = &UserInfo{
			UserId:   userId,
			Password: password,
			FullName: "Administrator",
			Role:     appValues.UserRoleOwner,
			AuthHash: hashing.GenerateAuthHash(),
		}
		usersInfoMap.Add(userId, user)

		return user
	}

	info := usersInfoMap.Get(userId)
	if info == valueUserNotFound {
		return nil
	} else if info == nil {
		info, err := getUserFromDB(userId)
		if err != nil {
			if err == ErrUserNotFound {
				usersInfoMap.Add(userId, valueUserNotFound)
			}
			return nil
		}

		usersInfoMap.Add(userId, info)
	}

	if !hashing.CheckPasswordHash(password, info.Password) {
		return nil
	}

	return info
}

func getUserFromDB(userId string) (*UserInfo, error) {
	info := &UserInfo{}
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT user_id, full_name, email, auth_hash, password, role
		FROM user_info WHERE user_id = $1`,
		userId,
	).Scan(
		&info.UserId, &info.FullName, &info.Email,
		&info.AuthHash, &info.Password, &info.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return info, nil
}

func CreateNewUser(data *NewUserData) error {
	data.UserId = strings.TrimSpace(strings.ToLower(data.UserId))
	if usersInfoMap.Exists(data.UserId) {
		return ErrUserAlreadyExists
	}

	info := &UserInfo{
		UserId:   data.UserId,
		FullName: data.FullName,
		Email:    data.Email,
		AuthHash: hashing.GenerateAuthHash(),
		Password: hashing.HashPassword(data.RawPassword),
		Role:     data.Role,
	}

	var newUserId string
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT create_user_info(
			p_user_id := $1,
			p_full_name := $2,
			p_email := $3,
			p_auth_hash := $4,
			p_password := $5,
			p_role := $6
		)`,
		info.UserId,
		info.FullName,
		info.Email,
		info.AuthHash,
		info.Password,
		info.Role,
	).Scan(&newUserId)
	if err != nil {
		return err
	} else if newUserId != info.UserId {
		logging.Error("CreateNewUser: failed to create user: ", newUserId, " != ", info.UserId)
		return ErrInternalDatabaseError
	}

	usersInfoMap.Add(newUserId, info)
	return nil
}
