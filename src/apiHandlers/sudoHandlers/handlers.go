package sudoHandlers

import (
	"ExamSphere/src/apiHandlers"
	"ExamSphere/src/core/appConfig"
	"ExamSphere/src/core/appValues"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ExitV1(c *fiber.Ctx) error {
	sudoToken := c.Get(SudoTokenHeaderName)
	if len(sudoToken) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if !appConfig.IsSudoToken(sudoToken) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	go func() {
		time.Sleep(time.Second)
		_ = appValues.ServerEngine.ShutdownWithTimeout(time.Second)
	}()
	return apiHandlers.SendResult(c, "Exiting...")
}
