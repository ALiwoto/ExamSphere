package database

import (
	"OnlineExams/src/core/appValues"
	"time"
)

type UserInfo struct {
	// UserId is the user-id of the user.
	UserId string `json:"user_id"`

	// AuthHash the hash used to invalidate when the user changes their
	// password. User CANNOT login only with this hash; this is just
	// a security measure and has no other usage.
	AuthHash string `json:"auth_hash"`

	// FullName is the user's full name.
	FullName string `json:"full_name"`

	// Email is the email of the user.
	Email string `json:"email"`

	// Password is the hashed password of the user.
	Password string `json:"password"`

	// Role is the role of the user.
	Role appValues.UserRole `json:"role"`

	// IsBanned is true if and only if the user is banned.
	IsBanned bool `json:"is_banned"`

	// BanReason is the user's ban reason (from the platform).
	BanReason string `json:"ban_reason"`

	// CreatedAt is the time when the user was created.
	CreatedAt time.Time `json:"created_at"`
}
