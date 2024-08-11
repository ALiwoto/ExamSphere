package userHandlers

import (
	"OnlineExams/src/apiHandlers"
	"OnlineExams/src/core/appConfig"
	"OnlineExams/src/core/appValues"
	"OnlineExams/src/core/utils/logging"
	"OnlineExams/src/database"
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

	userInfo := database.GetUserInfoByPass(
		loginInput.UserId, loginInput.Password,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidUsernamePass(c)
	} else if userInfo.IsBanned {
		return apiHandlers.SendErrUserBanned(c)
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

	return apiHandlers.SendResult(c, &LoginResult{
		UserId:       userInfo.UserId,
		FullName:     userInfo.FullName,
		Role:         userInfo.Role.ToString(),
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
		Role:     userInfo.Role.ToString(),
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

	newRole := appValues.ParseRole(newUserData.RoleStr)
	if !userInfo.CanCreateRole(newRole) {
		return apiHandlers.SendErrPermissionDenied(c)
	} else if newRole.IsInvalid() {
		return apiHandlers.SendErrInvalidBodyData(c)
	} else if newUserData.Email == "" { // email is mandatory
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	newUserData.Role = newRole
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
