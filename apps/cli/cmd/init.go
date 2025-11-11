/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cli/pkg/config"
	"cli/pkg/util"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// Device represents a device from the API response
type Device struct {
	ID   uint   `json:"ID"`
	Name string `json:"name"`
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize neuro-lab CLI configuration",
	Long: `Fetches available devices from the API server and creates a configuration file
at ~/.neurolab/config. This configuration file will be used by other CLI commands
to interact with registered devices.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing neuro-lab configuration...")

		// Initialize config system (creates directory if needed)
		if err := config.Initialize(cfgFile); err != nil {
			fmt.Printf("Error initializing config: %v\n", err)
			return
		}

		// Fetch devices from API
		apiEndpoint := config.GetAPIEndpoint()
		apiURL := apiEndpoint + "/device"
		resp, err := util.SendRequest("GET", apiURL, nil)
		if err != nil {
			fmt.Printf("Error fetching devices from API: %v\n", err)
			return
		}

		// Parse JSON response
		var devices []Device
		if err := json.Unmarshal(resp.Body, &devices); err != nil {
			fmt.Printf("Error parsing device response: %v\n", err)
			return
		}

		// Convert to config format
		var deviceWrappers []config.DeviceWrapper
		for _, device := range devices {
			deviceWrappers = append(deviceWrappers, config.DeviceWrapper{
				Device: config.DeviceInfo{
					Name:     device.Name,
					DeviceID: device.ID,
				},
			})
		}

		// Save devices to config
		if err := config.SetDevices(deviceWrappers); err != nil {
			fmt.Printf("Error saving devices to config: %v\n", err)
			return
		}

		// Fetch and cache API resources
		fmt.Println("Fetching API resources...")
		if err := config.FetchAndCacheDiscovery(apiEndpoint); err != nil {
			// Don't fail the whole init if discovery fails
			fmt.Printf("Warning: Failed to fetch API resources: %v\n", err)
			fmt.Println("  (You can continue without discovery cache)")
		} else {
			// Get cached resources to show count
			resources, _ := config.GetDiscoveryResources()
			fmt.Printf("✓ Cached %d API resource type(s)\n", len(resources))
		}

		fmt.Printf("✓ Configuration initialized successfully\n")
		fmt.Printf("✓ Found %d device(s)\n", len(devices))
		fmt.Printf("✓ Config saved to: %s\n", config.ConfigFileUsed())
		fmt.Printf("\nUse 'cli use <device-name>' to select a device\n")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
