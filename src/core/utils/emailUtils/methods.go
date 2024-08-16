package emailUtils

import (
	"net/smtp"

	"github.com/ALiwoto/ssg/ssg"
)

func (c *EmailSenderClient) GetHostAddress() string {
	return c.EmailHost + ":" + ssg.ToBase10(c.EmailPort)
}

func (c *EmailSenderClient) GetSmtpAuth() smtp.Auth {
	return smtp.PlainAuth("", c.EmailUser, c.EmailPass, c.EmailHost)
}
