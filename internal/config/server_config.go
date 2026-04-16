package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Server struct {
		HttpPort string
		GRPCPort string
	}
}

var instance *ServerConfig
var once sync.Once

// GetServerConfig - функция для получения конфига по "аппаратной части", доступы к бд, порты и тд
func GetServerConfig() *ServerConfig {
	once.Do(func() {
		instance = &ServerConfig{}

		if err := godotenv.Load(); err != nil {
			logrus.Fatalf("error loading env variables: %s", err.Error())
		}

		instance.Server.HttpPort = getEnv("HTTP_SERVER_PORT", "8080")
		instance.Server.GRPCPort = getEnv("GRPC_SERVER_PORT", "9054")
	})

	return instance
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
