package captchaHandlers

import (
	"OnlineExams/src/apiHandlers"
	"OnlineExams/src/core/appValues"
	"OnlineExams/src/core/utils/logging"
	"image/color"

	"github.com/gofiber/fiber/v2"
	fUtils "github.com/gofiber/fiber/v2/utils"
	"github.com/mojocn/base64Captcha"
)

// GetCaptchaV1 godoc
// @Summary Get a captcha
// @ID GenerateCaptchaV1
// @Description Allows a client (with Client-R-ID) to generate a captcha
// @Tags User
// @Produce json
// @Param Client-R-ID query string true "Client-R-ID"
// @Success 200 {object} apiHandlers.EndpointResponse{result=GetCaptchaResult}
// @Router /api/v1/captcha/generate [get]
func GenerateCaptchaV1(c *fiber.Ctx) error {
	clientRId := fUtils.CopyString(c.Query("Client-R-ID"))
	if !appValues.IsClientRIDValid(clientRId) {
		return apiHandlers.SendErrInvalidClientRId(c, clientRId)
	}

	// add rate-limiting here later in future...
	var driver base64Captcha.Driver
	switch captchaType {
	case "string":
		driver = base64Captcha.NewDriverString(
			CaptchaSizeHeight, CaptchaSizeWidth,
			CaptchaNoiseCount, 0, 6,
			StringCaptchaValues,
			&color.RGBA{0, 0, 0, 0},
			nil, []string{},
		)
	default:
		driver = base64Captcha.NewDriverDigit(
			CaptchaSizeHeight, CaptchaSizeWidth,
			CaptchaCharsLength,
			0,
			6,
		)
	}

	captcha := base64Captcha.NewCaptcha(driver, captchaStore)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		logging.Error("GetCaptchaV1: failed to generate captcha: ", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, &GetCaptchaResult{
		CaptchaId: id,
		Captcha:   b64s,
		ClientRId: clientRId,
	})
}
