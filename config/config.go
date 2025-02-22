package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	AiKey    string
	HttpPort string
	LogLevel string
}

// Load loads environment vars and inflates Config
func Load() Config {
	if err := godotenv.Load("/app/.env"); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Print("No .env file found")
		}
		log.Print("No .env file found")
	}

	config := Config{}
	fmt.Println(cast.ToString(getOrReturnDefault("AI_KEY", "aikey")))
	config.AiKey = cast.ToString(getOrReturnDefault("AI_KEY", "aikey"))
	config.HttpPort = cast.ToString(getOrReturnDefault("HTTP_PORT", "8081"))
	config.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))

	return config
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
