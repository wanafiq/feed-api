package email

import "embed"

//go:embed "*.tmpl"
var Templates embed.FS

var (
	ConfirmationEmail         = "confirmation_email"
	ConfirmationEmailTemplate = "confirmation_email.tmpl"
)
