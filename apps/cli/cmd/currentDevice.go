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

// currentDeviceCmd represents the currentDevice command
var currentDeviceCmd = &cobra.Command{
	Use:   "currentDevice",
	Short: "Show the current device",
	Long:  `Show the current device that is being used`,
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
		fmt.Println("Current device: ", config.CurrentDevice)
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
