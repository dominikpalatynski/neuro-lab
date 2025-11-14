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

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var completeScenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Complete a scenario",
	Long:  `Complete a scenario.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		scenarioId, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}
		resp, err := util.SendRequest("POST", config.GetAPIEndpoint()+"/scenario/complete/"+strconv.Itoa(scenarioId), nil)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("failed to complete scenario: %s", string(resp.Body))
		}

		fmt.Printf("Successfully completed scenario with ID %d\n", scenarioId)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)

	completeScenarioCmd.Flags().IntP("id", "i", 0, "The ID of the scenario to complete")
	completeCmd.AddCommand(completeScenarioCmd)
}
