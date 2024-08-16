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

// CreateUserData is the data required to create a new user.
type CreateUserData = database.NewUserData // @name CreateUserData

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

type SearchUserData = database.SearchUserData // @name SearchUserData

type SearchUserResult struct {
	Users []SearchedUserInfo `json:"users"`
} // @name SearchUserResult

type SearchedUserInfo struct {
	UserId    string             `json:"user_id"`
	FullName  string             `json:"full_name"`
	Role      appValues.UserRole `json:"role"`
	Email     string             `json:"email"`
	IsBanned  bool               `json:"is_banned"`
	BanReason *string            `json:"ban_reason"`
	CreatedAt time.Time          `json:"created_at"`
} // @name SearchedUserInfo

type EditUserData = database.UpdateUserData // @name EditUserData

type EditUserResult struct {
	UserId   string             `json:"user_id"`
	FullName string             `json:"full_name"`
	Email    string             `json:"email"`
	Role     appValues.UserRole `json:"role"`
} // @name EditUserResult

type GetUserInfoResult struct {
	UserId    string             `json:"user_id"`
	FullName  string             `json:"full_name"`
	Email     string             `json:"email"`
	Role      appValues.UserRole `json:"role"`
	IsBanned  bool               `json:"is_banned"`
	BanReason *string            `json:"ban_reason"`
	CreatedAt time.Time          `json:"created_at"`
} // @name GetUserInfoResult

type BanUserData = database.BanUserData // @name BanUserData

type BanUserResult struct {
	UserId    string  `json:"user_id"`
	IsBanned  bool    `json:"is_banned"`
	BanReason *string `json:"ban_reason"`
} // @name BanUserResult

type ChangePasswordData struct {
	UserId      string `json:"user_id"`
	NewPassword string `json:"new_password"`
	Lang        string `json:"lang" default:"en"`

	// if the user itself is trying to change their password,
	// they should get an email, click on the email link
	// get redirected to the special change password page
	// which contains their token, and then that page will
	// have to send the following parameters:
} // @name ChangePasswordData

type ChangePasswordResult struct {
	EmailSent       bool   `json:"email_sent"`
	PasswordChanged bool   `json:"password_changed"`
	Lang            string `json:"lang" default:"en"`
} // @name ChangePasswordResult

type ConfirmChangePasswordData struct {
	RTParam     string `json:"rt_param"`
	RTHash      string `json:"rt_hash"`
	RTVerifier  string `json:"rt_verifier"`
	RqId        string `json:"rq_id"`
	NewPassword string `json:"new_password"`
} // @name ConfirmChangePasswordData

// changePasswordRequestEntry is an internal type to keep track of
// password-change requests per-user.
type changePasswordRequestEntry struct {
	UserId    string
	LastTryAt time.Time
	TryCount  int
	mut       *sync.Mutex
	RqId      string
	RTParam   string
	LTNum     int32
}
