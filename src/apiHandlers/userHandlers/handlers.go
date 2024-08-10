package userHandlers

import (
	"OnlineExams/src/apiHandlers"
	"OnlineExams/src/core/appConfig"
	"OnlineExams/src/database"

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
// @Tags User
// @Accept json
// @Produce json
// @Param loginData body LoginData true "Login data"
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

	return apiHandlers.SendResult(c, &LoginResult{
		UserId:       userInfo.UserId,
		FullName:     userInfo.FullName,
		Role:         userInfo.Role.ToString(),
		AccessToken:  GenerateAccessToken(userInfo),
		RefreshToken: GenerateRefreshToken(userInfo),
		Expiration:   getLoginExpiration(),
	})
}

func AuthV1(c *fiber.Ctx) error {
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
