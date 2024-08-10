package swaggerHandlers

import (
	_ "embed"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	//go:embed swagger.json
	swaggerJsonContent []byte

	// go:embed swagger.yaml
	swaggerYamlContent []byte
)

func GetSwagger(c *fiber.Ctx) error {
	path := c.Path()
	if strings.HasSuffix(path, ".yaml") {
		c.Set(fiber.HeaderContentType, "application/x-yaml")
		return c.Send(swaggerYamlContent)
	}

	// otherwise, send JSON (e.g. for /swagger/swagger.json)
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Send(swaggerJsonContent)
}
