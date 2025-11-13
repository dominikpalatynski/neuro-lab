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
	"types"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update existing resources",
	Long: `Update existing resources in the system.

This command allows you to update devices, test sessions, conditions,
condition values, and scenarios.`,
}

var updateDeviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Update an existing device",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		deviceName, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		var req types.UpdateDeviceRequest
		req.ID = uint(id)
		req.Name = deviceName
		reqBytes, marshalErr := json.Marshal(req)
		if marshalErr != nil {
			fmt.Println("Error marshaling request: ", marshalErr)
			return
		}
		resp, err := util.SendRequest("PUT", config.GetAPIEndpoint()+"/device/"+strconv.Itoa(id), reqBytes)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var deviceResponse database.Device
		err = json.Unmarshal(resp.Body, &deviceResponse)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", deviceResponse)
	},
}

var updateTestSessionCmd = &cobra.Command{
	Use:   "test-session",
	Short: "Update an existing test session",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		deviceID, err := cmd.Flags().GetInt("device-id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		var req types.UpdateTestSessionRequest
		req.ID = uint(id)
		req.Name = name
		req.DeviceID = uint(deviceID)
		reqBytes, marshalErr := json.Marshal(req)
		if marshalErr != nil {
			fmt.Println("Error marshaling request: ", marshalErr)
			return
		}
		resp, err := util.SendRequest("PUT", config.GetAPIEndpoint()+"/test-session/"+strconv.Itoa(id), reqBytes)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var testSessionResponse database.TestSession
		err = json.Unmarshal(resp.Body, &testSessionResponse)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", testSessionResponse)
	},
}

var updateConditionCmd = &cobra.Command{
	Use:   "condition",
	Short: "Update an existing condition",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		var req types.UpdateConditionRequest
		req.ID = uint(id)
		req.Name = name
		reqBytes, marshalErr := json.Marshal(req)
		if marshalErr != nil {
			fmt.Println("Error marshaling request: ", marshalErr)
			return
		}
		resp, err := util.SendRequest("PUT", config.GetAPIEndpoint()+"/condition/"+strconv.Itoa(id), reqBytes)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var conditionResponse database.Condition
		err = json.Unmarshal(resp.Body, &conditionResponse)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", conditionResponse)
	},
}

var updateConditionValueCmd = &cobra.Command{
	Use:   "condition-value",
	Short: "Update an existing condition value",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		value, err := cmd.Flags().GetString("value")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		conditionID, err := cmd.Flags().GetInt("condition-id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		var req types.UpdateConditionValueRequest
		req.ID = uint(id)
		req.Value = value
		req.ConditionID = uint(conditionID)
		reqBytes, marshalErr := json.Marshal(req)
		if marshalErr != nil {
			fmt.Println("Error marshaling request: ", marshalErr)
			return
		}
		resp, err := util.SendRequest("PUT", config.GetAPIEndpoint()+"/condition-value/"+strconv.Itoa(id), reqBytes)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var conditionValueResponse database.ConditionValue
		err = json.Unmarshal(resp.Body, &conditionValueResponse)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", conditionValueResponse)
	},
}

var updateScenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Update an existing scenario or scenario condition",
	Long: `Update an existing scenario or scenario condition.

If you provide scenario-id and condition-value-id, it will update a scenario condition.
Otherwise, it will update the basic scenario properties (name and test-session-id).`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		// Check if this is a scenario condition update
		scenarioID, scenarioIDErr := cmd.Flags().GetInt("scenario-id")
		conditionValueID, conditionValueIDErr := cmd.Flags().GetInt("condition-value-id")

		var reqBytes []byte
		var marshalErr error
		var endpoint string

		if scenarioIDErr == nil && conditionValueIDErr == nil && scenarioID > 0 && conditionValueID > 0 {
			// Update scenario condition
			var req types.UpdateScenarioConditionRequest
			req.ID = uint(id)
			req.ScenarioID = uint(scenarioID)
			req.ConditionValueID = uint(conditionValueID)
			reqBytes, marshalErr = json.Marshal(req)
			endpoint = config.GetAPIEndpoint() + "/scenario-condition/" + strconv.Itoa(id)
		} else {
			// Update basic scenario
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			testSessionID, err := cmd.Flags().GetInt("test-session-id")
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			var req types.UpdateScenarioRequest
			req.ID = uint(id)
			req.Name = name
			req.TestSessionID = uint(testSessionID)
			reqBytes, marshalErr = json.Marshal(req)
			endpoint = config.GetAPIEndpoint() + "/scenario/" + strconv.Itoa(id)
		}

		if marshalErr != nil {
			fmt.Println("Error marshaling request: ", marshalErr)
			return
		}
		resp, err := util.SendRequest("PUT", endpoint, reqBytes)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}

		// Try to unmarshal as Scenario first, if that fails try ScenarioCondition
		var scenarioResponse database.Scenario
		err = json.Unmarshal(resp.Body, &scenarioResponse)
		if err != nil {
			var scenarioConditionResponse database.ScenarioCondition
			err = json.Unmarshal(resp.Body, &scenarioConditionResponse)
			if err != nil {
				fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
				return
			}
			fmt.Println("Response: ", scenarioConditionResponse)
			return
		}

		fmt.Println("Response: ", scenarioResponse)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateDeviceCmd.Flags().IntP("id", "i", 0, "The ID of the device to update")
	updateDeviceCmd.Flags().StringP("name", "n", "", "The new name of the device")
	updateCmd.AddCommand(updateDeviceCmd)

	updateTestSessionCmd.Flags().IntP("id", "i", 0, "The ID of the test session to update")
	updateTestSessionCmd.Flags().StringP("name", "n", "", "The new name of the test session")
	updateTestSessionCmd.Flags().IntP("device-id", "d", 0, "The new device ID")
	updateCmd.AddCommand(updateTestSessionCmd)

	updateConditionCmd.Flags().IntP("id", "i", 0, "The ID of the condition to update")
	updateConditionCmd.Flags().StringP("name", "n", "", "The new name of the condition")
	updateCmd.AddCommand(updateConditionCmd)

	updateConditionValueCmd.Flags().IntP("id", "i", 0, "The ID of the condition value to update")
	updateConditionValueCmd.Flags().StringP("value", "v", "", "The new condition value")
	updateConditionValueCmd.Flags().IntP("condition-id", "c", 0, "The new condition ID")
	updateCmd.AddCommand(updateConditionValueCmd)

	updateScenarioCmd.Flags().IntP("id", "i", 0, "The ID of the scenario or scenario condition to update")
	updateScenarioCmd.Flags().StringP("name", "n", "", "The new name of the scenario (for scenario updates)")
	updateScenarioCmd.Flags().IntP("test-session-id", "t", 0, "The new test session ID (for scenario updates)")
	updateScenarioCmd.Flags().IntP("scenario-id", "s", 0, "The scenario ID (for scenario condition updates)")
	updateScenarioCmd.Flags().IntP("condition-value-id", "c", 0, "The condition value ID (for scenario condition updates)")
	updateCmd.AddCommand(updateScenarioCmd)
}
