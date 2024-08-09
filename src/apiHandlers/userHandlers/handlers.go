package userHandlers

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthProtection() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: wcf.AccessTokenSigningKey},
		ErrorHandler: jwtError,
	})
}

func RefreshAuthProtection() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: wcf.RefreshTokenSigningKey},
		ErrorHandler: jwtError,
	})
}
