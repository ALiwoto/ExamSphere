package userHandlers

import (
	"ExamSphere/src/apiHandlers"
	"ExamSphere/src/core/appConfig"
	"ExamSphere/src/core/appValues"
	"ExamSphere/src/core/utils/logging"
	"ExamSphere/src/database"
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthProtection() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: appConfig.AccessTokenSigningKey},
		ErrorHandler: jwtError,
	})
}

func RefreshAuthProtection() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: appConfig.RefreshTokenSigningKey},
		ErrorHandler: jwtError,
	})
}

// LoginV1 godoc
// @Summary Login to the system
// @Description Allows a user to login to the system and obtain access/refresh tokens
// @ID loginV1
// @Tags User
// @Accept json
// @Produce json
// @Param loginData body LoginData true "Login data"
// @Success 200 {object} apiHandlers.EndpointResponse{result=LoginResult}
// @Router /api/v1/user/login [post]
func LoginV1(c *fiber.Ctx) error {
	if isRateLimited(c) {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	loginInput := &LoginData{}
	if err := c.BodyParser(loginInput); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	if appValues.VerifyCaptchaHandler != nil {
		// check captcha if and only if the captcha verifier is dependency-injected
		if !appValues.VerifyCaptchaHandler(
			loginInput.ClientRId,
			loginInput.CaptchaId, loginInput.CaptchaAnswer,
		) {
			return apiHandlers.SendErrInvalidCaptcha(c)
		}
	}

	userInfo := database.GetUserInfoByPass(
		loginInput.UserId, loginInput.Password,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidUsernamePass(c)
	} else if userInfo.IsBanned {
		return apiHandlers.SendErrUserBanned(c)
	}

	return apiHandlers.SendResult(c, &LoginResult{
		UserId:       userInfo.UserId,
		FullName:     userInfo.FullName,
		Role:         userInfo.Role,
		AccessToken:  GenerateAccessToken(userInfo),
		RefreshToken: GenerateRefreshToken(userInfo),
		Expiration:   getLoginExpiration(),
	})
}

// ReAuthV1 godoc
// @Summary Refresh the access token
// @Description Allows a user to refresh their access token
// @ID reAuthV1
// @Tags User
// @Produce json
// @Param Authorization header string true "Refresh token"
// @Success 200 {object} apiHandlers.EndpointResponse{result=AuthResult}
// @Router /api/v1/user/reAuth [post]
func ReAuthV1(c *fiber.Ctx) error {
	if isRateLimited(c) {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	if !claimInfo.Refresh {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	}

	return apiHandlers.SendResult(c, &AuthResult{
		UserId:       userInfo.UserId,
		FullName:     userInfo.FullName,
		Role:         userInfo.Role.ToString(),
		AccessToken:  GenerateAccessToken(userInfo),
		RefreshToken: GenerateRefreshToken(userInfo),
		Expiration:   getLoginExpiration(),
	})
}

// GetMeV1 godoc
// @Summary Get the user's information
// @Description Allows a user to get their own information
// @ID getMeV1
// @Tags User
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} apiHandlers.EndpointResponse{result=MeResult}
// @Router /api/v1/user/me [get]
func GetMeV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	}

	return apiHandlers.SendResult(c, &MeResult{
		UserId:   userInfo.UserId,
		FullName: userInfo.FullName,
		Role:     userInfo.Role,
	})
}

// CreateUserV1 godoc
// @Summary Create a new user
// @Description Allows a user to create a new user
// @ID createUserV1
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param createUserData body CreateUserData true "Create user data"
// @Success 200 {object} apiHandlers.EndpointResponse{result=CreateUserResult}
// @Router /api/v1/user/create [post]
func CreateUserV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	}

	newUserData := &CreateUserData{}
	if err := c.BodyParser(newUserData); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	if newUserData.Role.IsInvalid() {
		return apiHandlers.SendErrInvalidBodyData(c)
	} else if !userInfo.CanCreateRole(newUserData.Role) {
		return apiHandlers.SendErrPermissionDenied(c)
	} else if newUserData.Email == "" { // email is mandatory
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	createUserMutex.Lock()
	defer createUserMutex.Unlock()

	if IsInvalidPassword(newUserData.RawPassword) {
		return apiHandlers.SendErrInvalidInputPass(c)
	}

	newUserInfo, err := database.GetUserByUserId(newUserData.UserId)
	if err != nil && err != database.ErrUserNotFound {
		return apiHandlers.SendErrInternalServerError(c)
	}

	if newUserInfo != nil {
		return apiHandlers.SendErrUsernameExists(c)
	} else if !appValues.IsUserIdValid(newUserData.UserId) {
		return apiHandlers.SendErrInvalidUsername(c)
	}

	newUserInfo, err = database.CreateNewUser(newUserData)
	if err != nil {
		if strings.Contains(err.Error(), "violates check constraint") {
			return apiHandlers.SendErrInvalidBodyData(c)
		}

		logging.Error("CreateUserV1: failed to create new user: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	} else if newUserInfo == nil {
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, &CreateUserResult{
		UserId:   newUserInfo.UserId,
		Email:    newUserInfo.Email,
		FullName: newUserInfo.FullName,
		Role:     newUserInfo.Role.ToString(),
	})
}

// SearchUserV1 godoc
// @Summary Search users
// @Description Allows a user to search for users
// @ID searchUserV1
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param searchUserData body SearchUserData true "Search user data"
// @Success 200 {object} apiHandlers.EndpointResponse{result=SearchUserResult}
// @Router /api/v1/user/search [post]
func SearchUserV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	}

	if !userInfo.CanSearchUser() {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	searchUserData := &SearchUserData{}
	if err := c.BodyParser(searchUserData); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	} else if searchUserData.Query == "" || searchUserData.Limit <= 0 {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	users, err := database.SearchUser(searchUserData)
	if err != nil {
		logging.Error("SearchUserV1: failed to search users: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, &SearchUserResult{
		Users: toSearchedUsersResult(users),
	})
}

// EditUserV1 godoc
// @Summary Edit a user's basic information
// @Description Allows a user to edit another user. Users are not allowed to edit their own information.
// @ID editUserV1
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param editUserData body EditUserData true "Edit user data"
// @Success 200 {object} apiHandlers.EndpointResponse{result=EditUserResult}
// @Router /api/v1/user/edit [post]
func EditUserV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	}

	updateUserData := &EditUserData{}
	if err := c.BodyParser(updateUserData); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	} else if updateUserData.UserId == "" {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	if updateUserData.IsEmpty() {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	targetUserInfo, err := database.GetUserByUserId(updateUserData.UserId)
	if err != nil {
		if err == database.ErrUserNotFound {
			return apiHandlers.SendErrInvalidUsername(c)
		}

		logging.Error("EditUserV1: failed to get user: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	if !userInfo.CanEditUser(targetUserInfo.Role) {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	if updateUserData.UserId == userInfo.UserId {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	targetUserInfo, err = database.UpdateUserInfo(updateUserData)
	if err != nil {
		logging.Error("EditUserV1: failed to update user: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, &EditUserResult{
		UserId:   targetUserInfo.UserId,
		FullName: targetUserInfo.FullName,
		Email:    targetUserInfo.Email,
		Role:     targetUserInfo.Role,
	})
}
