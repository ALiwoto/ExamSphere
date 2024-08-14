package userHandlers

import (
	"ExamSphere/src/core/appValues"
	"ExamSphere/src/database"
	"sync"
	"time"
)

type LoginData struct {
	UserId        string `json:"user_id"`
	Password      string `json:"password"`
	ClientRId     string `json:"client_rid"`
	CaptchaId     string `json:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer"`
} // @name LoginData

type LoginResult struct {
	UserId   string `json:"user_id"`
	FullName string `json:"full_name"`

	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	Expiration   int64              `json:"expiration"`
	Role         appValues.UserRole `json:"role"`
} // @name LoginResult

type AuthResult struct {
	UserId   string `json:"user_id"`
	FullName string `json:"full_name"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"expiration"`
	Role         string `json:"role"`
} // @name AuthResult

type MeResult struct {
	UserId   string             `json:"user_id"`
	FullName string             `json:"full_name"`
	Role     appValues.UserRole `json:"role"`
} // @name GetMeResult

type CreateUserData = database.NewUserData

// CreateUserResult is the result of creating a new user.
type CreateUserResult struct {
	UserId   string `json:"user_id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
} // @name CreateUserResult

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
