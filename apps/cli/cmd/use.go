/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [device]",
	Short: "Switch to a device",
	Long:  `Use a device to run tests`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		configPath := filepath.Join(homeDir, ".neurolab", "config")
		var config Config
		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if len(args) == 0 {
			fmt.Println("Device name is required")
			return
		}

		deviceName := args[0]

		found := false
		for _, d := range config.Devices {
			if d.Device.Name == deviceName {
				config.CurrentDevice = d.Device.Name
				found = true
				break
			}
		}

		if !found {
			fmt.Println("Device not found")
			return
		}

		yamlData, err := yaml.Marshal(&config)
		if err != nil {
			fmt.Printf("Error creating YAML configuration: %v\n", err)
			return
		}

		if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
			fmt.Printf("Error writing config file: %v\n", err)
			return
		}

		fmt.Println("Current device: ", config.CurrentDevice)
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
