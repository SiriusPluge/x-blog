package main

import (
	"os"
	"voting-app/internal/handlers"
	"voting-app/internal/repository"
	"voting-app/internal/server"
	"voting-app/internal/service"

	loggo "github.com/bukerdevid/log-go-log"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	loggo.InitCastomLogger(&logrus.JSONFormatter{TimestampFormat: "15:04:05 02/01/2006"}, logrus.TraceLevel, false, true)

	if errConf := initConfig(); errConf != nil {
		logrus.Fatalf("error initializating configs: %s", errConf.Error())
	}

	if errLoadEnv := godotenv.Load(); errLoadEnv != nil {
		logrus.Fatalf("error loading env variables: %s", errLoadEnv.Error())
	}

	hlContract := initLedger()

	db := repository.NewConnectionPostgresDB(&repository.ConfigPostgres{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	// ethClient, err := ethclient.Dial(viper.GetString("blockchain.node"))
	// if err != nil {
	// 	logrus.Fatalf("error connect to ethereum: %s", err)
	// }

	repository := repository.NewRepository(db)
	service := service.NewService(repository, hlContract)
	handlers := handlers.NewHandler(service)

	srv := new(server.Server)
	if err := srv.RunServer(handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err)
	}

}
