package config

import (
	"fmt"
	"os"
	"path"

	"github.com/Phillezi/test-psql-conn/pkg/psql"
	"github.com/Phillezi/test-psql-conn/util"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPass     string
	DBName     string
	ServeHTTP  bool
	HTTPPort   int
	StackTrace bool
}

func InitConfig() {

	// Set up config file settings
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(GetConfigPath()) // User config dir
	viper.AddConfigPath(".")             // Current directory
	viper.AutomaticEnv()                 // Use environment variables

	// Set default values
	viper.SetDefault("db_host", getDefaultDBHost())
	viper.SetDefault("db_port", psql.DefaultDBPort)
	viper.SetDefault("db_user", psql.DefaultDBUser)
	viper.SetDefault("db_pass", psql.DefaultDBPass)
	viper.SetDefault("db_name", psql.DefaultDBName)
	viper.SetDefault("serve_http", true)
	viper.SetDefault("http_port", 8080)
	viper.SetDefault("stacktrace", false)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		zap.L().Debug("Config file not found, using defaults or environment variables.")
	} else {
		zap.L().Debug("Using config file", zap.String("file", viper.ConfigFileUsed()))
	}
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     viper.GetString("db_host"),
		DBPort:     viper.GetInt("db_port"),
		DBUser:     viper.GetString("db_user"),
		DBPass:     viper.GetString("db_pass"),
		DBName:     viper.GetString("db_name"),
		ServeHTTP:  viper.GetBool("serve_http"),
		HTTPPort:   viper.GetInt("http_port"),
		StackTrace: viper.GetBool("stacktrace"),
	}
}

func getDefaultDBHost() string {
	if os.Getenv("IS_DOCKER") == "true" {
		return "172.17.0.1"
	}
	return "localhost"
}

func GetConfigPath() string {
	basePath, err := os.UserConfigDir()
	if err != nil {
		zap.L().Error("Error getting config path", zap.Error(err))
		return "."
	}
	configPath := path.Join(basePath, ".test-psql-conn")
	fileDescr, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
			zap.L().Error("Error creating config directory", zap.Error(err))
			return "."
		}
	} else if err != nil || !fileDescr.IsDir() {
		zap.L().Error("Invalid config directory", zap.Error(err))
		return "."
	}
	return configPath
}

func (c *Config) LogConfig() {
	fmt.Println("Loaded Configuration:")
	fmt.Printf("%-20s %-40s %-10s\n", "ENV", "Desc", "Value")

	printConfigVar("DB_HOST", "The database hostname", c.DBHost)
	printConfigVar("DB_PORT", "The database port", fmt.Sprintf("%d", c.DBPort))
	printConfigVar("DB_USER", "The database user", c.DBUser)
	printConfigVar("DB_PASS", "The database password", c.DBPass)
	printConfigVar("DB_NAME", "The database name", c.DBName)
	printConfigVar("SERVE_HTTP", "Serve the status on /", util.BoolToStr(c.ServeHTTP))
	printConfigVar("HTTP_PORT", "The port to serve http on", fmt.Sprintf("%d", c.HTTPPort))
	printConfigVar("STACKTRACE", "Show the stack trace in error logs", util.BoolToStr(c.StackTrace))
}

func printConfigVar(envVar, description, value string) {
	fmt.Printf("%-20s %-40s %-10s\n", envVar, description, value)
}
