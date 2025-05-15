package cmd

import (
	"os"

	"go.uber.org/zap"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.L().Fatal("Failed", zap.Error(err))
		os.Exit(1)
	}
}
