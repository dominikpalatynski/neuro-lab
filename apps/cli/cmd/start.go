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

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var startScenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Start a scenario",
	Long:  `Start a scenario.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		scenarioId, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}
		resp, err := util.SendRequest("POST", config.GetAPIEndpoint()+"/scenario/activate/"+strconv.Itoa(scenarioId), nil)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("failed to start scenario: %s", string(resp.Body))
		}

		fmt.Printf("Successfully started scenario with ID %d\n", scenarioId)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startScenarioCmd.Flags().IntP("id", "i", 0, "The ID of the scenario to start")
	startCmd.AddCommand(startScenarioCmd)
}
