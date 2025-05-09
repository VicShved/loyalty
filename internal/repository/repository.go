package repository

import (
	"errors"
)

type BalanceWithDrawnType struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type RepoInterface interface {
	Register(login string, hashPassword string) (uint, error)
	Login(login string, hashPassword string) (uint, error)
	SaveOrder(orderNumber string, userID uint) (Order, bool, error)
	GetOrders(userID uint) (*[]Order, error)
	GetBalanceWithDrawn(user_id uint) (BalanceWithDrawnType, error)
	GetBalance(user_id uint) (float32, error)
	SaveWithDraw(userID uint, orderID string, withDrawSum float32) (float32, error)
	// PostWithdraw(user_id string, order_id string, sum float64) (string, error)
	// GetWithdrawals(user_id string) (float64, error)
	Ping() error
}

var ErrLoginConflict = errors.New("Login conflict")
var ErrLoginPassword = errors.New("Login/Password error")
var ErrOrderNumberUserConflict = errors.New("OrderNumber another user conflict")
