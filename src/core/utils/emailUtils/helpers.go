package emailUtils

import (
	"ExamSphere/src/core/appConfig"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/jordan-wright/email"
)

func LoadEmailClient() error {
	if appConfig.TheConfig.EmailFrom == "" ||
		appConfig.TheConfig.EmailUser == "" ||
		appConfig.TheConfig.EmailPass == "" ||
		appConfig.TheConfig.EmailHost == "" ||
		appConfig.TheConfig.EmailPort == 0 {
		return ErrInvalidEmailConfig
	}

	decodedPass, err := DecodeSpecificPassword(appConfig.TheConfig.EmailPass)
	if err != nil {
		return err
	}

	emailSenderClient = &EmailSenderClient{
		EmailFrom: appConfig.TheConfig.EmailFrom,
		EmailUser: appConfig.TheConfig.EmailUser,
		EmailPass: decodedPass,
		EmailHost: appConfig.TheConfig.EmailHost,
		EmailPort: appConfig.TheConfig.EmailPort,
	}

	return nil
}

func IsEmailClientLoaded() bool {
	return emailSenderClient != nil
}

func SendChangePasswordEmail(data *ChangePasswordEmailData) error {
	if emailSenderClient == nil {
		return ErrEmailClientNotLoaded
	}

	e := email.NewEmail()
	e.From = emailSenderClient.EmailFrom
	e.To = []string{data.EmailTo}
	e.Subject = "Change Password"

	htmlTemplate, ok := PasswordChangeTemplateMap[data.Lang]
	if !ok {
		// default to english
		htmlTemplate = PasswordChangeTemplateMap[DefaultEmailLanguage]
	}

	e.HTML = []byte(fmt.Sprintf(htmlTemplate, data.UserFullName, data.ChangeLink))

	err := e.Send(emailSenderClient.GetHostAddress(), emailSenderClient.GetSmtpAuth())
	if err != nil {
		return err
	}

	return nil
}

func SendConfirmAccountEmail(data *ConfirmAccountEmailData) error {
	if emailSenderClient == nil {
		return ErrEmailClientNotLoaded
	}

	e := email.NewEmail()
	e.From = emailSenderClient.EmailFrom
	e.To = []string{data.EmailTo}
	e.Subject = "Account Confirmation"

	htmlTemplate, ok := ConfirmAccountTemplateMap[data.Lang]
	if !ok {
		// default to english
		htmlTemplate = ConfirmAccountTemplateMap[DefaultEmailLanguage]
	}

	e.HTML = []byte(fmt.Sprintf(htmlTemplate, data.UserFullName, data.ChangeLink))

	err := e.Send(emailSenderClient.GetHostAddress(), emailSenderClient.GetSmtpAuth())
	if err != nil {
		return err
	}

	return nil
}

func DecodeSpecificPassword(encodedPassword string) (string, error) {
	// Base64 decode the whole string
	decoded, err := base64.StdEncoding.DecodeString(encodedPassword)
	if err != nil {
		return "", err
	}

	// Convert to string and split to get the encoded parts
	parts := ssg.Split(string(decoded), []string{
		"passM66QFT_",
		"_s5rS58O0O3ML_",
		"_RendPassTS5S",
	}...)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid encoded password format")
	}

	passResult := ""
	for _, part := range parts {
		currentPart, err := base64.StdEncoding.DecodeString(part)
		if err != nil {
			return "", err
		}

		passResult += string(currentPart)
	}

	return passResult, nil
}

func fixTemplateFormatting(template string) string {
	template = strings.ReplaceAll(template, "100%", "100%%")

	return template
}
