/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"cli/pkg/util"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var createDeviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Create a new device",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create device called")
		resp, err := util.SendRequest("POST", "http://localhost:3002/api/v1/device", []byte(`{"name": "test"}`))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Response: ", resp.Body)
	},
}

var createTestSessionCmd = &cobra.Command{
	Use:   "test-session",
	Short: "Create a new test session",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create test session called")
		resp, err := util.SendRequest("POST", "http://localhost:3002/api/v1/test-session", []byte(`{"name": "test", "device_id": 1}`))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Response: ", resp.Body)
	},
}

var createConditionCmd = &cobra.Command{
	Use:   "condition",
	Short: "Create a new condition",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create condition called")
		resp, err := util.SendRequest("POST", "http://localhost:3002/api/v1/condition", []byte(`{"name": "test"}`))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Response: ", resp.Body)
	},
}

var createConditionValueCmd = &cobra.Command{
	Use:   "condition-value",
	Short: "Create a new condition value",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create condition value called")
		resp, err := util.SendRequest("POST", "http://localhost:3002/api/v1/condition-value", []byte(`{"value": "test", "condition_id": 1}`))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Response: ", resp.Body)
	},
}

var createScenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Create a new scenario",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create scenario called")
		resp, err := util.SendRequest("POST", "http://localhost:3002/api/v1/scenario", []byte(`{"name": "test", "test_session_id": 2}`))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Response: ", resp.Body)
	},
}

var createScenarioWithConditionValuesCmd = &cobra.Command{
	Use:   "scenario-with-condition-values",
	Short: "Create a new scenario with condition values",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create scenario with condition values called")
		resp, err := util.SendRequest("POST", "http://localhost:3002/api/v1/scenario/with-condition-values", []byte(`{"name": "test", "test_session_id": 1, "condition_value_ids": [1, 2]}`))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Response: ", resp.Body)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createDeviceCmd)
	createCmd.AddCommand(createTestSessionCmd)
	createCmd.AddCommand(createConditionCmd)
	createCmd.AddCommand(createConditionValueCmd)
	createCmd.AddCommand(createScenarioCmd)
	createCmd.AddCommand(createScenarioWithConditionValuesCmd)
}
