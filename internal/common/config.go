package common

import (
	"flag"
	"os"
)

type ServerConfigStruct struct {
	RunAddress           string
	AccuralSystemAddress string
	DBDSN                string
	SecretKey            string
	LogLevel             string
}

var ServerConfig ServerConfigStruct

func GetServerConfig() *ServerConfigStruct {
	flag.StringVar(&ServerConfig.RunAddress, "a", "localhost:8080", "Run address")
	flag.StringVar(&ServerConfig.AccuralSystemAddress, "r", "", "AccuralSystemAddress")
	// host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable
	flag.StringVar(&ServerConfig.DBDSN, "d", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", "DataBase DSN")
	flag.StringVar(&ServerConfig.SecretKey, "s", "VeryImpotantSecretKey.YesYes", "Secret key")
	flag.StringVar(&ServerConfig.LogLevel, "l", "INFO", "Log level")
	flag.Parse()

	value, exists := os.LookupEnv("RUN_ADDRESS")
	if exists {
		ServerConfig.RunAddress = value
	}

	value, exists = os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if exists {
		ServerConfig.AccuralSystemAddress = value
	}

	value, exists = os.LookupEnv("DATABASE_DSN")
	if exists {
		ServerConfig.DBDSN = value
	}

	value, exists = os.LookupEnv("SECRET_KEY")
	if exists {
		ServerConfig.SecretKey = value
	}

	value, exists = os.LookupEnv("LOG_LEVEL")
	if exists {
		ServerConfig.LogLevel = value
	}

	return &ServerConfig
}
