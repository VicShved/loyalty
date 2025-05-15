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

func (r GormRepository) GetOrders(userID uint) (*[]OrderAccrual, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var orders []OrderAccrual
	result := r.DB.WithContext(ctx).Table("orders").Select("orders.order_number, orders.status, orders.updated_at, transactions.value as Accrual").Joins("left join transactions on transactions.order_id = orders.id").Order(`"orders"."updated_at" DESC`).Where(`"orders"."user_id" = ? AND "transactions"."transaction_type" = ?`, userID, "a").Scan(&orders)
	logger.Log.Debug("(r GormRepository) GetOrders", zap.Any("orders", orders))
	return &orders, result.Error
}

func (r GormRepository) GetBalanceWithDrawn(userID uint) (BalanceWithDrawnType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var transactions []Transaction
	var current float32
	var withdrawn float32
	result := r.DB.WithContext(ctx).Table("transactions").Where(`"transactions"."order_id" IN (?)`, r.DB.Table("orders").Where(Order{UserID: userID}).Select(`"orders"."id"`)).Scan(&transactions)
	if result.Error != nil {
		return BalanceWithDrawnType{}, result.Error
	}
	for _, trans := range transactions {
		current += trans.Value
		if trans.TransactionType == "w" {
			withdrawn += trans.Value
		}
	}
	return BalanceWithDrawnType{Current: current, Withdrawn: withdrawn}, nil
}

func (r GormRepository) GetBalance(userID uint) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var current float32
	result := r.DB.WithContext(ctx).Table("transactions").Where(`"transactions"."order_id" IN (?)`, r.DB.Table("orders").Where(Order{UserID: userID}).Select(`"orders"."id"`)).Select(`Coalesce(sum("transactions"."value"), 0) as sm`).Pluck("sm", &current)
	if result.Error != nil {
		return current, result.Error
	}
	return current, nil
}

func (r GormRepository) SaveWithDraw(userID uint, orderNumber string, withDrawSum float32) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return 0, err
	}
	var order Order
	result := r.DB.WithContext(ctx).Where("order_number = ? AND user_id = ?", orderNumber, userID).First(&order)
	logger.Log.Debug("", zap.Any("order", order))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Если заказ у другого пользователя
			return 0, ErrOrderNumberUserConflict
		}
		return 0, result.Error
	}
	result = r.DB.WithContext(ctx).Create(&Transaction{OrderID: order.ID, Value: withDrawSum, TransactionType: "w"})
	if result.Error != nil {
		return 0, result.Error
	}
	var current float32
	result = r.DB.WithContext(ctx).Table("transactions").Where(`"transactions"."order_id" IN (?)`, r.DB.Table("orders").Where(Order{UserID: userID}).Select(`"orders"."id"`)).Select(`Coalesce(sum("transactions"."value"), 0) as sm`).Pluck("sm", &current)
	if result.Error != nil {
		return 0, result.Error
	}
	if (current - withDrawSum) < 0 {
		tx.Rollback()
		return current, nil
	}
	err := tx.Commit().Error
	return current, err

}

func (r GormRepository) GetWithdrawals(userID uint) (*[]OrderTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var orderTrs []OrderTransaction
	result := r.DB.WithContext(ctx).Table("orders").Select("orders.order_number, orders.updated_at, transactions.value as Sum, transactions.transaction_type").Joins("join transactions on transactions.order_id = orders.id").Order(`"orders"."updated_at" DESC`).Where(`"orders"."user_id" = ? AND "transactions"."transaction_type" = ?`, userID, "w").Scan(&orderTrs)
	return &orderTrs, result.Error
}

func (r GormRepository) UpdateOrderStatus(orderNumber string, userID uint, status string, accrual float32) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	updateResult := r.DB.WithContext(ctx).Model(&Order{}).Where("order_number = ? AND user_id = ?", orderNumber, userID).Updates(Order{Status: status})
	if updateResult.Error != nil {
		return updateResult.Error
	}
	if status == "PROCESSED" {
		order, err := r.GetOrderbyNumber(orderNumber)
		if err != nil {
			tx.Rollback()
			return err
		}
		trans := Transaction{OrderID: (*order).ID, TransactionType: "a", Value: accrual}
		saveResult := r.DB.WithContext(ctx).Create(&trans)
		if saveResult.Error != nil {
			tx.Rollback()
			return saveResult.Error
		}
	}
	return tx.Commit().Error
}

func (r GormRepository) GetOrderbyNumber(orderNumber string) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var order Order
	result := r.DB.WithContext(ctx).Where("order_number = ?", orderNumber).First(&order)
	return &order, result.Error
}
