package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/Phillezi/test-psql-conn/util"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DBHost        string
	DBPort        int
	DBUser        string
	DBPass        string
	DBName        string
	SleepWhenDone bool
	ServeHTTP     bool
}

func Load() *Config {
	isDocker := os.Getenv("IS_DOCKER")
	// Default DBHost is localhost, but change to Docker internal IP if IS_DOCKER is set
	defaultHost := "localhost"
	if isDocker == "true" {
		defaultHost = "172.17.0.1"
	}
	config := &Config{
		DBHost:        util.GetEnv("DB_HOST", defaultHost),
		DBPort:        util.GetEnvAsInt("DB_PORT", 5432),
		DBUser:        util.GetEnv("DB_USER", "N/A"),
		DBPass:        util.GetEnv("DB_PASS", "N/A"),
		DBName:        util.GetEnv("DB_NAME", "N/A"),
		SleepWhenDone: strings.ToLower(util.GetEnv("SLEEP_INF", "true")) == "true",
		ServeHTTP:     strings.ToLower(util.GetEnv("SERVE_HTTP", "true")) == "true",
	}

	return config
}

func (c *Config) LogConfig() {
	logrus.Infoln("Loaded Configuration:")
	fmt.Printf("%-20s %-40s %-10s\n", "ENV", "Desc", "Value")

	printConfigVar("DB_HOST", "The database hostname", c.DBHost)
	printConfigVar("DB_PORT", "The database port", fmt.Sprintf("%d", c.DBPort))
	printConfigVar("DB_USER", "The database user", c.DBUser)
	printConfigVar("DB_PASS", "The database password", c.DBPass)
	printConfigVar("DB_NAME", "The database name", c.DBName)

	printConfigVar("SLEEP_INF", "Sleep infinitely when done", func(val bool) string {
		if val {
			return "true"
		}
		return "false"
	}(c.SleepWhenDone))

	printConfigVar("SERVE_HTTP", "Serve the status on /", func(val bool) string {
		if val {
			return "true"
		}
		return "false"
	}(c.ServeHTTP))
}

func printConfigVar(envVar, description, value string) {
	fmt.Printf("%-20s %-40s %-10s\n", envVar, description, value)
}
