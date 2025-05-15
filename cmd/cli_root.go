package cmd

import (
	"sync"

	"github.com/Phillezi/test-psql-conn/config"
	"github.com/Phillezi/test-psql-conn/internal/interrupt"
	"github.com/Phillezi/test-psql-conn/internal/log"
	"github.com/Phillezi/test-psql-conn/pkg/psql"
	"github.com/Phillezi/test-psql-conn/pkg/web"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "test-psql-conn",
	Short: "App to test postgreSQL connections",
	Long: `
  __                   __                               .__                                        
_/  |_  ____   _______/  |_          ______  ___________|  |             ____  ____   ____   ____  
\   __\/ __ \ /  ___/\   __\  ______ \____ \/  ___/ ____/  |    ______ _/ ___\/  _ \ /    \ /    \ 
 |  | \  ___/ \___ \  |  |   /_____/ |  |_> >___ < <_|  |  |__ /_____/ \  \__(  <_> )   |  \   |  \
 |__|  \___  >____  > |__|           |   __/____  >__   |____/          \___  >____/|___|  /___|  /
           \/     \/                 |__|       \/   |__|                   \/           \/     \/ `,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Setup()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		if viper.GetBool("logconfig") {
			cfg.LogConfig()
		}

		dbClient := psql.New().WithContext(
			interrupt.GetInstance().Context(),
		).WithOptions(psql.ClientOpts{
			DBHost: &cfg.DBHost,
			DBPort: &cfg.DBPort,
			DBUser: &cfg.DBUser,
			DBPass: &cfg.DBPass,
			DBName: &cfg.DBName,
		})

		webServer := web.New(web.ServerOpts{
			Port:              &cfg.HTTPPort,
			ConnectionChannel: dbClient.ConnectChannel(),
			TablesChannel:     dbClient.TableChannel(),
		}).WithContext(interrupt.GetInstance().Context())

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			if err := dbClient.Probe(); err != nil {
				zap.L().Error("prober exited with error", zap.Error(err))
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			if err := webServer.Start(); err != nil {
				zap.L().Error("webserver exited with error", zap.Error(err))
			}
			wg.Done()
		}()

		wg.Wait()
	},
}

func init() {

	cobra.OnInitialize(config.InitConfig)

	// Persistent flags
	rootCmd.PersistentFlags().String("loglevel", "info", "Set the logging level (info, warn, error, debug)")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.PersistentFlags().String("profile", "", "Set the logging profile (production or empty)")
	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))

	rootCmd.PersistentFlags().Bool("logconfig", false, "Log the current config on boot")
	viper.BindPFlag("logconfig", rootCmd.PersistentFlags().Lookup("logconfig"))

	rootCmd.PersistentFlags().Bool("stacktrace", false, "Show the stack trace in error logs")
	viper.BindPFlag("stacktrace", rootCmd.PersistentFlags().Lookup("stacktrace"))

	rootCmd.PersistentFlags().String("host", viper.GetString("db_host"), "Database hostname")
	viper.BindPFlag("db_host", rootCmd.PersistentFlags().Lookup("host"))

	rootCmd.PersistentFlags().Int("port", viper.GetInt("db_port"), "Database port")
	viper.BindPFlag("db_port", rootCmd.PersistentFlags().Lookup("port"))

	rootCmd.PersistentFlags().String("user", viper.GetString("db_user"), "Database username")
	viper.BindPFlag("db_user", rootCmd.PersistentFlags().Lookup("user"))

	rootCmd.PersistentFlags().String("pass", viper.GetString("db_pass"), "Database password")
	viper.BindPFlag("db_pass", rootCmd.PersistentFlags().Lookup("pass"))

	rootCmd.PersistentFlags().String("name", viper.GetString("db_name"), "Database name")
	viper.BindPFlag("db_name", rootCmd.PersistentFlags().Lookup("name"))

	rootCmd.PersistentFlags().Bool("serve-http", viper.GetBool("serve_http"), "Serve the status on /")
	viper.BindPFlag("serve_http", rootCmd.PersistentFlags().Lookup("serve-http"))

}
