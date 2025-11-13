/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"cli/pkg/config"
	"cli/pkg/util"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete existing resources",
	Long: `Delete existing resources from the system.

This command allows you to delete devices, test sessions, conditions,
condition values, scenarios, and scenario conditions.`,
}

var deleteDeviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Delete a device",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("DELETE", config.GetAPIEndpoint()+"/device/"+strconv.Itoa(id), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}

		fmt.Printf("Successfully deleted device with ID %d\n", id)
	},
}

var deleteTestSessionCmd = &cobra.Command{
	Use:   "test-session",
	Short: "Delete a test session",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("DELETE", config.GetAPIEndpoint()+"/test-session/"+strconv.Itoa(id), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}

		fmt.Printf("Successfully deleted test session with ID %d\n", id)
	},
}

var deleteConditionCmd = &cobra.Command{
	Use:   "condition",
	Short: "Delete a condition",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("DELETE", config.GetAPIEndpoint()+"/condition/"+strconv.Itoa(id), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}

		fmt.Printf("Successfully deleted condition with ID %d\n", id)
	},
}

var deleteConditionValueCmd = &cobra.Command{
	Use:   "condition-value",
	Short: "Delete a condition value",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("DELETE", config.GetAPIEndpoint()+"/condition-value/"+strconv.Itoa(id), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}

		fmt.Printf("Successfully deleted condition value with ID %d\n", id)
	},
}

var deleteScenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Delete a scenario",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("DELETE", config.GetAPIEndpoint()+"/scenario/"+strconv.Itoa(id), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}

		fmt.Printf("Successfully deleted scenario with ID %d\n", id)
	},
}

var deleteScenarioConditionCmd = &cobra.Command{
	Use:   "scenario-condition",
	Short: "Delete a scenario condition",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("DELETE", config.GetAPIEndpoint()+"/scenario-condition/"+strconv.Itoa(id), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}

		fmt.Printf("Successfully deleted scenario condition with ID %d\n", id)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteDeviceCmd.Flags().IntP("id", "i", 0, "The ID of the device to delete")
	deleteCmd.AddCommand(deleteDeviceCmd)

	deleteTestSessionCmd.Flags().IntP("id", "i", 0, "The ID of the test session to delete")
	deleteCmd.AddCommand(deleteTestSessionCmd)

	deleteConditionCmd.Flags().IntP("id", "i", 0, "The ID of the condition to delete")
	deleteCmd.AddCommand(deleteConditionCmd)

	deleteConditionValueCmd.Flags().IntP("id", "i", 0, "The ID of the condition value to delete")
	deleteCmd.AddCommand(deleteConditionValueCmd)

	deleteScenarioCmd.Flags().IntP("id", "i", 0, "The ID of the scenario to delete")
	deleteCmd.AddCommand(deleteScenarioCmd)

	deleteScenarioConditionCmd.Flags().IntP("id", "i", 0, "The ID of the scenario condition to delete")
	deleteCmd.AddCommand(deleteScenarioConditionCmd)
}
