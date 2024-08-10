package appValues

const (
	UserRoleOwner   UserRole = "owner"
	UserRoleAdmin   UserRole = "admin"
	UserRoleStudent UserRole = "student"
	UserRoleTeacher UserRole = "teacher"
	UserRoleUnknown UserRole = ""
)

const (
	MinUserIdLength = 2
	MaxUserIdLength = 16
)
