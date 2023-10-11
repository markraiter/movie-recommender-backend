package models

import "errors"

var (
	ErrNotFound             = errors.New("not found")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrWrongTokenClaimType  = errors.New("wrong token claim type")
	ErrUniqueViolation      = errors.New("this user is already exist")
	ErrWrongCredentials     = errors.New("invalid login or password")
	ErrNoSession            = errors.New("there is no such session")
	ErrEmptyCookie          = errors.New("empty cookie")
)
