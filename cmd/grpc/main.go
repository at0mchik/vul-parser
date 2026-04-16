package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vul-parser/internal/config"
	"vul-parser/pkg/server"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.GetServerConfig()

	grpcServer, err := server.NewServer(cfg.Server.GRPCPort)
	if err != nil {
		logrus.Fatalf("Failed to create server: %v", err)
	}

	// Канал для ошибок сервера
	errChan := make(chan error, 1)

	go func() {
		if err := grpcServer.Start(); err != nil {
			errChan <- err
		}
	}()

	// Канал для сигналов остановки
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		logrus.Fatalf("Server error: %v", err)
	case <-stopChan:
		logrus.Info("Received shutdown signal")
		
		// Graceful shutdown с таймаутом
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		done := make(chan struct{})
		go func() {
			grpcServer.Stop()
			close(done)
		}()
		
		select {
		case <-done:
			logrus.Info("Server stopped gracefully")
		case <-ctx.Done():
			logrus.Info("Server shutdown timeout")
		}
	}
}