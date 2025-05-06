package main

import (
	"log"
	"net/http"

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

	orderChan := make(chan string, 1000) // TODO 1000 replace to config
	defer close(orderChan)

	// Bussiness layer (empty)
	serv := service.GetService(repo, config.AccuralSystemAddress, &orderChan)

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
