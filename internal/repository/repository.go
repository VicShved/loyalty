package repository

import (
	"errors"
)

type KeyLongURLStr struct {
	Key     string
	LongURL string
}

type RepoInterface interface {
	Register()
	Login()
	GetOrders()
	PostOrders()
	GetBalance()
	PostWithdraw()
	GetWithdrawals()
	Ping() error
}

var ErrPKConflict = errors.New("PK conflict")
