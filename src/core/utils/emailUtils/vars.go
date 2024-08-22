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

	//go:embed templates/ConfirmAccount.en.html
	ConfirmAccountEmailTemplate_en string
)

var (
	PasswordChangeTemplateMap = map[string]string{
		"en": PasswordChangeEmailTemplate_en,
	}

	ConfirmAccountTemplateMap = map[string]string{
		"en": ConfirmAccountEmailTemplate_en,
	}
)