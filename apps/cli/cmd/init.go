/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cli/pkg/util"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Device represents a device from the API response
type Device struct {
	ID   uint   `json:"ID"`
	Name string `json:"name"`
}

// Config represents the structure of the config file
type Config struct {
	Devices       []DeviceWrapper `yaml:"devices"`
	CurrentDevice string          `yaml:"current_device,omitempty"`
}

// DeviceWrapper wraps a device in the YAML structure
type DeviceWrapper struct {
	Device DeviceInfo `yaml:"device"`
}

// DeviceInfo contains the device information for config
type DeviceInfo struct {
	Name     string `yaml:"name"`
	DeviceID uint   `yaml:"device_id"`
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

		// Fetch devices from API
		apiURL := "http://localhost:3002/api/v1/device"
		body, err := util.SendRequest("GET", apiURL, nil)
		if err != nil {
			fmt.Printf("Error fetching devices from API: %v\n", err)
			return
		}

		// Parse JSON response
		var devices []Device
		if err := json.Unmarshal(body, &devices); err != nil {
			fmt.Printf("Error parsing device response: %v\n", err)
			return
		}

		// Convert to config format
		var deviceWrappers []DeviceWrapper
		for _, device := range devices {
			deviceWrappers = append(deviceWrappers, DeviceWrapper{
				Device: DeviceInfo{
					Name:     device.Name,
					DeviceID: device.ID,
				},
			})
		}

		config := Config{
			Devices: deviceWrappers,
		}

		// Marshal to YAML
		yamlData, err := yaml.Marshal(&config)
		if err != nil {
			fmt.Printf("Error creating YAML configuration: %v\n", err)
			return
		}

		// Get home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			return
		}

		// Create config directory
		configDir := filepath.Join(homeDir, ".neurolab")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("Error creating config directory: %v\n", err)
			return
		}

		// Write config file
		configPath := filepath.Join(configDir, "config")
		if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
			fmt.Printf("Error writing config file: %v\n", err)
			return
		}

		fmt.Printf("✓ Configuration initialized successfully\n")
		fmt.Printf("✓ Found %d device(s)\n", len(devices))
		fmt.Printf("✓ Config saved to: %s\n", configPath)
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
