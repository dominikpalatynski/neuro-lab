/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cli/pkg/config"
	"fmt"

	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [device-name]",
	Short: "Switch to a specific device",
	Long: `Switch the current device context. This device will be used
for subsequent commands like creating test sessions.

Example:
  cli use device-name`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deviceName := args[0]

		// Set the current device (validates it exists)
		if err := config.SetCurrentDevice(deviceName); err != nil {
			return fmt.Errorf("failed to set current device: %w", err)
		}

		fmt.Printf("✓ Current device set to: %s\n", deviceName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
