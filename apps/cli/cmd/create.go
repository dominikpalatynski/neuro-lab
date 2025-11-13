/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"cli/pkg/config"
	"cli/pkg/util"

	"database"
	"types"

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
		deviceName, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		var req types.CreateDeviceRequest
		req.Name = deviceName
		reqBytes, marshalErr := json.Marshal(req)
		if marshalErr != nil {
			fmt.Println("Error marshaling request: ", marshalErr)
			return
		}
		resp, err := util.SendRequest("POST", config.GetAPIEndpoint()+"/device", reqBytes)
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

var createTestSessionCmd = &cobra.Command{
	Use:   "test-session",
	Short: "Create a new test session",
	Run: func(cmd *cobra.Command, args []string) {
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
		var req types.CreateTestSessionRequest
		req.Name = name
		req.DeviceID = uint(deviceID)
		reqBytes, marshalErr := json.Marshal(req)
		if marshalErr != nil {
			fmt.Println("Error marshaling request: ", marshalErr)
			return
		}
		resp, err := util.SendRequest("POST", config.GetAPIEndpoint()+"/test-session", reqBytes)
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

var createConditionCmd = &cobra.Command{
	Use:   "condition",
	Short: "Create a new condition",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		var req types.CreateConditionRequest
		req.Name = name
		reqBytes, marshalErr := json.Marshal(req)
		if marshalErr != nil {
			fmt.Println("Error marshaling request: ", marshalErr)
			return
		}
		resp, err := util.SendRequest("POST", config.GetAPIEndpoint()+"/condition", reqBytes)
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

var createConditionValueCmd = &cobra.Command{
	Use:   "condition-value",
	Short: "Create a new condition value",
	Run: func(cmd *cobra.Command, args []string) {
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
		var req types.CreateConditionValueRequest
		req.Value = value
		req.ConditionID = uint(conditionID)
		reqBytes, marshalErr := json.Marshal(req)
		if marshalErr != nil {
			fmt.Println("Error marshaling request: ", marshalErr)
			return
		}
		resp, err := util.SendRequest("POST", config.GetAPIEndpoint()+"/condition-value", reqBytes)
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

var createScenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Create a new scenario with optional condition values",
	Run: func(cmd *cobra.Command, args []string) {
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
		conditionValueIDs, err := cmd.Flags().GetIntSlice("condition-value-ids")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		var reqBytes []byte
		var marshalErr error
		var endpoint string

		if len(conditionValueIDs) > 0 {
			// Use scenario with condition values endpoint
			var req types.CreateScenarioWithConditionValuesRequest
			req.Name = name
			req.TestSessionID = uint(testSessionID)
			req.ConditionValueIDs = make([]uint, len(conditionValueIDs))
			for i, id := range conditionValueIDs {
				req.ConditionValueIDs[i] = uint(id)
			}
			reqBytes, marshalErr = json.Marshal(req)
			endpoint = config.GetAPIEndpoint() + "/scenario/with-condition-values"
		} else {
			// Use basic scenario endpoint
			var req types.CreateScenarioRequest
			req.Name = name
			req.TestSessionID = uint(testSessionID)
			reqBytes, marshalErr = json.Marshal(req)
			endpoint = config.GetAPIEndpoint() + "/scenario"
		}

		if marshalErr != nil {
			fmt.Println("Error marshaling request: ", marshalErr)
			return
		}
		resp, err := util.SendRequest("POST", endpoint, reqBytes)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(resp.Body))
			return
		}
		var scenarioResponse database.Scenario
		err = json.Unmarshal(resp.Body, &scenarioResponse)
		if err != nil {
			fmt.Printf("Error unmarshaling response: %v\nResponse body: %s\n", err, string(resp.Body))
			return
		}

		fmt.Println("Response: ", scenarioResponse)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createDeviceCmd.Flags().StringP("name", "n", "", "The name of the device")
	createCmd.AddCommand(createDeviceCmd)

	createTestSessionCmd.Flags().StringP("name", "n", "", "The name of the test session")
	createTestSessionCmd.Flags().IntP("device-id", "d", 0, "The device ID")
	createCmd.AddCommand(createTestSessionCmd)

	createConditionCmd.Flags().StringP("name", "n", "", "The name of the condition")
	createCmd.AddCommand(createConditionCmd)

	createConditionValueCmd.Flags().StringP("value", "v", "", "The condition value")
	createConditionValueCmd.Flags().IntP("condition-id", "c", 0, "The condition ID")
	createCmd.AddCommand(createConditionValueCmd)

	createScenarioCmd.Flags().StringP("name", "n", "", "The name of the scenario")
	createScenarioCmd.Flags().IntP("test-session-id", "t", 0, "The test session ID")
	createScenarioCmd.Flags().IntSliceP("condition-value-ids", "c", []int{}, "The condition value IDs (comma-separated, optional)")
	createCmd.AddCommand(createScenarioCmd)
}
