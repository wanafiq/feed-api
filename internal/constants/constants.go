package constants

import "time"

var (
	DevEnv = "development"

	QueryTimeout = time.Second * 5

	RoleUser      = "user"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"

	ConfirmationToken           = "confirmation_token"
	ConfirmationTokenExpireTime = time.Hour * 24 * 3 // 3 days
)
