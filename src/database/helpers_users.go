package database

import (
	"ExamSphere/src/core/appConfig"
	"ExamSphere/src/core/appValues"
	"ExamSphere/src/core/utils/hashing"
	"ExamSphere/src/core/utils/logging"
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
	var err error
	if info == valueUserNotFound {
		return nil
	} else if info == nil {
		info, err = getUserFromDB(userId)
		if err != nil {
			if err == ErrUserNotFound {
				usersInfoMap.Add(userId, valueUserNotFound)
			}
			return nil
		}

		usersInfoMap.Add(userId, info)
	}

	if info == nil || info.AuthHash != authHash {
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

		return createArtificialOwnerUser(userId)
	}

	info := usersInfoMap.Get(userId)
	var err error
	if info == valueUserNotFound {
		return nil
	} else if info == nil {
		info, err = getUserFromDB(userId)
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

// GetUserByUserId returns a user by their user-id.
// If the user is not found (e.g. the id does not exist), the error will always be
// ErrUserNotFound; anything other than that should be treated as an internal error.
func GetUserByUserId(userId string) (*UserInfo, error) {
	userId = appValues.NormalizeUserId(userId)
	if userId == "" {
		return nil, ErrUserNotFound
	}

	info := usersInfoMap.Get(userId)
	var err error
	if info == valueUserNotFound {
		return nil, ErrUserNotFound
	} else if info == nil {
		if appConfig.IsOwnerUsername(userId) {
			return createArtificialOwnerUser(userId), nil
		}

		info, err = getUserFromDB(userId)
		if err != nil {
			if err == ErrUserNotFound {
				usersInfoMap.Add(userId, valueUserNotFound)
			}

			return nil, err
		}

		if info == nil || info.UserId != userId {
			return nil, ErrInternalDatabaseError
		}

		usersInfoMap.Add(userId, info)
	}

	return info, nil
}

func createArtificialOwnerUser(userId string) *UserInfo {
	info := &UserInfo{
		UserId:   userId,
		FullName: "Administrator",
		Role:     appValues.UserRoleOwner,
		AuthHash: hashing.GenerateAuthHash(),
	}

	usersInfoMap.Add(userId, info)
	return info
}

func getUserFromDB(userId string) (*UserInfo, error) {
	info := &UserInfo{}
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT user_id, full_name, email, auth_hash, password, role, is_banned, ban_reason, created_at
		FROM user_info WHERE user_id = $1`,
		userId,
	).Scan(
		&info.UserId,
		&info.FullName,
		&info.Email,
		&info.AuthHash,
		&info.Password,
		&info.Role,
		&info.IsBanned,
		&info.BanReason,
		&info.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return info, nil
}

func CreateNewUser(data *NewUserData) (*UserInfo, error) {
	data.UserId = strings.TrimSpace(strings.ToLower(data.UserId))
	info := usersInfoMap.Get(data.UserId)
	if info != nil && info != valueUserNotFound && info.UserId == data.UserId {
		return nil, ErrUserAlreadyExists
	}

	info = &UserInfo{
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
		return nil, err
	} else if newUserId != info.UserId {
		logging.Error("CreateNewUser: failed to create user: ", newUserId, " != ", info.UserId)
		return nil, ErrInternalDatabaseError
	}

	usersInfoMap.Add(newUserId, info)
	return info, nil
}

// SearchUser searches for users based on the query.
func SearchUser(searchData *SearchUserData) ([]*UserInfo, error) {
	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT user_id, full_name, email, role, is_banned, ban_reason, created_at
		FROM user_info
		WHERE user_id ILIKE $1 OR full_name ILIKE $1 OR email ILIKE $1
		ORDER BY user_id ASC
		LIMIT $2 OFFSET $3`,
		"%"+searchData.Query+"%",
		searchData.Limit,
		searchData.Offset,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []*UserInfo
	for rows.Next() {
		info := &UserInfo{}
		err = rows.Scan(
			&info.UserId,
			&info.FullName,
			&info.Email,
			&info.Role,
			&info.IsBanned,
			&info.BanReason,
			&info.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, info)
	}

	return users, nil
}

// UpdateUserInfo updates the user's basic information.
func UpdateUserInfo(data *UpdateUserData) (*UserInfo, error) {
	data.UserId = appValues.NormalizeUserId(data.UserId)
	if data.UserId == "" {
		return nil, ErrUserNotFound
	}

	info, err := GetUserByUserId(data.UserId)
	if err != nil {
		return nil, err
	}

	if data.FullName != "" {
		info.FullName = data.FullName
	}

	if data.Email != "" {
		info.Email = data.Email
	}

	_, err = DefaultContainer.db.Exec(context.Background(),
		`UPDATE user_info
		SET full_name = $1, email = $2
		WHERE user_id = $3`,
		info.FullName,
		info.Email,
		info.UserId,
	)
	if err != nil {
		return nil, err
	}

	return info, nil
}
