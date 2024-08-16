package emailUtils

import (
	"ExamSphere/src/core/appConfig"
	"encoding/base64"
	"fmt"

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

func SendChangePassword(data *ChangePasswordEmailData) error {
	if emailSenderClient == nil {
		return fmt.Errorf("email client not loaded")
	}

	e := email.NewEmail()
	e.From = emailSenderClient.EmailFrom
	e.To = []string{data.EmailTo}
	e.Subject = "Change Password"

	htmlTemplate, ok := PasswordChangeTemplateMap[data.Lang]
	if !ok {
		// default to english
		htmlTemplate = PasswordChangeTemplateMap["en"]
	}

	e.HTML = []byte(fmt.Sprintf(htmlTemplate, data.UserFullName, data.ChangeLink))
	print(string(e.HTML))

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
