package masterServer

import (
	"ExamSphere/src/apiHandlers/captchaHandlers"
	"ExamSphere/src/apiHandlers/sudoHandlers"
	"ExamSphere/src/apiHandlers/swaggerHandlers"
	"ExamSphere/src/apiHandlers/userHandlers"
	"ExamSphere/src/core/appConfig"
	"ExamSphere/src/core/appValues"
	"ExamSphere/src/core/utils/logging"
	"ExamSphere/src/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func RunServer() error {
	err := LoadDatabase()
	if err != nil {
		logging.Error("RunMasterServer: failed to load database: ", err)
		return err
	}

	appValues.ServerEngine = fiber.New(fiber.Config{
		ProxyHeader:   appConfig.GetIPProxyHeader(),
		CaseSensitive: CaseSensitive,
	})

	if appConfig.IsDebug() {
		LoadSwaggerHandler(appValues.ServerEngine)
	}

	LoadMiddlewares(appValues.ServerEngine)
	LoadHandlersV1(appValues.ServerEngine)

	// make sure to load the UI files at the end
	// because it uses the pattern * which will
	// match all the routes
	LoadUIFiles(appValues.ServerEngine)

	if appConfig.TheConfig.CertFile != "" {
		return appValues.ServerEngine.ListenTLS(
			appConfig.TheConfig.BindAddress,
			appConfig.GetCertFile(),
			appConfig.GetCertKeyFile(),
		)
	}

	return appValues.ServerEngine.Listen(appConfig.TheConfig.BindAddress)
}

func LoadHandlersV1(app *fiber.App) {
	authProtection := userHandlers.AuthProtection()
	refreshAuthProtection := userHandlers.RefreshAuthProtection()

	v1 := app.Group("/api/v1")

	// user handlers
	v1.Post("/user/login", userHandlers.LoginV1)
	v1.Post("/user/reAuth", refreshAuthProtection, userHandlers.ReAuthV1)
	v1.Get("/user/me", authProtection, userHandlers.GetMeV1)
	v1.Post("/user/create", authProtection, userHandlers.CreateUserV1)
	v1.Post("/user/search", authProtection, userHandlers.SearchUserV1)

	// captcha handlers
	v1.Get("/captcha/generate", captchaHandlers.GenerateCaptchaV1)

	// sudo handlers
	v1.Post("/sudo/exit", sudoHandlers.ExitV1)
}

// @securityDefinitions.basic BasicAuth
func LoadSwaggerHandler(app *fiber.App) {
	app.Get("/swagger/swagger.*", swaggerHandlers.GetSwagger)

	app.Get("/swagger/*", swagger.New(
		swagger.Config{
			InstanceName:    appConfig.GetSwaggerInstanceName(),
			Title:           appConfig.GetSwaggerTitle(),
			URL:             appConfig.GetSwaggerBaseURL() + "/swagger/swagger.json",
			TryItOutEnabled: true,
		},
	))
}

func LoadMiddlewares(app *fiber.App) {
	app.Use(cors.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Server", "Microsoft-IIS/10.0")
		c.Set("X-Powered-By", "PHP/8.2.8")

		return c.Next()
	})
}

func LoadUIFiles(app *fiber.App) {
	app.Static("/", "./ui")
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("ui/index.html")
	})
}

func LoadDatabase() error {
	return database.StartDatabase()
}
