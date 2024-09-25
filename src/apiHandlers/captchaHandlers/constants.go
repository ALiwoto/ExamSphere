package captchaHandlers

import "time"

const (
	// CaptchaLimitAmount is the number of captchas created that triggers
	// garbage collection used by the captcha store.
	CaptchaLimitAmount = 8192
	// CaptchaExpirationTime is expiration time of captchas used by default store.
	CaptchaExpirationTime = 10 * time.Minute

	// captchaType is the type of captcha sent to users.
	// Add this to config later (maybe).
	captchaType = "digit"

	// StringCaptchaValues is the set of characters that can be used in the captcha.
	StringCaptchaValues = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// CaptchaSizeHeight is the height of the captcha image.
	CaptchaSizeHeight = 80

	// CaptchaSizeWidth is the width of the captcha image.
	CaptchaSizeWidth = 240

	// CaptchaCharsLength is the number of characters in the captcha.
	CaptchaCharsLength = 6

	CaptchaNoiseCount = 5
)
