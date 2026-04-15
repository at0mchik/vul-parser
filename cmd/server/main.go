package main

import (
	"os"
	"os/signal"
	"syscall"
	"vul-parser/internal/config"
	"vul-parser/internal/handler"
	"vul-parser/internal/service"
	"vul-parser/pkg/server"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.GetServerConfig()

	service := service.NewService()

	handlers := handler.NewHandler(service)


	srv := new(server.Server)


	go func() {
		if err := srv.Run(cfg.Server.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error while running http server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Println("Shutting down server...")
}