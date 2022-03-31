package main

import (
	"X-Blog/internal/handlers"
	"X-Blog/internal/repository"
	"X-Blog/internal/server"
	"X-Blog/internal/service"
	"os"

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

	db := repository.NewConnectionPostgresDB(&repository.ConfigPostgres{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handlers := handlers.NewHandler(service)

	srv := new(server.Server)
	if err := srv.RunServer(handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err)
	}

}

func initConfig() error {
	viper.AddConfigPath("source/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
