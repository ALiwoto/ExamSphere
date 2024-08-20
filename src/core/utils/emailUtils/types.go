package emailUtils

type EmailSenderClient struct {
	EmailFrom string
	EmailUser string
	EmailPass string
	EmailHost string
	EmailPort int
}

type ChangePasswordEmailData struct {
	UserFullName string
	ChangeLink   string
	EmailTo      string
	Lang         string
}

type ConfirmAccountEmailData struct {
	UserFullName string
	ChangeLink   string
	EmailTo      string
	Lang         string
}
