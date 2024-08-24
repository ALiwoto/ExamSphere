package userHandlers

import (
	"ExamSphere/src/apiHandlers"
	"ExamSphere/src/core/appConfig"
	"ExamSphere/src/core/appValues"
	"ExamSphere/src/core/utils/emailUtils"
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
	} else if !IsEmailValid(newUserData.Email) { // email is mandatory
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
		return apiHandlers.SendErrInvalidUserID(c)
	}

	newUserInfo, err = database.CreateNewUser(&database.NewUserData{
		UserId:         newUserData.UserId,
		FullName:       newUserData.FullName,
		Email:          newUserData.Email,
		RawPassword:    newUserData.RawPassword,
		Role:           newUserData.Role,
		UserAddress:    newUserData.UserAddress,
		PhoneNumber:    newUserData.PhoneNumber,
		SetupCompleted: newUserData.SetupCompleted,
	})
	if err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "unique key") && strings.Contains(errStr, "email") {
			return apiHandlers.SendErrEmailAlreadyExists(c)
		}

		logging.Error("CreateUserV1: failed to create new user: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	} else if newUserInfo == nil {
		return apiHandlers.SendErrInternalServerError(c)
	}

	if !newUserInfo.SetupCompleted && emailUtils.IsEmailClientLoaded() {
		// we need to send account confirmation email for them
		entry, err := newConfirmAccountRequest(newUserInfo)
		if err != nil {
			logging.Error("CreateUserV1: failed to create confirm account request: ", err)
			return apiHandlers.SendErrInternalServerError(c)
		} else if entry == nil {
			return apiHandlers.SendErrInternalServerError(c)
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					logging.Error("CreateUserV1: failed to send confirm account email: ", r)
				}
			}()

			err := emailUtils.SendConfirmAccountEmail(&emailUtils.ConfirmAccountEmailData{
				UserFullName: newUserInfo.FullName,
				ChangeLink:   entry.GetRedirectAddress(appConfig.GetConfirmAccountBaseURL()),
				EmailTo:      newUserInfo.Email,
				Lang:         newUserData.PrimaryLanguage,
			})
			if err != nil {
				logging.Error("CreateUserV1: failed to send confirm account email: ", err)
			}
		}()
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
	} else if searchUserData.Limit <= 0 {
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
	} else if updateUserData.Email != "" && !IsEmailValid(updateUserData.Email) {
		return apiHandlers.SendErrInvalidEmail(c)
	}

	if updateUserData.IsEmpty() {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	targetUserInfo, err := database.GetUserByUserId(updateUserData.UserId)
	if err != nil {
		if err == database.ErrUserNotFound {
			return apiHandlers.SendErrInvalidUserID(c)
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
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "unique key") && strings.Contains(errStr, "email") {
			return apiHandlers.SendErrEmailAlreadyExists(c)
		}

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

// GetUserInfoV1 godoc
// @Summary Get a user's information
// @Description Allows a user to get another user's information by their user ID
// @ID getUserInfoV1
// @Tags User
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id query string true "User ID"
// @Success 200 {object} apiHandlers.EndpointResponse{result=GetUserInfoResult}
// @Router /api/v1/user/info [get]
func GetUserInfoV1(c *fiber.Ctx) error {
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

	targetUserId := c.Query("id")
	if targetUserId == "" {
		return apiHandlers.SendErrQueryParameterNotProvided(c, "id")
	}

	targetUserInfo, err := database.GetUserByUserId(targetUserId)
	if err != nil {
		if err == database.ErrUserNotFound {
			return apiHandlers.SendErrInvalidUserID(c)
		}

		logging.Error("GetUserInfoV1: failed to get user: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	if !userInfo.CanGetUserInfo(targetUserInfo) {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	return apiHandlers.SendResult(c, &GetUserInfoResult{
		UserId:   targetUserInfo.UserId,
		FullName: targetUserInfo.FullName,
		Email:    targetUserInfo.Email,
		Role:     targetUserInfo.Role,
	})
}

// BanUserV1 godoc
// @Summary Ban a user
// @Description Allows a user to ban another user
// @ID banUserV1
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param banUserData body BanUserData true "Ban user data"
// @Success 200 {object} apiHandlers.EndpointResponse{result=BanUserResult}
// @Router /api/v1/user/ban [post]
func BanUserV1(c *fiber.Ctx) error {
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

	banData := &BanUserData{}
	if err := c.BodyParser(banData); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	targetUserInfo, err := database.GetUserByUserId(banData.UserId)
	if err != nil {
		if err == database.ErrUserNotFound {
			return apiHandlers.SendErrInvalidUserID(c)
		}

		logging.Error("BanUserV1: failed to get user: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	if !userInfo.CanBanUser(targetUserInfo) {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	targetUserInfo, err = database.BanUser(banData)
	if err != nil {
		logging.Error("BanUserV1: failed to ban user: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, &BanUserResult{
		UserId:    targetUserInfo.UserId,
		IsBanned:  targetUserInfo.IsBanned,
		BanReason: targetUserInfo.BanReason,
	})
}

// ChangePasswordV1 godoc
// @Summary Change a user's password
// @Description Allows a user to change a user's password.
// If the user is trying to change their own password, they should get an email,
// click on the email link, get redirected to the special change password page
// which contains their token-parameters, and then that page will
// have to send the new password in confirmChangePassword endpoint.
// @ID changePasswordV1
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param changePasswordData body ChangePasswordData true "Change password data"
// @Success 200 {object} apiHandlers.EndpointResponse{result=ChangePasswordResult}
// @Router /api/v1/user/changePassword [post]
func ChangePasswordV1(c *fiber.Ctx) error {
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

	newPasswordData := &ChangePasswordData{}
	if err := c.BodyParser(newPasswordData); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	targetUserInfo, err := database.GetUserByUserId(newPasswordData.UserId)
	if err != nil {
		if err == database.ErrUserNotFound {
			return apiHandlers.SendErrInvalidUserID(c)
		}

		logging.Error("ChangePasswordV1: failed to get user: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	shouldSendEmail := targetUserInfo.UserId == userInfo.UserId

	if shouldSendEmail {
		// the user is trying to change their own password
		// so they should get an email, click on the email link
		// get redirected to the special change password page
		// which contains their token, and then that page will
		// call the change password endpoint with the token
		entry, err := newChangePasswordRequest(userInfo)
		if err != nil {
			if err == ErrTooManyPasswordChangeAttempts {
				return apiHandlers.SendErrTooManyPasswordChangeAttempts(c)
			}

			logging.Error("ChangePasswordV1: failed to create change password request: ", err)
		}

		if entry == nil {
			return apiHandlers.SendErrInternalServerError(c)
		}

		err = emailUtils.SendChangePasswordEmail(&emailUtils.ChangePasswordEmailData{
			UserFullName: userInfo.FullName,
			ChangeLink:   entry.GetRedirectAddress(appConfig.GetChangePasswordBaseURL()),
			EmailTo:      userInfo.Email,
			Lang:         newPasswordData.Lang,
		})
		if err != nil {
			logging.Error("ChangePasswordV1: failed to send change password email: ", err)
			return apiHandlers.SendErrInternalServerError(c)
		}

		return apiHandlers.SendResult(c, &ChangePasswordResult{
			EmailSent:       true,
			PasswordChanged: false,
			Lang:            newPasswordData.Lang,
		})
	}

	if !userInfo.CanChangePassword(targetUserInfo) {
		return apiHandlers.SendErrPermissionDenied(c)
	} else if IsInvalidPassword(newPasswordData.NewPassword) {
		return apiHandlers.SendErrInvalidInputPass(c)
	}

	err = database.UpdateUserPassword(&database.UpdateUserPasswordData{
		UserId:      targetUserInfo.UserId,
		RawPassword: newPasswordData.NewPassword,
	})
	if err != nil {
		logging.Error("ChangePasswordV1: failed to update password: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, &ChangePasswordResult{
		EmailSent:       false,
		PasswordChanged: true,
		Lang:            newPasswordData.Lang,
	})
}

// ConfirmChangePasswordV1 godoc
// @Summary Confirm changing your own's password
// @Description Allows a user to confirm changing their own's password (from redirected page)
// @ID confirmChangePasswordV1
// @Tags User
// @Accept json
// @Produce json
// @Param confirmChangePasswordData body ConfirmChangePasswordData true "Confirm change password data"
// @Success 200 {object} apiHandlers.EndpointResponse{result=bool}
// @Router /api/v1/user/confirmChangePassword [post]
func ConfirmChangePasswordV1(c *fiber.Ctx) error {
	confirmData := &ConfirmChangePasswordData{}
	if err := c.BodyParser(confirmData); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	if confirmData.RqId == "" {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	entry := getChangePasswordRequest(confirmData.RqId)
	if entry == nil {
		// this basically means the request is expired, the user
		// should try sending another request later
		return apiHandlers.SendErrRequestExpired(c)
	}

	if !entry.Verify(confirmData) {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	if IsInvalidPassword(confirmData.NewPassword) {
		return apiHandlers.SendErrInvalidInputPass(c)
	}

	// all is fine, we can now update the password in the database
	err := database.UpdateUserPassword(&database.UpdateUserPasswordData{
		UserId:      entry.UserId,
		RawPassword: confirmData.NewPassword,
	})
	if err != nil {
		logging.Error("ConfirmChangePasswordV1: failed to update password: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, true)
}

// ConfirmAccountV1 godoc
// @Summary Confirm account
// @Description Allows a user to confirm their account
// @ID confirmAccountV1
// @Tags User
// @Accept json
// @Produce json
// @Param confirmAccountData body ConfirmAccountData true "Confirm account data"
// @Success 200 {object} apiHandlers.EndpointResponse{result=bool}
// @Router /api/v1/user/confirmAccount [post]
func ConfirmAccountV1(c *fiber.Ctx) error {
	confirmData := &ConfirmAccountData{}
	if err := c.BodyParser(confirmData); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	if !confirmData.IsValid() {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	userInfo, err := database.GetUserByUserId(confirmData.UserId)
	if err != nil {
		if err == database.ErrUserNotFound {
			return apiHandlers.SendErrInvalidUserID(c)
		}

		return apiHandlers.SendErrInternalServerError(c)
	} else if userInfo == nil {
		return apiHandlers.SendErrInvalidUserID(c)
	}

	if userInfo.SetupCompleted {
		return apiHandlers.SendErrAccountAlreadyConfirmed(c)
	}

	if !verifyAccountConfirmation(userInfo, confirmData) {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	if !appValues.IsPasswordValid(confirmData.RawPassword) {
		return apiHandlers.SendErrInvalidInputPass(c)
	}

	// all is fine, we can now update the password in the database
	err = database.UpdateUserPassword(&database.UpdateUserPasswordData{
		UserId:      userInfo.UserId,
		RawPassword: confirmData.RawPassword,
	})
	if err != nil {
		logging.Error("ConfirmAccountV1: failed to update password: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	err = database.ConfirmUserAccount(userInfo.UserId)
	if err != nil {
		logging.Error("ConfirmAccountV1: failed to confirm account: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, true)
}
