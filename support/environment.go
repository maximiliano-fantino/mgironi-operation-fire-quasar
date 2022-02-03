package support

import (
	"log"
	"os"
)

func WebServerPort() string {
	return getEnv("PORT", "8080")
}

func RedisHost() string {
	return getEnv("REDISHOST", "redis")
}

func RedisPort() string {
	return getEnv("REDISPORT", "6379")
}

func getEnv(envkey string, envDefaultValue string) string {
	value := os.Getenv(envkey)
	if value == "" {
		log.Printf("WARN env variable %s not present. Setting default to '%s'", envkey, envDefaultValue)
		value = envDefaultValue
	}
	return value
}
