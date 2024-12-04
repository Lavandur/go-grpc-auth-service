package common

import "errors"

var (
	ErrNotFound    = errors.New("Not found")
	ErrLoginExists = errors.New("Login already exists")
	ErrEmailExists = errors.New("Email already exists")

	ErrBuildQuery   = errors.New("Build query error")
	ErrConnectionDB = errors.New("Database connection error")

	ErrPermissionDenied = errors.New("Permission denied")
)
