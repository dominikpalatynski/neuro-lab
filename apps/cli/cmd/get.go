/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"cli/pkg/config"
	"cli/pkg/util"

	"database"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources from the system",
	Long: `Get resources from the system.

This command allows you to retrieve devices, test sessions, conditions,
condition values, scenarios, and scenario conditions.`,
}

var getDeviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Get a device by ID",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("GET", config.GetAPIEndpoint()+"/device/"+strconv.Itoa(id), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var device database.Device
		err = json.Unmarshal(resp.Body, &device)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", device)
	},
}

var getDevicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Get all devices",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := util.SendRequest("GET", config.GetAPIEndpoint()+"/device/", nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var devices []database.Device
		err = json.Unmarshal(resp.Body, &devices)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", devices)
	},
}

var getTestSessionCmd = &cobra.Command{
	Use:   "test-session",
	Short: "Get a test session by ID",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("GET", config.GetAPIEndpoint()+"/test-session/"+strconv.Itoa(id), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var testSession database.TestSession
		err = json.Unmarshal(resp.Body, &testSession)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", testSession)
	},
}

var getTestSessionsCmd = &cobra.Command{
	Use:   "test-sessions",
	Short: "Get test sessions by device ID",
	Run: func(cmd *cobra.Command, args []string) {
		deviceID, err := cmd.Flags().GetInt("device-id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("GET", config.GetAPIEndpoint()+"/test-session/list/"+strconv.Itoa(deviceID), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var testSessions []database.TestSession
		err = json.Unmarshal(resp.Body, &testSessions)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", testSessions)
	},
}

var getConditionCmd = &cobra.Command{
	Use:   "condition",
	Short: "Get a condition by ID",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("GET", config.GetAPIEndpoint()+"/condition/"+strconv.Itoa(id), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var condition database.Condition
		err = json.Unmarshal(resp.Body, &condition)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", condition)
	},
}

var getConditionsCmd = &cobra.Command{
	Use:   "conditions",
	Short: "Get all conditions",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := util.SendRequest("GET", config.GetAPIEndpoint()+"/condition/", nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var conditions []database.Condition
		err = json.Unmarshal(resp.Body, &conditions)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", conditions)
	},
}

var getConditionValueCmd = &cobra.Command{
	Use:   "condition-value",
	Short: "Get a condition value by ID",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("GET", config.GetAPIEndpoint()+"/condition-value/"+strconv.Itoa(id), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var conditionValue database.ConditionValue
		err = json.Unmarshal(resp.Body, &conditionValue)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", conditionValue)
	},
}

var getConditionValuesCmd = &cobra.Command{
	Use:   "condition-values",
	Short: "Get all condition values or filter by condition ID",
	Run: func(cmd *cobra.Command, args []string) {
		conditionID, err := cmd.Flags().GetInt("condition-id")
		var endpoint string
		if err == nil && conditionID > 0 {
			// Filter by condition ID
			endpoint = config.GetAPIEndpoint() + "/condition-value/list/" + strconv.Itoa(conditionID)
		} else {
			// Get all
			endpoint = config.GetAPIEndpoint() + "/condition-value/"
		}
		resp, err := util.SendRequest("GET", endpoint, nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var conditionValues []database.ConditionValue
		err = json.Unmarshal(resp.Body, &conditionValues)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", conditionValues)
	},
}

var getScenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Get a scenario by ID (includes related data)",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("GET", config.GetAPIEndpoint()+"/scenario/"+strconv.Itoa(id), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var scenario database.Scenario
		err = json.Unmarshal(resp.Body, &scenario)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", scenario)
	},
}

var getScenariosCmd = &cobra.Command{
	Use:   "scenarios",
	Short: "Get scenarios by test session ID",
	Run: func(cmd *cobra.Command, args []string) {
		testSessionID, err := cmd.Flags().GetInt("test-session-id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("GET", config.GetAPIEndpoint()+"/scenario/list/"+strconv.Itoa(testSessionID), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var scenarios []database.Scenario
		err = json.Unmarshal(resp.Body, &scenarios)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", scenarios)
	},
}

var getScenarioConditionCmd = &cobra.Command{
	Use:   "scenario-condition",
	Short: "Get a scenario condition by ID",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		resp, err := util.SendRequest("GET", config.GetAPIEndpoint()+"/scenario-condition/"+strconv.Itoa(id), nil)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var scenarioCondition database.ScenarioCondition
		err = json.Unmarshal(resp.Body, &scenarioCondition)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", scenarioCondition)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getDeviceCmd.Flags().IntP("id", "i", 0, "The ID of the device")
	getCmd.AddCommand(getDeviceCmd)

	getCmd.AddCommand(getDevicesCmd)

	getTestSessionCmd.Flags().IntP("id", "i", 0, "The ID of the test session")
	getCmd.AddCommand(getTestSessionCmd)

	getTestSessionsCmd.Flags().IntP("device-id", "d", 0, "The device ID to filter by")
	getCmd.AddCommand(getTestSessionsCmd)

	getConditionCmd.Flags().IntP("id", "i", 0, "The ID of the condition")
	getCmd.AddCommand(getConditionCmd)

	getCmd.AddCommand(getConditionsCmd)

	getConditionValueCmd.Flags().IntP("id", "i", 0, "The ID of the condition value")
	getCmd.AddCommand(getConditionValueCmd)

	getConditionValuesCmd.Flags().IntP("condition-id", "c", 0, "The condition ID to filter by (optional)")
	getCmd.AddCommand(getConditionValuesCmd)

	getScenarioCmd.Flags().IntP("id", "i", 0, "The ID of the scenario")
	getCmd.AddCommand(getScenarioCmd)

	getScenariosCmd.Flags().IntP("test-session-id", "t", 0, "The test session ID to filter by")
	getCmd.AddCommand(getScenariosCmd)

	getScenarioConditionCmd.Flags().IntP("id", "i", 0, "The ID of the scenario condition")
	getCmd.AddCommand(getScenarioConditionCmd)
}
