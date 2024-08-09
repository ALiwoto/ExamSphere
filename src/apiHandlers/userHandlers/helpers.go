package userHandlers

import "github.com/gofiber/fiber/v2"

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
