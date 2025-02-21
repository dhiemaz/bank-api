package cmd

import (
	"fmt"
	"github.com/dhiemaz/bank-api/cmd/gapi"
	"github.com/dhiemaz/bank-api/cmd/gateway"
	"github.com/dhiemaz/bank-api/cmd/migration"
	"github.com/dhiemaz/bank-api/cmd/rest"
	"github.com/dhiemaz/bank-api/config"
	"github.com/dhiemaz/bank-api/infrastructure/logger"
	"github.com/spf13/cobra"
)

type Command struct {
	rootCmd *cobra.Command
}

var text = `BANK API`

// NewCommandEngine the command line boot loader
func NewCommand() *Command {
	var rootCmd = &cobra.Command{
		Use:   "BANK API",
		Short: "Bank API made with Go",
		Long:  "Bank API made with Go",
	}

	return &Command{
		rootCmd: rootCmd,
	}
}

// Run the all command line
func (c *Command) Run() {
	var rootCommands = []*cobra.Command{
		{
			Use:   "rest",
			Short: "Run Banking API HTTP server (rest-API)",
			Long:  "Run Banking API HTTP server (rest-API)",
			PreRun: func(cmd *cobra.Command, args []string) {
				// Show display text
				fmt.Println(fmt.Sprintf(text))
				config.InitLogger()
			},
			Run: func(cmd *cobra.Command, args []string) {
				rest.Run()
				logger.WithFields(logger.Fields{"component": "command", "action": "serve HTTP rest API"}).
					Infof("PreRun command done")
			},
			PostRun: func(cmd *cobra.Command, args []string) {
				logger.WithFields(logger.Fields{"component": "command", "action": "serve HTTP rest API"}).
					Infof("PostRun command done")
			},
		},
		{
			Use:   "gapi",
			Short: "Run Banking API HTTP server (gRPC)",
			Long:  "Run Banking API HTTP server (gRPC)",
			PreRun: func(cmd *cobra.Command, args []string) {
				// Show display text
				fmt.Println(fmt.Sprintf(text))
				config.InitLogger()
			},
			Run: func(cmd *cobra.Command, args []string) {
				gapi.RunGRPCAPIServer()
				logger.WithFields(logger.Fields{"component": "command", "action": "serve HTTP gRPC"}).
					Infof("PreRun command done")
			},
			PostRun: func(cmd *cobra.Command, args []string) {
				logger.WithFields(logger.Fields{"component": "command", "action": "serve HTTP gRPC"}).
					Infof("PostRun command done")
			},
		},
		{
			Use:   "gateway",
			Short: "Run Banking API HTTP server (gRPC Gateway)",
			Long:  "Run Banking API HTTP server (gRPC Gateway)",
			PreRun: func(cmd *cobra.Command, args []string) {
				// Show display text
				fmt.Println(fmt.Sprintf(text))
				config.InitLogger()
			},
			Run: func(cmd *cobra.Command, args []string) {
				gateway.RunGateway()
				logger.WithFields(logger.Fields{"component": "command", "action": "serve HTTP gRPC Gateway"}).
					Infof("PreRun command done")
			},
			PostRun: func(cmd *cobra.Command, args []string) {
				logger.WithFields(logger.Fields{"component": "command", "action": "serve HTTP gRPC Gateway"}).
					Infof("PostRun command done")
			},
		},
		{
			Use:   "migrate",
			Short: "Run Banking API migration",
			Long:  "Run Banking API migration",
			PreRun: func(cmd *cobra.Command, args []string) {
				// Show display text
				fmt.Println(fmt.Sprintf(text))
				config.InitLogger()

				logger.WithFields(logger.Fields{"component": "command", "action": "running database migration"}).
					Infof("starting migration...")
			},
			Run: func(cmd *cobra.Command, args []string) {
				migration.RunMigration()
			},
			PostRun: func(cmd *cobra.Command, args []string) {
				logger.WithFields(logger.Fields{"component": "command", "action": "running database migration"}).
					Infof("done migration...")
			},
		},
	}

	for _, command := range rootCommands {
		c.rootCmd.AddCommand(command)
	}

	c.rootCmd.Execute()
}

// GetRoot the command line service
func (c *Command) GetRoot() *cobra.Command {
	return c.rootCmd
}
