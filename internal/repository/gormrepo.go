package repository

import (
	"context"
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
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
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
	err := r.DB.WithContext(ctx).AutoMigrate(&KeyOriginalURL{})
	return err
}

func (r GormRepository) Ping() error {
	sqlDB, _ := r.DB.DB()
	return sqlDB.Ping()
}

func (r GormRepository) Register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return nil
}

func (r GormRepository) Login() {
	logger.Log.Debug("Read", zap.String("UserID", userID))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
}

func (r GormRepository) GetOrders() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
}

func (r GormRepository) PostOrders() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return nil
}

func (r GormRepository) GetBalance() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
}

func (r GormRepository) PostWithdraw() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
}

func (r GormRepository) PostWithdrawals() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
}
