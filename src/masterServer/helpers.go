package masterServer

import (
	"ExamSphere/src/apiHandlers"
	"ExamSphere/src/apiHandlers/captchaHandlers"
	"ExamSphere/src/apiHandlers/courseHandlers"
	"ExamSphere/src/apiHandlers/examHandlers"
	"ExamSphere/src/apiHandlers/sudoHandlers"
	"ExamSphere/src/apiHandlers/swaggerHandlers"
	"ExamSphere/src/apiHandlers/topicHandlers"
	"ExamSphere/src/apiHandlers/userHandlers"
	"ExamSphere/src/core/appConfig"
	"ExamSphere/src/core/appValues"
	"ExamSphere/src/core/utils/emailUtils"
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

	LoadEmailClient()

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
	v1.Get("/user/info", authProtection, userHandlers.GetUserInfoV1)
	v1.Post("/user/search", authProtection, userHandlers.SearchUserV1)
	v1.Post("/user/edit", authProtection, userHandlers.EditUserV1)
	v1.Post("/user/ban", authProtection, userHandlers.BanUserV1)
	v1.Post("/user/changePassword", authProtection, userHandlers.ChangePasswordV1)
	v1.Post("/user/confirmChangePassword", userHandlers.ConfirmChangePasswordV1)
	v1.Post("/user/confirmAccount", userHandlers.ConfirmAccountV1)

	// captcha handlers
	v1.Get("/captcha/generate", captchaHandlers.GenerateCaptchaV1)

	// topic handlers
	v1.Post("/topic/create", authProtection, topicHandlers.CreateTopicV1)
	v1.Post("/topic/search", authProtection, topicHandlers.SearchTopicV1)
	v1.Get("/topic/info", authProtection, topicHandlers.GetTopicInfoV1)
	v1.Post("/topic/userTopicStat", authProtection, topicHandlers.GetUserTopicStatV1)
	v1.Get("/topic/allUserTopicStats", authProtection, topicHandlers.GetAllUserTopicStatsV1)
	v1.Delete("/topic/delete", authProtection, topicHandlers.DeleteTopicV1)

	// course handlers
	v1.Post("/course/create", authProtection, courseHandlers.CreateCourseV1)
	v1.Post("/course/edit", authProtection, courseHandlers.EditCourseV1)
	v1.Get("/course/info", authProtection, courseHandlers.GetCourseInfoV1)
	v1.Post("/course/search", authProtection, courseHandlers.SearchCourseV1)
	v1.Post("/course/CreatedCourses", authProtection, courseHandlers.GetCreatedCoursesV1)
	v1.Post("/course/userCourses", authProtection, courseHandlers.GetUserCoursesV1)
	v1.Post("/course/courseParticipants", authProtection, courseHandlers.GetCourseParticipantsV1)

	// exam handlers
	v1.Post("/exam/create", authProtection, examHandlers.CreateExamV1)
	v1.Get("/exam/info", authProtection, examHandlers.GetExamInfoV1)
	v1.Post("/exam/edit", authProtection, examHandlers.EditExamV1)
	v1.Post("/exam/questions", authProtection, examHandlers.GetExamQuestionsV1)
	v1.Post("/exam/answer", authProtection, examHandlers.AnswerExamQuestionV1)
	v1.Post("/exam/setScore", authProtection, examHandlers.SetScoreV1)
	v1.Post("/exam/givenExam", authProtection, examHandlers.GetGivenExamV1)
	v1.Get("/exam/userOngoingExams", authProtection, examHandlers.GetUserOngoingExamsV1)
	v1.Post("/exam/userExamsHistory", authProtection, examHandlers.GetUserExamsHistoryV1)

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

		defer func() {
			if r := recover(); r != nil {
				_ = apiHandlers.ApiPanicHandler(c, r)
			}
		}()

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

func LoadEmailClient() {
	err := emailUtils.LoadEmailClient()
	if err != nil {
		logging.Warn("LoadEmailClient: failed to load email client: ", err)
		logging.Warn("Without an email client, some features of the platform will not work correctly.")
		logging.Warn("Please check the email configuration in the config file.")
	}
}
