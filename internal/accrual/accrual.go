package accrual

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/VicShved/loyalty/internal/logger"
	"github.com/VicShved/loyalty/internal/repository"
	"go.uber.org/zap"
)

type OrderUser struct {
	OrderNumber string
	UserID      uint
}

type AccrualService struct {
	orderChan      chan OrderUser
	accrualAddress string
	repo           repository.RepoInterface
}

func GetAccrualService(orderChan chan OrderUser, accrualAddress string, repo *repository.RepoInterface) *AccrualService {
	return &AccrualService{orderChan: orderChan, accrualAddress: accrualAddress, repo: *repo}
}

type processOrderAccrualBody struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

func (s *AccrualService) processOrderAccrual(orderNumber string, userID uint) {
	client := &http.Client{}
	uri := s.accrualAddress + "/api/orders/" + orderNumber
	logger.Log.Debug("processOrderAccrual", zap.String("uri", uri))
	response, err := client.Get(uri)
	logger.Log.Debug("processOrderAccrual", zap.String("response", response.Status))
	if err != nil {
		s.orderChan <- OrderUser{OrderNumber: orderNumber, UserID: userID}
		logger.Log.Error("processOrderAccrual", zap.String("client.Get", err.Error()))
		return
	}
	if response.StatusCode == http.StatusNoContent {
		// s.orderChan <- OrderUser{OrderNumber: orderNumber, UserID: userID}
		return
	}

	if (response.StatusCode == http.StatusTooManyRequests) || (response.StatusCode == http.StatusInternalServerError) {
		sleepString := response.Header.Get("Retry-After")
		if len(sleepString) > 0 {
			i, err := strconv.Atoi(sleepString)
			if err != nil {
				time.Sleep(time.Second)
			}
			time.Sleep(time.Duration(i) * time.Second)
		} else {
			time.Sleep(time.Second)
		}
		s.orderChan <- OrderUser{OrderNumber: orderNumber, UserID: userID}
		return
	}

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		logger.Log.Error("", zap.String("ReadBody", err.Error()))
		s.orderChan <- OrderUser{OrderNumber: orderNumber, UserID: userID}
		return
	}

	var bodyData processOrderAccrualBody
	err = json.Unmarshal(body, &bodyData)
	logger.Log.Debug("", zap.Any("body", bodyData))
	if err != nil {
		logger.Log.Error("", zap.String("json.Unmarshal", err.Error()))
		s.orderChan <- OrderUser{OrderNumber: orderNumber, UserID: userID}
		return
	}
	if bodyData.Status == "REGISTERED" {
		bodyData.Status = "PROCESSING"
	}
	err = s.repo.UpdateOrderStatus(orderNumber, userID, bodyData.Status, bodyData.Accrual)
	if err != nil {
		s.orderChan <- OrderUser{OrderNumber: orderNumber, UserID: userID}
		return
	}
	if slices.Contains([]string{"PROCESSED", "INVALID"}, bodyData.Status) {
		return
	}
	s.orderChan <- OrderUser{OrderNumber: orderNumber, UserID: userID}

}

func (s *AccrualService) ProcessAccrualService() {
	logger.Log.Info("ProcessAccrualService", zap.String("orderUser", "Enter to goroutine"))
	for orderUser := range s.orderChan {
		logger.Log.Info("ProcessAccrualService", zap.Any("orderUser", orderUser))
		go s.processOrderAccrual(orderUser.OrderNumber, orderUser.UserID)
	}
	logger.Log.Error("ProcessAccrualService", zap.String("orderUser", "EXIT from goroutine"))
}
