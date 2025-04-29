package repository

import (
	"errors"
)

type BalanceType struct {
	Accural   string `json:"current"`
	Withdrawn string `json:"withdrawn"`
}

type OrderAccuralType struct {
	Order       string `json:"number"`
	Status      string `json:"status"`
	Accural     string `json:"accural"`
	Uploaded_at string `json:"uploaded_at"`
}

type OrderWithdrawType struct {
	Order        string `json:"order"`
	Withdraw     string `json:"sum"`
	Processed_at string `json:"processed_at"`
}

type RepoInterface interface {
	Register(login string, hashPassword string) (string, error)
	Login(login string, hashPassword string) (string, error)
	PostOrders(user_id string, order_id string) error
	GetOrders(user_id string) (*[]OrderAccuralType, error)
	GetBalance(user_id string) (BalanceType, error)
	PostWithdraw(user_id string, order_id string, sum float64) (string, error)
	GetWithdrawals(user_id string) (float64, error)
	Ping() error
}

var ErrPKConflict = errors.New("PK conflict")
