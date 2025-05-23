package errors

import "errors"

var (
    ErrUserNotFound = errors.New("user not found")
    ErrWrongPassword = errors.New("wrong password")
    ErrUserExists = errors.New("user already exists")
)