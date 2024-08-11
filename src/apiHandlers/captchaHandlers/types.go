package captchaHandlers

type GetCaptchaResult struct {
	CaptchaId string `json:"captcha_id"`
	Captcha   string `json:"captcha"`
	ClientRId string `json:"client_r_id"`
}

type storedCaptchaInfo struct {
	ClientRID string
	CaptchaID string
}
