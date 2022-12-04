// Package cmd implements the service CLI.
// It will run the sensor-api server.
package cmd

/*
Copyright (c) 2021 John Dewey <john@dewey.ws>

*/

import (
	"os"

	"github.com/retr0h/task-migrate/pkg/migrate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	directory string
	verbose   bool
	m         *migrate.Migration
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "task-migrate",
	Short: "A Task file migratio tool",
	Long: `A general purpose migration tool built around [Task][] runner.

_/_ _   _ /___ _ _  . _  _ _ _/_ _
/  /_|_\ /\   / / // /_// /_|/  /_'
                     _/
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger := &logrus.Logger{
			Out:       os.Stderr,
			Formatter: &logrus.JSONFormatter{PrettyPrint: true},
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.InfoLevel,
		}

		db, err := migrate.OpenDB()
		if err != nil {
			m.Logger.WithError(err).Fatal("failed to get db")
		}

		repository := migrate.NewRepository(db)
		m = migrate.NewMigration(logger, repository, directory, verbose)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().
		StringVarP(&directory, "directory", "d", "versions", "Directory containing task files.")
	rootCmd.PersistentFlags().
		BoolVarP(&verbose, "verbose", "v", false, "Verbose")
}
