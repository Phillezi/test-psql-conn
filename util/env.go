package util

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetEnvAsInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	var intValue int
	fmt.Sscanf(value, "%d", &intValue)
	return intValue
}

func GetEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	var durationValue time.Duration
	fmt.Sscanf(value, "%d", &durationValue)
	return durationValue
}

func LogEnvTable() {
	logrus.Infoln("Environment Variables (ENV) and Descriptions (Desc)")

	fmt.Printf("%-20s %-40s %-10s\n", "ENV", "Desc", "Default")

	printEnvVar("DB_HOST", "The database hostname", "localhost")
	printEnvVar("DB_PORT", "The database port", "5432")
	printEnvVar("DB_USER", "The database user", "N/A")
	printEnvVar("DB_PASS", "The database password", "N/A")
	printEnvVar("DB", "The database name", "N/A")
}

func printEnvVar(envVar, description, defaultValue string) {
	fmt.Printf("%-20s %-40s %-10s\n", envVar, description, defaultValue)
}
