package appValues

import "github.com/gofiber/fiber/v2"

var (
	ServerEngine *fiber.App
)

var (
	VerifyCaptchaHandler func(clientRId, captchaId, captchaAnswer string) bool
)
