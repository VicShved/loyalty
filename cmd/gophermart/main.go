package main

import (
	"log"
	"net/http"

	"github.com/VicShved/loyalty/internal/accrual"
	"github.com/VicShved/loyalty/internal/common"
	"github.com/VicShved/loyalty/internal/handler"
	"github.com/VicShved/loyalty/internal/logger"
	"github.com/VicShved/loyalty/internal/middware"
	"github.com/VicShved/loyalty/internal/repository"
	"github.com/VicShved/loyalty/internal/service"
	"go.uber.org/zap"
)

func main() {
	// Get app config
	var config = common.GetServerConfig()
	// Init custom logger
	logger.InitLogger(config.LogLevel)

	// repo choice
	var repo repository.RepoInterface
	repo, err := repository.GetGormRepo(config.DBDSN)
	if err != nil {
		panic(err)
	}

	logger.Log.Info("Connect to db", zap.String("DSN", config.DBDSN))

	orderChan := make(chan accrual.OrderUser, 100) // TODO 100 replace to config
	defer close(orderChan)

	accrualService := accrual.GetAccrualService(orderChan, config.AccuralSystemAddress, &repo)
	go accrualService.ProcessAccrualService()

	// Bussiness layer
	serv := service.GetService(repo, config.AccuralSystemAddress, &orderChan)

	// Add unprocessed orders from DB to accrual process
	// err = serv.InitAccrualProcess()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Handlers
	handler := handler.GetHandler(serv)

	// Middlewares chain
	middlewares := []func(http.Handler) http.Handler{
		middware.AuthMiddlewareHeader,
		middware.Logger,
		middware.GzipMiddleware,
	}

	//	Create Router
	router := handler.InitRouter(middlewares)

	// Run server
	server := new(common.Server)
	err = server.Run(config.RunAddress, router)
	if err != nil {
		log.Fatal(err)
	}
}
