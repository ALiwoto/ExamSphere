package appValues

import "strings"

func NormalizeUserId(userId string) string {
	return strings.ToLower(strings.TrimSpace(userId))
}

func ParseRole(value string) UserRole {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "owner":
		return UserRoleOwner
	case "admin":
		return UserRoleAdmin
	case "student":
		return UserRoleStudent
	case "teacher":
		return UserRoleTeacher
	default:
		return UserRoleUnknown
	}
}

func IsUserIdValid(userId string) bool {
	userId = NormalizeUserId(userId)
	if strings.Contains(userId, " ") {
		return false
	}

	return len(userId) >= MinUserIdLength &&
		len(userId) <= MaxUserIdLength
}
