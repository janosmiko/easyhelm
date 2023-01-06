package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"easyhelm/cmd/generate"
	"easyhelm/internal/config"
	"easyhelm/internal/logger"
)

// rootCmd represents the root command.
var rootCmd = &cobra.Command{
	Use:   "easyhelm",
	Short: "EasyHelm",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RootCmd(cmd)
	},
}

// init initializes cobra commands, sets up the persistent command line flags
// and attaches the environment variables to them.
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(generate.Cmd)
}

// Execute runs the rootCmd itself.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

// RootCmd represents the root command.
func RootCmd(cmd *cobra.Command) error {
	printSettings()

	return cmd.Help() // nolint: wrapcheck
}

func initConfig() {
	logger.Init()
	config.Init()
}

// printSettings is a helper function, which prints out all the environment variables.
func printSettings() {
	if viper.GetBool("PRINT_ENVIRONMENT") {
		for i, v := range viper.AllSettings() {
			zap.S().Infof("%v=%v", strings.ToUpper(i), v)
		}

		os.Exit(0)
	}
}
