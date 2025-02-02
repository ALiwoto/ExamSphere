package apiUtils

import (
	"ExamSphere/src/apiHandlers"
	"ExamSphere/src/database"

	"github.com/gofiber/fiber/v2"
)

func GetUserInfo(c *fiber.Ctx) (*database.UserInfo, error) {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return nil, apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return nil, apiHandlers.SendErrInvalidAuth(c)
	}

	return userInfo, nil
}
