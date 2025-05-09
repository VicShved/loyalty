package service

import (
	"github.com/VicShved/loyalty/internal/common"
	"github.com/VicShved/loyalty/internal/logger"
	"github.com/VicShved/loyalty/internal/repository"
	"go.uber.org/zap"
)

type BatchReqJSON struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchRespJSON struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UserURLRespJSON struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ShortenService struct {
	repo           repository.RepoInterface
	accrualAddress string
	orderChan      chan string
}

func GetService(repo repository.RepoInterface, accrualAddress string, orderChan *chan string) *ShortenService {
	return &ShortenService{repo: repo, accrualAddress: accrualAddress, orderChan: *orderChan}
}

func (s *ShortenService) Ping() error {
	return s.repo.Ping()
}

func (s *ShortenService) Register(login string, password string) (uint, error) {
	userID, err := (*s).repo.Register(login, common.HashSha256(password))
	logger.Log.Debug("", zap.Uint("userID", userID), zap.Any("err", err))
	return userID, err
}

func (s *ShortenService) Login(login string, password string) (uint, error) {
	userID, err := (*s).repo.Login(login, common.HashSha256(password))
	return userID, err
}

func (s *ShortenService) SaveOrder(orderNumber string, userID uint) (repository.Order, bool, error) {
	order, isNew, err := (*s).repo.SaveOrder(orderNumber, userID)
	if isNew {
		s.orderChan <- order.OrderNumber
	}
	return order, isNew, err
}

func (s *ShortenService) GetOrders(userID uint) (*[]repository.Order, error) {
	orders, err := (*s).repo.GetOrders(userID)
	return orders, err

}

func (s *ShortenService) GetBalanceWithDrawn(userID uint) (repository.BalanceWithDrawnType, error) {
	balance, err := (*s).repo.GetBalanceWithDrawn(userID)
	return balance, err

}

func (s *ShortenService) GetBalance(userID uint) (float32, error) {
	balance, err := (*s).repo.GetBalance(userID)
	return balance, err
}

func (s *ShortenService) SaveWithDraw(userID uint, orderID string, withDrawSum float32) (float32, error) {
	current, err := (*s).repo.SaveWithDraw(userID, orderID, -withDrawSum)
	return current, err
}
