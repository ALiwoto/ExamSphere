package apiHandlers

import (
	"ExamSphere/src/core/utils/logging"

	"github.com/gofiber/fiber/v2"
)

func ApiPanicHandler(c *fiber.Ctx, err any) error {
	logging.UnexpectedPanic(err)

	return SendErrInternalServerError(c)
}
