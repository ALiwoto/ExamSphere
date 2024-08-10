package masterServer

import (
	"OnlineExams/src/apiHandlers/userHandlers"
	"OnlineExams/src/core/appConfig"
	"OnlineExams/src/core/appValues"
	"OnlineExams/src/core/utils/logging"
	"OnlineExams/src/database"

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

	LoadMiddlewares(appValues.ServerEngine)
	LoadHandlersV1(appValues.ServerEngine)

	if appConfig.IsDebug() {
		LoadSwaggerHandler(appValues.ServerEngine)
	}

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
	v1.Post("/user/auth", refreshAuthProtection, userHandlers.AuthV1)
	v1.Get("/user/me", authProtection, userHandlers.GetMeV1)
}

func LoadSwaggerHandler(app *fiber.App) {
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

func LoadDatabase() error {
	return database.StartDatabase()
}
