package userHandlers

import (
	"ExamSphere/src/apiHandlers"
	"ExamSphere/src/core/appConfig"
	"ExamSphere/src/database"

	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	fUtils "github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userInfo *database.UserInfo) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userInfo.UserId
	claims["refresh"] = false
	claims["auth_hash"] = userInfo.AuthHash
	claims["exp"] = time.Now().Add(appConfig.AccessTokenExpiration).Unix()
	accessToken, _ := token.SignedString(appConfig.AccessTokenSigningKey)
	return accessToken
}

func GenerateRefreshToken(userInfo *database.UserInfo) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userInfo.UserId
	claims["refresh"] = true
	claims["auth_hash"] = userInfo.AuthHash
	claims["exp"] = time.Now().Add(appConfig.RefreshTokenExpiration).Unix()
	refreshToken, _ := token.SignedString(appConfig.RefreshTokenSigningKey)
	return refreshToken
}

// getLoginExpiration returns the expiration time of the login
// done on the master server. Do note that this will send the
// little bit less, so that the client can refresh the token
// before it expires.
func getLoginExpiration() int64 {
	return time.Now().Add(
		appConfig.AccessTokenExpiration - time.Hour,
	).Unix()
}

func IsInvalidPassword(value string) bool {
	return len(value) < MinPasswordLength ||
		len(value) > MaxPasswordLength
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return apiHandlers.SendErrMalformedJWT(c)
	}

	return apiHandlers.SendErrInvalidJWT(c)
}

func isRateLimited(c *fiber.Ctx) bool {
	path := strings.ToLower(fUtils.CopyString(c.Path()))
	entryKey := strings.ToLower(fUtils.CopyString(c.IP())) + "_" +
		path

	entryValue := requestRateLimitMap.Get(entryKey)
	if entryValue == nil {
		entryValue = &userRequestEntry{
			RequestPath: path,
			LastTryAt:   time.Now(),
			TryCount:    1,
			mut:         &sync.Mutex{},
		}
		requestRateLimitMap.Add(entryKey, entryValue)
		return false
	}

	entryValue.mut.Lock()
	defer entryValue.mut.Unlock()

	// check time
	if time.Since(entryValue.LastTryAt) > appConfig.GetMaxRateLimitDuration() {
		entryValue.TryCount = 1
		entryValue.LastTryAt = time.Now()
		return false
	}

	// is already rate limited?
	if entryValue.TryCount > appConfig.GetMaxRequestTillRateLimit() {
		if time.Since(entryValue.LastTryAt) > appConfig.GetRateLimitPunishmentDuration() {
			// it should get released now
			entryValue.TryCount = 1
			entryValue.LastTryAt = time.Now()
			return false
		}
		return true
	}

	entryValue.TryCount++
	if entryValue.TryCount > appConfig.GetMaxRequestTillRateLimit() {
		entryValue.LastTryAt = time.Now()
		return true
	}

	return false
}

func toSearchedUsersResult(users []*database.UserInfo) []SearchedUserInfo {
	searchedUsers := make([]SearchedUserInfo, 0, len(users))

	for _, user := range users {
		searchedUsers = append(searchedUsers, SearchedUserInfo{
			UserId:   user.UserId,
			FullName: user.FullName,
			Role:     user.Role,
		})
	}

	return searchedUsers
}
