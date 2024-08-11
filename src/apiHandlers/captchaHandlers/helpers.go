package captchaHandlers

import "OnlineExams/src/core/appValues"

func init() {
	appValues.VerifyCaptchaHandler = VerifyCaptcha
}

func VerifyCaptcha(clientRId, captchaId, captchaAnswer string) bool {
	if !appValues.IsClientRIDValid(clientRId) ||
		captchaId == "" || captchaAnswer == "" {
		return false
	}

	return clientRIDMap.Exists(clientRId) &&
		captchaStore.Verify(captchaId, captchaAnswer, true)
}
