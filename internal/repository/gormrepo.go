package repository

import (
	"context"
	"errors"
	"time"

	"github.com/VicShved/loyalty/internal/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormRepository struct {
	DB *gorm.DB
}

func GetGormDB(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{TranslateError: true})
	return db, err
}

func GetGormRepo(dns string) (*GormRepository, error) {
	db, err := GetGormDB(dns)
	if err != nil {
		return nil, err
	}
	repo := &GormRepository{
		DB: db,
	}
	err = repo.Migrate()
	if err != nil {
		return nil, err
	}
	return repo, err
}

func (r *GormRepository) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := r.DB.WithContext(ctx).AutoMigrate(&User{}, &Order{}, &Transaction{})
	return err
}

func (r GormRepository) Ping() error {
	sqlDB, _ := r.DB.DB()
	return sqlDB.Ping()
}

func (r GormRepository) Register(login string, hashPassword string) (uint, error) {
	logger.Log.Debug("", zap.String("login", login), zap.String("hashPassword", hashPassword))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	user := User{Login: login, HashPassword: hashPassword}
	result := r.DB.WithContext(ctx).Create(&user)
	if result.Error != nil {
		// проверяем на ошибка дублирования логина
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			logger.Log.Debug("login exists", zap.String("login", login))
			return 0, ErrLoginConflict
		}
		return 0, result.Error
	}
	return user.ID, result.Error
}

func (r GormRepository) Login(login string, hashPassword string) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	user := User{}
	result := r.DB.WithContext(ctx).Where(&User{Login: login, HashPassword: hashPassword}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Log.Debug("login|Password not found", zap.String("login", login), zap.String("hashPassword", hashPassword))
			return 0, ErrLoginPassword
		}
		return 0, result.Error
	}
	logger.Log.Debug("", zap.Any("User", user))
	return user.ID, result.Error
}

func (r GormRepository) SaveOrder(orderNumber string, userID uint) (Order, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	order := Order{OrderNumber: orderNumber, UserID: userID}
	result := r.DB.WithContext(ctx).Create(&order)
	if result.Error != nil {
		// проверяем на ошибку дублирования
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			logger.Log.Debug("OrderNumber exists", zap.String("OrderNumber", orderNumber))
			result = r.DB.WithContext(ctx).Where("user_id = ? AND order_number = ?", userID, orderNumber).First(&order)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					// Если заказ у другого пользователя
					return Order{}, false, ErrOrderNumberUserConflict
				}
				// другая ошибка
				return Order{}, false, result.Error
			}
			// уже есть такой заказ у этого пользователя
			return order, false, nil
		}
		return Order{}, false, result.Error
	}
	return order, true, result.Error
}

func (r GormRepository) GetOrders(userID uint) (*[]Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var orders []Order
	result := r.DB.WithContext(ctx).Order("uploaded_at desc").Where("user_id = ?", userID).Find(&orders)
	logger.Log.Debug("", zap.Any("orders", orders))
	// if result.Error != nil {
	// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 		return Order{}, nil
	// 	}
	return &orders, result.Error
}

// func (r GormRepository) GetBalance() {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
// }

// func (r GormRepository) PostWithdraw() {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
// }

// func (r GormRepository) PostWithdrawals() {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
// }
