package main

import (
	"context"
	"github.com/evstrART/taskFlow"
	"github.com/evstrART/taskFlow/pkg/handler"
	"github.com/evstrART/taskFlow/pkg/repository"
	"github.com/evstrART/taskFlow/pkg/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := InitConfig(); err != nil {
		logrus.Fatalf("Error reading config file, %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("Error connecting to database, %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(taskFlow.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Error occurred while running HTTP server: %s", err.Error())
		}
	}()

	logrus.Println("Server started")

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Println("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("Error occurred while shutting down server, %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Fatalf("Error occurred while closing database connection, %s", err.Error())
	}

}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
