/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cli/pkg/config"
	"cli/pkg/util"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var stopScenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Stop a scenario",
	Long:  `Stop a scenario.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		scenarioId, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}
		resp, err := util.SendRequest("POST", config.GetAPIEndpoint()+"/scenario/deactivate/"+strconv.Itoa(scenarioId), nil)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("failed to stop scenario: %s", string(resp.Body))
		}

		fmt.Printf("Successfully stopped scenario with ID %d\n", scenarioId)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopScenarioCmd.Flags().IntP("id", "i", 0, "The ID of the scenario to stop")
	stopCmd.AddCommand(stopScenarioCmd)
}
