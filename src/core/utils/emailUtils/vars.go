package emailUtils

import (
	_ "embed"
)

var (
	emailSenderClient *EmailSenderClient
)

var (
	//go:embed templates/PasswordChange.en.html
	PasswordChangeEmailTemplate_en string
)

var (
	PasswordChangeTemplateMap = map[string]string{
		"en": PasswordChangeEmailTemplate_en,
	}
)
