package userHandlers

import (
	"sync"
	"time"
)

type LoginData struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

type LoginResult struct {
	UserId   string `json:"user_id"`
	FullName string `json:"full_name"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"expiration"`
	Role         string `json:"role"`
}

type AuthResult struct {
	UserId   string `json:"user_id"`
	FullName string `json:"full_name"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"expiration"`
	Role         string `json:"role"`
}

type MeResult struct {
	UserId   string `json:"user_id"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type userRequestEntry struct {
	RequestPath string
	LastTryAt   time.Time
	TryCount    int
	mut         *sync.Mutex
}

type ChangePasswordData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
