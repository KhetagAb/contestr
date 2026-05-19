package repository

import "errors"

var (
	ErrContestAlreadyRegistered = errors.New("contest already registered")
	ErrContestNotRegistered     = errors.New("contest not registered")
	ErrHandleNotFound           = errors.New("handle mapping not found")
)
