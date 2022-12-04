package cmd

/*
Copyright (c) 2021 John Dewey <john@dewey.ws>

*/

import (
	"github.com/spf13/cobra"
)

// statusCmd represents the status command.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long: `Show migration status.
`,
	Run: func(cmd *cobra.Command, args []string) {
		t, err := m.Status()
		if err != nil {
			m.Logger.WithError(err).Fatal("failed to render status")
		}
		t.Render()
	},
}

// upCmd represents the up command.
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Applies migrations",
	Long: `Migrates to the most recent version available.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := m.Up()
		if err != nil {
			m.Logger.WithError(err).Fatal("failed to run migrations")
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(statusCmd)
}
