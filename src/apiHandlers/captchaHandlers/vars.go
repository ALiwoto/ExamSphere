package captchaHandlers

import (
	"time"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/mojocn/base64Captcha"
)

var (
	captchaStore = base64Captcha.NewMemoryStore(CaptchaLimitAmount, CaptchaExpirationTime)
	clientRIDMap = func() *ssg.SafeEMap[string, storedCaptchaInfo] {
		m := ssg.NewSafeEMap[string, storedCaptchaInfo]()
		m.SetExpiration(CaptchaExpirationTime)
		m.SetInterval(time.Hour)
		m.EnableChecking()

		return m
	}()
)
