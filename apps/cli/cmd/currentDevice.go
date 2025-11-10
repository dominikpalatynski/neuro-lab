/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cli/pkg/config"
	"fmt"

	"github.com/spf13/cobra"
)

// currentDeviceCmd represents the currentDevice command
var currentDeviceCmd = &cobra.Command{
	Use:   "current-device",
	Short: "Display the currently selected device",
	Long: `Display the currently selected device context.
This device is used for creating test sessions and other operations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get full device info
		deviceInfo, err := config.GetCurrentDeviceInfo()
		if err != nil {
			return fmt.Errorf("no device selected: %w\nRun 'cli init' to fetch devices and 'cli use <device-name>' to select one", err)
		}

		fmt.Printf("Current device: %s (ID: %d)\n", deviceInfo.Name, deviceInfo.DeviceID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(currentDeviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// currentDeviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// currentDeviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
