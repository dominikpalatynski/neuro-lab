/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"cli/pkg/config"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Neuro Lab CLI - Manage devices and test sessions",
	Long: `Neuro Lab CLI is a command-line tool for managing devices,
test sessions, conditions, and scenarios for neurological testing.`,
	// PersistentPreRunE is called after flags are parsed but before
	// the command's RunE function is called.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization for the init command to avoid requiring config
		if cmd.Name() == "init" {
			return nil
		}

		// Initialize config system
		if err := config.Initialize(cfgFile); err != nil {
			return err
		}

		// Check if discovery cache is valid, fetch if not
		if !config.IsDiscoveryCacheValid() {
			// Cache is invalid or doesn't exist - fetch from API
			apiEndpoint := config.GetAPIEndpoint()
			if err := config.FetchAndCacheDiscovery(apiEndpoint); err != nil {
				fmt.Println("Error fetching and caching discovery:", err)
			}
		}

		return nil
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
	// Add the persistent --config flag to the root command
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.neurolab/config.yaml)")
}
