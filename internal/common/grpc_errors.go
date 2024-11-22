package common

import "errors"

var (
	ErrNotFound    = errors.New("Not found")
	ErrLoginExists = errors.New("Login already exists")
	ErrEmailExists = errors.New("Email already exists")
)
