package users

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrEmailExists        = errors.New("email_already_exists")
	ErrUserNotFound       = errors.New("user_not_found")
)
